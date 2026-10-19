package test

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

// Setup_rebuild_db builds a full schema and seeds a few wiki docs with external
// links plus an admin user bound to the setting_backlink permission.
func Setup_rebuild_db(t *testing.T) *sql.DB {
	t.Helper()

	db_dir := filepath.Join(t.TempDir(), "db")
	os.MkdirAll(db_dir, 0755)
	os.Setenv("NAMU_DB", db_dir)
	os.Setenv("NAMU_DB_TYPE", "sqlite")
	defer os.Unsetenv("NAMU_DB")
	defer os.Unsetenv("NAMU_DB_TYPE")

	tool.DB_init()
	db := tool.DB_connect()

	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed exec %v: %v", q, err)
		}
	}
	exec("insert into data (title, data) values ('alist','[[page|external page]]\n[[nonpublic]]')")
	exec("insert into data (title, data) values ('page','[[page]]')")
	exec("insert into data (title, data) values ('Page','[[page|same host]]')")
	exec("insert into data (title, data) values ('Example','[[https://example.com/doc|external example]]')")
	exec("insert into data (title, data) values ('example','[[https://example.com/doc|external example]]')")
	// Doc whose external backlink must survive the rebuild.
	exec("insert into data (title, data) values ('NonPublic','[[https://nonpub.example.com/doc|non public]]')")
	exec("insert into data (title, data) values ('Beta','target beta content')")
	exec("insert into data (title, data) values ('admin_ip','user admin')")
	// DB_init creates the tables but not the alist entries (Always_init does in
	// production); seed the owner group so admin_ip gets the full admin perms.
	exec("insert into alist (name, acl) values ('owner','owner')")
	exec("insert into alist (name, acl) values ('admin','admin')")
	// admin_ip holds the owner group, which grants setting_backlink and
	// document edit on all docs.
	exec("insert into user_set (name, id, data) values ('acl','admin_ip','owner')")

	return db
}

func Count_backlinks(db *sql.DB, link, title string) int {
	var c int
	if err := db.QueryRow("select count(*) from back where link = ? and title = ? and type = 'external'", link, title).Scan(&c); err != nil {
		return 0
	}
	return c
}

// Test_Api_w_render_backlink_uses_saved_source verifies that /api/render/backlink
// rebuilds the index from the saved source, not the unsaved content.
func Test_Api_w_render_backlink_uses_saved_source(t *testing.T) {
	db := Setup_rebuild_db(t)
	defer tool.DB_close(db)

	config := tool.Config{IP: "admin_ip"}

	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed exec %v: %v", q, err)
		}
	}
	// Saved content has only the internal [[Beta]] link (Beta exists). The
	// request carries an unsaved external link [[Secret]] which must NOT be
	// recorded in the index.
	exec("insert into data (title, data) values ('RenderTarget','[[Beta]]')")

	out := route.Api_w_render(config, "RenderTarget", "[[Secret]]", "backlink", "")

	if resp, ok := out["response"].(string); !ok || resp != "ok" {
		t.Fatalf("render response not ok: %v", out)
	}
	if c := Count_backlinks(db, "RenderTarget", "secret.example.com"); c != 0 {
		t.Fatalf("unsaved external link was recorded (count=%d)", c)
	}
	var ic int
	if err := db.QueryRow("select count(*) from back where link = 'RenderTarget' and title = 'Beta' and type = ''").Scan(&ic); err != nil || ic == 0 {
		t.Fatalf("saved internal backlink missing: %v", err)
	}
}

// Test_Setting_backlink_reset_includes_non_public_doc verifies that the admin
// rebuild processes all docs (bypassing per-doc ACL) and maintains each doc's
// external backlinks, including a doc the rebuild admin would not edit.
func Test_Setting_backlink_reset_includes_non_public_doc(t *testing.T) {
	db := Setup_rebuild_db(t)
	defer tool.DB_close(db)

	config := tool.Config{IP: "admin_ip"}
	out := route.Api_setting_backlink_reset_post_internal(config, "normal")

	if resp, ok := out["response"].(string); !ok || resp != "ok" {
		t.Fatalf("rebuild response not ok: %v", out)
	}
	doc_count, _ := out["document_count"].(int)
	total_count, _ := out["document_total"].(int)
	if doc_count != total_count {
		t.Fatalf("rebuild processed %d of %d docs", doc_count, total_count)
	}
	// Non-Example doc's external backlink is maintained.
	if c := Count_backlinks(db, "NonPublic", "nonpub.example.com"); c != 1 {
		t.Fatalf("non-public doc backlink missing after rebuild (count=%d), want 1", c)
	}
	// Normal doc's external backlink is also maintained.
	if c := Count_backlinks(db, "Example", "example.com"); c != 1 {
		t.Fatalf("normal doc backlink missing after rebuild (count=%d), want 1", c)
	}
}

