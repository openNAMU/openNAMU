package test

import (
	"database/sql"
	"opennamu/route/tool/markup"
	"os"
	"path/filepath"
	"testing"

	"opennamu/route/tool"
)

// Setup_namumark_db builds a full schema and returns a connected sqlite DB.
func Setup_namumark_db(t *testing.T) *sql.DB {
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
	return tool.DB_connect()
}

func Count_rows(db *sql.DB, t *testing.T, query string, args ...any) int {
	var c int
	if err := db.QueryRow(query, args...).Scan(&c); err != nil {
		t.Errorf("count query failed: %v", err)
	}
	return c
}

// Test_Namumark_external_backlink_collect verifies that a namumark document
// containing an inline external link produces an external backlink row with the
// normalized host title, and that link_count is NOT incremented for it.
func Test_Namumark_external_backlink_collect(t *testing.T) {
	db := Setup_namumark_db(t)
	defer tool.DB_close(db)

	doc_name := "ExtLinkDoc"
	content := "intro line\n[[https://Example.Com/Page|label text]]\nfinal"

	if _, err := db.Exec("insert into data (title, data) values (?, ?)", doc_name, content); err != nil {
		t.Fatalf("seed: %v", err)
	}

	markup.Render_namumark_compat_backlink(db, doc_name, content)

	if Count_rows(db, t, "select count(*) from back where title=? and link=? and type='external'", "example.com", doc_name) != 1 {
		t.Fatalf("external backlink not persisted for %s", doc_name)
	}
	if c := Count_rows(db, t, "select count(*) from back where link=? and type='external'", doc_name); c != 1 {
		t.Fatalf("expected exactly 1 external backlink, got %d", c)
	}

	var lc int
	if err := db.QueryRow("select set_data from data_set where doc_name=? and set_name='link_count'", doc_name).Scan(&lc); err == nil && lc != 0 {
		t.Fatalf("link_count should be 0 for an external-only document, got %d", lc)
	}
}

// Test_Namumark_external_link_with_anchor_collect verifies that a namumark
// external link with query and fragment normalizes to just the host (the
// host is the only thing stored in V5).
func Test_Namumark_external_link_with_anchor_collect(t *testing.T) {
	db := Setup_namumark_db(t)
	defer tool.DB_close(db)

	doc_name := "AnchorDoc"
	content := "[[https://example.com/a?b=1#c]]\n"

	if _, err := db.Exec("insert into data (title, data) values (?, ?)", doc_name, content); err != nil {
		t.Fatalf("seed: %v", err)
	}

	markup.Render_namumark_compat_backlink(db, doc_name, content)

	if Count_rows(db, t, "select count(*) from back where title=? and link=? and type='external'", "example.com", doc_name) != 1 {
		t.Fatalf("external backlink with host-normalized title not persisted for %s", doc_name)
	}
}

// Test_Namumark_external_redirect_collect verifies that a namumark document
// whose first line redirects to an external URL produces an external backlink
// row (the redirect branch calls Add_backlink before the collect_only
// early-return). The host is normalized (case lowercased) and stored.
func Test_Namumark_external_redirect_collect(t *testing.T) {
	db := Setup_namumark_db(t)
	defer tool.DB_close(db)

	doc_name := "RedirectDoc"
	content := "#redirect https://Example.Com/Article#frag\ntail"

	if _, err := db.Exec("insert into data (title, data) values (?, ?)", doc_name, content); err != nil {
		t.Fatalf("seed: %v", err)
	}

	markup.Render_namumark_compat_backlink(db, doc_name, content)

	if Count_rows(db, t, "select count(*) from back where title=? and link=? and type='external'", "example.com", doc_name) != 1 {
		t.Fatalf("external redirect backlink not persisted for %s", doc_name)
	}
}

// Test_Namumark_external_and_internal_counts verifies that a document with both
// internal and external links has both backlinks collected, and link_count
// reflects only the internal link.
func Test_Namumark_external_and_internal_counts(t *testing.T) {
	db := Setup_namumark_db(t)
	defer tool.DB_close(db)

	if _, err := db.Exec("insert into data (title, data) values ('RealTarget','target content')"); err != nil {
		t.Fatalf("seed: %v", err)
	}
	doc_name := "MixedDoc"
	content := "intro\n[[RealTarget]]\n[[https://example.com/Other]]\n"
	if _, err := db.Exec("insert into data (title, data) values (?, ?)", doc_name, content); err != nil {
		t.Fatalf("seed: %v", err)
	}

	markup.Render_namumark_compat_backlink(db, doc_name, content)

	if Count_rows(db, t, "select count(*) from back where title=? and link=? and type=''", "RealTarget", doc_name) != 1 {
		t.Fatalf("internal backlink to RealTarget not persisted")
	}
	if Count_rows(db, t, "select count(*) from back where title=? and link=? and type='external'", "example.com", doc_name) != 1 {
		t.Fatalf("external backlink to example.com not persisted")
	}

	var lc int
	if err := db.QueryRow("select set_data from data_set where doc_name=? and set_name='link_count'", doc_name).Scan(&lc); err != nil {
		t.Fatalf("link_count query failed: %v", err)
	}
	if lc != 1 {
		t.Fatalf("link_count should be 1 (only the internal link), got %d", lc)
	}
}

// Test_Merge_child_skips_external verifies that Merge_child merges the child's
// internal backlink into the parent but skips the child's external backlink.
// (Merge_child is only invoked from full-render include/inter paths, so this is
// tested directly rather than through the backlink-collection path.)
func Test_Merge_child_skips_external(t *testing.T) {
	db := Setup_namumark_db(t)
	defer tool.DB_close(db)

	parent := markup.New_namumark_compat_renderer(db, "Parent", "", "view", nil, 0, false)
	child := markup.New_namumark_compat_renderer(db, "Child", "", "include", nil, 0, false)

	child.Add_backlink("Internal", "", "")
	child.Add_backlink("https://child.example.com/e", "external", "")

	parent.Merge_child(child)

	if !parent.Backlink_exists("Internal", "") {
		t.Fatalf("child's internal backlink not merged into Parent")
	}
	if parent.Backlink_exists("https://child.example.com/e", "external") {
		t.Fatalf("child's external backlink was merged into Parent")
	}
	if parent.Has_external_backlink() {
		t.Fatalf("an external backlink was merged into Parent")
	}
}
