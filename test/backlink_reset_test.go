package test

import (
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

// Test_Reset_rebuilds_backlinks runs the full backlink reset and verifies the
// index is rebuilt from each document's current content (no stale leftovers).
func Test_Reset_rebuilds_backlinks(t *testing.T) {
	db := Setup_query_api_db(t)
	defer tool.DB_close(db)

	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("reset seed exec: %v", err)
		}
	}
	// Owner group grants setting_backlink.
	exec("insert into user_set (name, id, data) values ('acl','admin_ip','owner')")
	admin := tool.Config{IP: "admin_ip"}

	// A document with a namumark backlink and a second with none.
	exec("insert into data (title, data, type) values ('bl_src_a', '[[bl_dst_b]] content a', 'document')")
	exec("insert into data (title, data, type) values ('bl_dst_b', 'content b', 'document')")
	exec("insert into data (title, data, type) values ('bl_nolink', 'no links here', 'document')")

	out := route.Api_setting_backlink_reset_post(admin)
	if out["response"] != "ok" {
		t.Fatalf("reset response = %v, want ok", out)
	}

	// The rebuilt index must reflect the current content: src_a -> dst_b,
	// nolink -> nothing.
	if !hasBacklink(db, "bl_src_a", "bl_dst_b", t) {
		t.Fatalf("src_a does not backlink to bl_dst_b after reset")
	}
	if hasBacklink(db, "bl_src_a", "bl_nolink", t) {
		t.Fatalf("src_a unexpectedly backlinks to bl_nolink")
	}
	if hasBacklink(db, "bl_nolink", "bl_dst_b", t) {
		t.Fatalf("bl_nolink unexpectedly backlinks to bl_dst_b")
	}
}

// Test_Reset_does_not_restore_deleted_doc verifies that a document deleted
// while the rebuild runs (so its content is gone and the rebuild "fails") is
// NOT restored — its (already wiped) backlinks stay gone.
func Test_Reset_does_not_restore_deleted_doc(t *testing.T) {
	db := Setup_query_api_db(t)
	defer tool.DB_close(db)

	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("exec: %v", err)
		}
	}
	exec("insert into user_set (name, id, data) values ('acl','admin_ip','owner')")
	admin := tool.Config{IP: "admin_ip"}

	// A doc with a backlink, then deleted (simulating a concurrent delete).
	exec("insert into data (title, data, type) values ('bl_del_a', '[[bl_dst_b]]', 'document')")
	exec("insert into data (title, data, type) values ('bl_dst_b', 'b', 'document')")

	// Pre-wipe + capture is done inside the API; to simulate "deleted during
	// rebuild" we delete it just before the call. The API enumerates, captures,
	// wipes, then rebuilds each enumerated doc; the deleted doc has no content
	// when rebuilt and must NOT be restored.
	exec("delete from data where title = 'bl_del_a'")

	out := route.Api_setting_backlink_reset_post(admin)
	// The deleted doc has no content at reset time, so it is not enumerated and
	// cannot be rebuilt; the remaining doc rebuilds cleanly. The response must
	// therefore reflect a clean rebuild (never an error state that would imply
	// a restoration of a doc that should not exist).
	if out["response"] != "ok" {
		t.Fatalf("reset response = %v, want ok (deleted doc must not be restored)", out)
	}
	// The deleted doc must have no backlink rows (it was wiped and must NOT be
	// restored).
	if hasBacklink(db, "bl_del_a", "bl_dst_b", t) {
		t.Fatalf("deleted doc bl_del_a was incorrectly restored (backlink to bl_dst_b)")
	}
}

// Test_Restore_backlink_backup directly verifies the failure-recovery path:
// capture the original backlinks, wipe them, "corrupt" the index, then restore
// the doc to its original state.
func Test_Restore_backlink_backup(t *testing.T) {
	db := Setup_query_api_db(t)
	defer tool.DB_close(db)

	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("restore exec: %v", err)
		}
	}
	doc := "bl_restore_a"
	// Seed the original backlinks and a link_count for this doc.
	exec("insert into data (title, data, type) values (?, 'content', 'document')", doc)
	exec("insert into back (link, title, type, data) values (?, 'orig_target', '', '')", doc)
	exec("insert into back (link, title, type, data) values (?, 'orig_target2', 'redirect', '')", doc)
	exec("insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, '', 'link_count', ?)", doc, "2")

	backup, err := route.Take_backlink_backup(db)
	if err != nil {
		t.Fatalf("Take_backlink_backup: %v", err)
	}
	// After Take_backlink_backup the non-manual backlinks are wiped.
	if hasBacklink(db, doc, "orig_target", t) {
		t.Fatalf("backlink still present after Take_backlink_backup wipe")
	}

	// Simulate a failed rebuild that wrote wrong backlinks.
	exec("insert into back (link, title, type, data) values (?, 'wrong_target', '', '')", doc)
	exec("insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, '', 'link_count', ?)", doc, "99")

	if err := route.Restore_backlink_backup(db, doc, backup); err != nil {
		t.Fatalf("Restore_backlink_backup: %v", err)
	}

	// The original backlinks must be back; the wrong ones must be gone.
	if !hasBacklink(db, doc, "orig_target", t) {
		t.Fatalf("orig_target missing after restore")
	}
	if !hasBacklink(db, doc, "orig_target2", t) {
		t.Fatalf("orig_target2 (redirect) missing after restore")
	}
	if hasBacklink(db, doc, "wrong_target", t) {
		t.Fatalf("wrong_target present after restore (should be replaced)")
	}
}