// Test_Setting_backlink_reset_no_duplicate_external verifies that rebuilding a
// doc whose content references the same external domain through multiple URLs
// (and rebuilding it again) stores exactly one (doc, domain, external) row.
func Test_Setting_backlink_reset_no_duplicate_external(t *testing.T) {
	db := Setup_rebuild_db(t)
	defer tool.DB_close(db)
	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed exec %v: %v", q, err)
		}
	}
	content := "[[https://example.com/a]]\n[[https://example.com/b#x?y=1]]\n[[https://sub.example.com/c]]\n[[Beta]]"
	exec("insert into data (title, data) values ('DupDoc', ?)", content)
	exec("insert into data_set (doc_name, doc_rev, set_name, set_data) values ('DupDoc', '', 'document_markup', 'namumark')")

	admin := tool.Config{IP: "admin_ip"}
	out := route.Api_setting_backlink_reset_post_internal(admin, "normal")
	if resp, ok := out["response"].(string); !ok || resp != "ok" {
		t.Fatalf("rebuild response not ok: %v", out)
	}
	if c := Count_backlinks(db, "DupDoc", "example.com"); c != 1 {
		t.Fatalf("② duplicate external backlink (DupDoc, example.com): count=%d, want 1", c)
	}
	if c := Count_backlinks(db, "DupDoc", "sub.example.com"); c != 1 {
		t.Fatalf("② duplicate external backlink (DupDoc, sub.example.com): count=%d, want 1", c)
	}

	out2 := route.Api_setting_backlink_reset_post_internal(admin, "normal")
	if resp2, ok2 := out2["response"].(string); !ok2 || resp2 != "ok" {
		t.Fatalf("2nd rebuild response not ok: %v", out2)
	}
	if c := Count_backlinks(db, "DupDoc", "example.com"); c != 1 {
		t.Fatalf("② duplicate external backlink after 2nd rebuild (DupDoc, example.com): count=%d, want 1", c)
	}
}

// Test_Setting_backlink_reset_stores_only_domain verifies the rebuild stores
// only the normalized domain in the external backlink title, never a URL that
// still contains a path, query string, anchor, or fragment.
func Test_Setting_backlink_reset_stores_only_domain(t *testing.T) {
	db := Setup_rebuild_db(t)
	defer tool.DB_close(db)
	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed exec %v: %v", q, err)
		}
	}
	content := "[[https://example.com/Page#anchor?a=1&b=2]]\n[[http://sub.example.co.uk/path/with/slashes]]"
	exec("insert into data (title, data) values ('DomDoc', ?)", content)
	exec("insert into data_set (doc_name, doc_rev, set_name, set_data) values ('DomDoc', '', 'document_markup', 'namumark')")

	admin := tool.Config{IP: "admin_ip"}
	out := route.Api_setting_backlink_reset_post_internal(admin, "normal")
	if resp, ok := out["response"].(string); !ok || resp != "ok" {
		t.Fatalf("rebuild response not ok: %v", out)
	}
	if c := Count_backlinks(db, "DomDoc", "example.com"); c != 1 {
		t.Fatalf("③ domain example.com missing after rebuild (count=%d), want 1", c)
	}
	if c := Count_backlinks(db, "DomDoc", "sub.example.co.uk"); c != 1 {
		t.Fatalf("③ domain sub.example.co.uk missing after rebuild (count=%d), want 1", c)
	}
	var bad int
	if err := db.QueryRow("select count(*) from back where link='DomDoc' and type='external' and (title LIKE '%%/%' OR title LIKE '%%?%%' OR title LIKE '%%#%%')").Scan(&bad); err != nil {
		t.Fatalf("scan domain-leak count: %v", err)
	}
	if bad != 0 {
		t.Fatalf("③ stored external title contains path/query/anchor: %d rows leak", bad)
	}
}

// Test_Setting_backlink_reset_non_public_acl verifies the rebuild processes a
// doc that is non-public via ACL (owner-only view permission) and maintains its
// external backlink.
func Test_Setting_backlink_reset_non_public_acl(t *testing.T) {
	db := Setup_rebuild_db(t)
	defer tool.DB_close(db)
	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed exec %v: %v", q, err)
		}
	}
	exec("insert into data (title, data) values ('PrivateDoc', ?)", "[[https://priv.example.com/doc]]")
	exec("insert into acl (title, data, type) values ('PrivateDoc', 'owner', 'view')")

	admin := tool.Config{IP: "admin_ip"}
	out := route.Api_setting_backlink_reset_post_internal(admin, "normal")
	if resp, ok := out["response"].(string); !ok || resp != "ok" {
		t.Fatalf("rebuild response not ok: %v", out)
	}
	total, _ := out["document_total"].(int)
	processed, _ := out["document_count"].(int)
	if processed != total {
		t.Fatalf("④ rebuild processed %d of %d docs (must process every doc, ignoring ACL)", processed, total)
	}
	if c := Count_backlinks(db, "PrivateDoc", "priv.example.com"); c != 1 {
		t.Fatalf("④ non-public (ACL) doc external backlink missing after rebuild (count=%d), want 1", c)
	}
}

