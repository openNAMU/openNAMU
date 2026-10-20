package test

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

// histRev describes one history row: id (revision number), content, hide
// ("" visible, "O" hidden).
type histRev struct {
	id   string
	data string
 hide string
}

// Setup_history_tool_db builds a DB with a regular group (view + history_view,
// no hidel) and a trusted group (admin, full access including hidel).
func Setup_history_tool_db(t *testing.T) *sql.DB {
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
			t.Fatalf("seed: %v", err)
		}
	}
	// regular: view history but cannot unhide.
	exec("insert into alist (name, acl) values ('regular_group','view')")
	exec("insert into alist (name, acl) values ('regular_group','history_view')")
	// trusted: admin => everything incl. hidel.
	exec("insert into alist (name, acl) values ('trusted_group','admin')")
	exec("insert into user_set (id, name, data) values ('regular_user','acl','regular_group')")
	exec("insert into user_set (id, name, data) values ('trusted_user','acl','trusted_group')")

	return db
}

// seed_revisions inserts a doc and its history rows (in id order).
func seed_revisions(db *sql.DB, t *testing.T, doc string, revs []histRev) {
	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed: %v", err)
		}
	}
	exec("insert into data (title, data) values (?, '')", doc)
	for i, r := range revs {
		date := fmt.Sprintf("2026-01-%02d 00:00:00", i+1)
		exec(`insert into history (id, title, data, date, ip, send, leng, hide, type)
			values (?,?,?,?,?,?,?,?,?)`, r.id, doc, r.data, date, "test_ip", "edit", "+1", r.hide, "")
	}
}

// viewTool returns the rendered history tool page.
func viewTool(db *sql.DB, user, doc, rev string) string {
	session := &emptySession{}
	return route.View_history_tool(tool.Config{IP: user, Session: session}, doc, rev)
}

// wantPrev asserts the previous-revision diff and revert links point to
// previous (for current=revision, doc).
func wantPrev(t *testing.T, out string, previous, revision, doc string) {
	diff := "/diff/" + previous + "/" + revision + "/" + doc + "\""
	if !strings.Contains(out, diff) {
		t.Fatalf("want diff link %q, got: %q", diff, out[:500])
	}
	revert := "/revert/" + previous + "/" + doc + "\""
	if !strings.Contains(out, revert) {
		t.Fatalf("want revert link %q, got: %q", revert, out[:500])
	}
}

// wantNoPrev asserts no previous-revision links are generated (no /diff/ link).
func wantNoPrev(t *testing.T, out string, doc string) {
	if strings.Contains(out, "/diff/") {
		t.Fatalf("want no previous/diff links for %s, got: %q", doc, out[:500])
	}
	// The revert-to-current link is always present; ensure the only /revert/
	// is the revert-to-current one (no /revert/0/ or other previous).
	if strings.Contains(out, "/revert/0/") {
		t.Fatalf("found /revert/0/ for %s (revision 1 has no previous)", doc)
	}
}

func Test_view_history_tool_previous_revision(t *testing.T) {
	db := Setup_history_tool_db(t)
	defer tool.DB_close(db)

	// 1) consecutive, all visible: r3's previous is r2.
	seed_revisions(db, t, "PrevDoc", []histRev{
		{"1", "a", ""},
		{"2", "b", ""},
		{"3", "c", ""},
	})
	out := viewTool(db, "regular_user", "PrevDoc", "3")
	wantPrev(t, out, "2", "3", "PrevDoc")

	// 2) gap (r2 deleted): r3's previous must be r1, not the missing r2.
	seed_revisions(db, t, "GapDoc", []histRev{
		{"1", "a", ""},
		{"3", "c", ""},
	})
	out = viewTool(db, "regular_user", "GapDoc", "3")
	wantPrev(t, out, "1", "3", "GapDoc")

	// 3) only one revision: no previous links.
	seed_revisions(db, t, "FirstDoc", []histRev{
		{"1", "a", ""},
	})
	out = viewTool(db, "regular_user", "FirstDoc", "1")
	wantNoPrev(t, out, "FirstDoc")

	// 3b) all smaller revisions deleted (only r5 exists): the old code would
	// point to a nonexistent r4; now no previous links are generated.
	seed_revisions(db, t, "OnlyRev5", []histRev{
		{"5", "e", ""},
	})
	out = viewTool(db, "regular_user", "OnlyRev5", "5")
	wantNoPrev(t, out, "OnlyRev5")

	// 4) previous hidden, regular user (no hidel): previous excluded.
	seed_revisions(db, t, "HideDoc", []histRev{
		{"1", "a", "O"},
		{"3", "c", ""},
	})
	out = viewTool(db, "regular_user", "HideDoc", "3")
	wantNoPrev(t, out, "HideDoc")

	// 5) previous hidden, trusted user (has hidel): previous shown.
	out = viewTool(db, "trusted_user", "HideDoc", "3")
	wantPrev(t, out, "1", "3", "HideDoc")

	// 6) mixed: r1 visible, r2 hidden, r3 visible.
	seed_revisions(db, t, "MixDoc", []histRev{
		{"1", "a", ""},
		{"2", "b", "O"},
		{"3", "c", ""},
	})
	out = viewTool(db, "regular_user", "MixDoc", "3")
	// regular: closest non-hidden smaller than 3 is r1 (skips hidden r2).
	wantPrev(t, out, "1", "3", "MixDoc")
	out = viewTool(db, "trusted_user", "MixDoc", "3")
	// trusted: closest smaller (incl. hidden) is r2.
	wantPrev(t, out, "2", "3", "MixDoc")
}
