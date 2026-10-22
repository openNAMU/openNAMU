package test

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

// Setup_child_exists_db builds a schema and seeds docs to cover:
// - real children, orphaned children (no parent), and wildcard
// 	(%, _, !) false-positive cases.
func Setup_child_exists_db(t *testing.T) *sql.DB {
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
			t.Fatalf("seed %q: %v", q, err)
		}
	}

	// Basic: parent and child exist.
	exec("insert into data (title, data) values ('Basic', '')")
	exec("insert into data (title, data) values ('Basic/Child', '')")

	// Orphan: no parent doc, only a child.
	exec("insert into data (title, data) values ('Orphan/Child', '')")

	// Underscore: no A_st child, only a similar-name (AXt) child.
	exec("insert into data (title, data) values ('AXt/Child', '')")

	// Percent: no B%st child, only a similar-name (BAst) child.
	exec("insert into data (title, data) values ('BAst/Child', '')")

	// Real children whose names contain %, _, !.
	exec("insert into data (title, data) values ('C_st/RealChild', '')")
	exec("insert into data (title, data) values ('D%st/RealChild', '')")
	exec("insert into data (title, data) values ('E!st/RealChild', '')")

	return db
}

func Test_view_w_child_exists(t *testing.T) {
	db := Setup_child_exists_db(t)
	defer tool.DB_close(db)

	check := func(name, docName string, want bool) {
		got := route.View_w_child_exists(db, docName)
		if got != want {
			t.Errorf("%s: View_w_child_exists(%q) = %v, want %v", name, docName, got, want)
		}
	}

	// 1. Child exists -> menu shown.
	check("basic child exists", "Basic", true)

	// 2. Parent doc absent, only child -> menu still shown.
	check("no parent only child", "Orphan", true)

	// 3. No A_st child; only similar-name AXt/Child -> no menu (underscore escape).
	check("underscore no false positive", "A_st", false)

	// 4. No B%st child; only similar-name BAst/Child -> no menu (percent escape).
	check("percent no false positive", "B%st", false)

	// 5. Real children with %, _, ! in the name -> menu shown.
	check("underscore real child", "C_st", true)
	check("percent real child", "D%st", true)
	check("bang real child", "E!st", true)
}