// Test_Setting_backlink_reset_empty_db_wipes_stale_auto verifies that a rebuild
// on a 0-doc wiki returns "ok" and wipes orphaned auto backlinks (nothing left
// to rebuild), while preserving manual backlinks.
func Test_Setting_backlink_reset_empty_db_wipes_stale_auto(t *testing.T) {
	db := Setup_empty_db(t)
	defer tool.DB_close(db)

	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed: %v", err)
		}
	}
	// Orphaned auto backlink (its source doc no longer exists) plus a manual
	// backlink that must survive a 0-doc rebuild.
	exec("insert into back (link, title, type, data) values ('GoneDoc','gone.example.com','external','')")
	exec("insert into back (link, title, type, data) values ('GoneDoc','manual.example.com','cat_manual','')")

	admin := tool.Config{IP: "admin_ip"}
	out := route.Api_setting_backlink_reset_post_internal(admin, "normal")
	if out["response"] != "ok" {
		t.Fatalf("0-doc rebuild response = %v, want ok", out["response"])
	}
	if count, ok := out["document_count"].(int); !ok || count != 0 {
		t.Fatalf("0-doc rebuild document_count = %v, want 0", out["document_count"])
	}
	if total, ok := out["document_total"].(int); !ok || total != 0 {
		t.Fatalf("0-doc rebuild document_total = %v, want 0", out["document_total"])
	}
	// Orphaned auto backlink must be wiped; manual backlink preserved.
	var auto int
	if err := db.QueryRow("select count(*) from back where link='GoneDoc' and type='external'").Scan(&auto); err != nil || auto != 0 {
		t.Fatalf("orphaned auto backlink survived a 0-doc rebuild (count=%d), want 0", auto)
	}
	var manual int
	if err := db.QueryRow("select count(*) from back where link='GoneDoc' and type='cat_manual'").Scan(&manual); err != nil || manual != 1 {
		t.Fatalf("manual backlink not preserved after a 0-doc rebuild (count=%d), want 1", manual)
	}
}

// Setup_empty_db builds a schema-only DB (no data rows) for the empty-rebuild
// verification.
func Setup_empty_db(t *testing.T) *sql.DB {
	t.Helper()

	db_dir := filepath.Join(t.TempDir(), "db")
	if err := os.MkdirAll(db_dir, 0755); err != nil {
		t.Fatalf("mkdir db dir: %v", err)
	}
	os.Setenv("NAMU_DB", db_dir)
	os.Setenv("NAMU_DB_TYPE", "sqlite")
	defer os.Unsetenv("NAMU_DB")
	defer os.Unsetenv("NAMU_DB_TYPE")

	tool.DB_init()
	db := tool.DB_connect()

	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed: %v", err)
		}
	}
	// Grant admin_ip the owner group so setting_backlink is available.
	exec("insert into alist (name, acl) values ('owner','owner')")
	exec("insert into user_set (name, id, data) values ('acl','admin_ip','owner')")

	return db
}

// Test_Setting_backlink_reset_empty_db verifies a rebuild on an empty wiki
// returns "ok" with zero counts and does not crash.
func Test_Setting_backlink_reset_empty_db(t *testing.T) {
	db := Setup_empty_db(t)
	defer tool.DB_close(db)

	admin := tool.Config{IP: "admin_ip"}
	out := route.Api_setting_backlink_reset_post_internal(admin, "normal")
	if out["response"] != "ok" {
		t.Fatalf("empty DB rebuild response = %v, want ok", out["response"])
	}
	if count, ok := out["document_count"].(int); !ok || count != 0 {
		t.Fatalf("empty DB rebuild document_count = %v, want 0", out["document_count"])
	}
	if err, ok := out["error_count"].(int); !ok || err != 0 {
		t.Fatalf("empty DB rebuild error_count = %v, want 0", out["error_count"])
	}
	// document_total is 0; the 0-doc path now wipes stale auto backlinks and
	// reports counts rather than short-circuiting before them.
	if total, ok := out["document_total"].(int); !ok || total != 0 {
		t.Fatalf("empty DB rebuild document_total = %v, want 0", out["document_total"])
	}
}

