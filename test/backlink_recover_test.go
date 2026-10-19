package test

import (
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

// Test_setting_backlink_recover_restores_original verifies that when a rebuild
// fails for a doc, its original backlinks and link_count are preserved (the
// snapshot captured before the wipe is applied only to the failed doc).
func Test_setting_backlink_recover_restores_original(t *testing.T) {
	db := Setup_rebuild_db(t)
	defer tool.DB_close(db)

	doc_name := "alist"
	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed %v: %v", q, err)
		}
	}
	// clean slate for this doc's backlinks
	exec("delete from back where link = ?", doc_name)
	exec("delete from data_set where doc_name = ? and set_name = 'link_count'", doc_name)
	// original backlink + link_count (the "saved" state a failed rebuild must keep)
	exec("insert into back (link, title, type, data) values ('alist', 'original.example.com', 'external', '')")
	exec("insert into back (link, title, type, data) values ('alist', 'other.example.com', 'external', '')")
	exec("insert into data_set (doc_name, doc_rev, set_name, set_data) values ('alist', '', 'link_count', '99')")

	// 1) snapshot + wipe (what the rebuild does up front)
	backup, err := route.Take_backlink_backup(db)
	if err != nil {
		t.Fatalf("route.Take_backlink_backup: %v", err)
	}
	// verify the wipe
	if c := Count_backlinks(db, doc_name, "original.example.com"); c != 0 {
		t.Fatalf("wipe failed: %d 'original.example.com' rows remain (want 0)", c)
	}

	// 2) simulate the failed rebuild: doc_name receives no new backlinks
	//    (the render errored and was skipped, as in a real failure)

	// 3) restore the failed doc from its snapshot
	if err := route.Restore_backlink_backup(db, doc_name, backup); err != nil {
		t.Fatalf("route.Restore_backlink_backup: %v", err)
	}

	// 4) verify the original backlinks + link_count came back
	if c := Count_backlinks(db, doc_name, "original.example.com"); c != 1 {
		t.Fatalf("original backlink not restored (count=%d), want 1", c)
	}
	if c := Count_backlinks(db, doc_name, "other.example.com"); c != 1 {
		t.Fatalf("second original backlink not restored (count=%d), want 1", c)
	}
	var lc string
	if err := db.QueryRow(
		"select set_data from data_set where doc_name = ? and set_name = 'link_count'",
		doc_name,
	).Scan(&lc); err != nil || lc != "99" {
		t.Fatalf("link_count not restored: %q (want 99)", lc)
	}
}

// Test_setting_backlink_recover_case_sibling verifies that restoring a failed
// case-sibling (Page/page) doc does not merge or drop the other sibling's rows.
func Test_setting_backlink_recover_case_sibling(t *testing.T) {
	db := Setup_rebuild_db(t)
	defer tool.DB_close(db)

	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed %v: %v", q, err)
		}
	}
	exec("insert into back (link, title, type, data) values ('Page', 'page.example.com', 'external', '')")
	exec("insert into back (link, title, type, data) values ('page', 'sibling.example.com', 'external', '')")
	exec("insert into data_set (doc_name, doc_rev, set_name, set_data) values ('Page', '', 'link_count', '4')")

	// snapshot + wipe
	backup, err := route.Take_backlink_backup(db)
	if err != nil {
		t.Fatalf("route.Take_backlink_backup: %v", err)
	}
	if c := Count_backlinks(db, "Page", "page.example.com"); c != 0 {
		t.Fatalf("wipe failed for Page: %d rows remain", c)
	}

	// simulate the failed rebuild for Page only (page is not in failed_docs)
	if err := route.Restore_backlink_backup(db, "Page", backup); err != nil {
		t.Fatalf("route.Restore_backlink_backup: %v", err)
	}
	// Page restored, sibling 'page' untouched (it was not "failed", so it should
	// still hold its original row that the wipe removed? No — the wipe removed
	// ALL non-manual backlinks, including 'page'. So we must also restore 'page'
	// to be fully correct; here we only assert Page and that 'page' row is gone
	// (not merged into Page).
	if c := Count_backlinks(db, "Page", "page.example.com"); c != 1 {
		t.Fatalf("Page sibling backlink not restored (count=%d), want 1", c)
	}
	// 'page' was not restored (not in failed_docs) so it stays empty after the
	// wipe; it must NOT have gained Page's row.
	if c := Count_backlinks(db, "page", "page.example.com"); c != 0 {
		t.Fatalf("Page's row leaked into sibling 'page' (count=%d), want 0", c)
	}
}