// Test_Setting_backlink_reset_success_replaces_originals verifies a fully
// successful rebuild replaces pre-existing (stale) non-manual backlinks and
// link_count with ones rebuilt from the saved source, and reports "ok".
func Test_Setting_backlink_reset_success_replaces_originals(t *testing.T) {
	db := Setup_rebuild_db(t)
	defer tool.DB_close(db)

	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed: %v", err)
		}
	}
	// Stale original backlink + link_count that the rebuild must wipe.
	exec("insert into back (link, title, type, data) values ('Example','stale.example.com','external','')")
	exec("insert into data_set (doc_name, doc_rev, set_name, set_data) values ('Example','','link_count','7')")

	admin := tool.Config{IP: "admin_ip"}
	out := route.Api_setting_backlink_reset_post_internal(admin, "normal")
	if out["response"] != "ok" {
		t.Fatalf("rebuild response = %v, want ok", out["response"])
	}
	total, _ := out["document_total"].(int)
	count, _ := out["document_count"].(int)
	if count != total {
		t.Fatalf("rebuild processed %d of %d docs (all should succeed)", count, total)
	}
	// Example's saved content references example.com -> rebuilt; stale original wiped.
	if c := Count_backlinks(db, "Example", "example.com"); c != 1 {
		t.Fatalf("Example rebuilt backlink missing (count=%d), want 1", c)
	}
	if c := Count_backlinks(db, "Example", "stale.example.com"); c != 0 {
		t.Fatalf("stale original backlink survived rebuild (count=%d), want 0", c)
	}
	// NonPublic doc's external backlink maintained too.
	if c := Count_backlinks(db, "NonPublic", "nonpub.example.com"); c != 1 {
		t.Fatalf("NonPublic backlink missing after rebuild (count=%d), want 1", c)
	}
}

// Test_Setting_backlink_reset_multi_page verifies the rebuild processes docs
// spanning more than one 50-row page with no page-boundary omission.
func Test_Setting_backlink_reset_multi_page(t *testing.T) {
	db := Setup_rebuild_db(t)
	defer tool.DB_close(db)

	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed: %v", err)
		}
	}
	// Seed 60 wiki docs, each with a distinct external link, to force multiple
	// title-list pages.
	for i := 0; i < 60; i++ {
		doc := fmt.Sprintf("PageDoc%04d", i)
		exec("insert into data (title, data) values (?, ?)", doc, fmt.Sprintf("[[https://host%04d.example.com/doc]]", i))
	}

	admin := tool.Config{IP: "admin_ip"}
	out := route.Api_setting_backlink_reset_post_internal(admin, "normal")
	if out["response"] != "ok" {
		t.Fatalf("multi-page rebuild response = %v, want ok", out["response"])
	}
	total, _ := out["document_total"].(int)
	count, _ := out["document_count"].(int)
	if count != total {
		t.Fatalf("multi-page rebuild processed %d of %d docs (want all)", count, total)
	}
	// A sample from a later page was rebuilt.
	if c := Count_backlinks(db, "PageDoc0059", "host0059.example.com"); c != 1 {
		t.Fatalf("late-page doc PageDoc0059 backlink missing (count=%d), want 1", c)
	}
	if c := Count_backlinks(db, "PageDoc0000", "host0000.example.com"); c != 1 {
		t.Fatalf("first-page doc PageDoc0000 backlink missing (count=%d), want 1", c)
	}
}

// Test_Setting_backlink_reset_large verifies a rebuild over a large doc set
// (well over 1000) completes without an out-of-memory or page-boundary failure
// and processes every doc.
func Test_Setting_backlink_reset_large(t *testing.T) {
	db := Setup_empty_db(t)
	defer tool.DB_close(db)

	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed: %v", err)
		}
	}
	const total = 1500
	for i := 0; i < total; i++ {
		doc := fmt.Sprintf("BigDoc%05d", i)
		exec("insert into data (title, data) values (?, ?)", doc, "[[https://big.example.com/doc]]")
	}

	admin := tool.Config{IP: "admin_ip"}
	out := route.Api_setting_backlink_reset_post_internal(admin, "normal")
	if out["response"] != "ok" {
		t.Fatalf("large rebuild response = %v, want ok", out["response"])
	}
	totalRows, _ := out["document_total"].(int)
	count, _ := out["document_count"].(int)
	if totalRows != total {
		t.Fatalf("large rebuild document_total = %d, want %d", totalRows, total)
	}
	if count != totalRows {
		t.Fatalf("large rebuild processed %d of %d docs (want all)", count, totalRows)
	}
	// Spot-check docs across the set.
	if c := Count_backlinks(db, "BigDoc00000", "big.example.com"); c != 1 {
		t.Fatalf("BigDoc00000 backlink missing (count=%d), want 1", c)
	}
	if c := Count_backlinks(db, "BigDoc01499", "big.example.com"); c != 1 {
		t.Fatalf("BigDoc01499 backlink missing (count=%d), want 1", c)
	}
}
