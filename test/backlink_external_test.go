package test

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

// Setup_external_backlink_db builds a full schema and seeds two documents that
// differ only by case (Page / page) plus the internal and external backlinks
// needed to verify that Stage 1 SQL changes never corrupt a case-sibling doc's
// rows.
func Setup_external_backlink_db(t *testing.T) *sql.DB {
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

	exec := func(query string, args ...any) {
		if _, err := db.Exec(query, args...); err != nil {
			t.Fatalf("seed: %v", err)
		}
	}

	// Real documents.
	exec("insert into data (title, data) values ('Page','content Page')")
	exec("insert into data (title, data) values ('page','content page')")
	exec("insert into data (title, data) values ('A','content A')")
	exec("insert into data (title, data) values ('B','content B')")
	exec("insert into data (title, data) values ('C','content C')")
	exec("insert into data (title, data) values ('D','content D')")

	// Internal backlinks:
	//   A -> Page   (title=Page,  link=A,  type=cat)
	//   Page -> C   (title=C,     link=Page, type=cat)
	//   B -> page   (title=page,  link=B,  type=cat)
	//   page -> D   (title=D,     link=page, type=cat)
	exec("insert into back (title, link, type, data) values ('Page','A','cat','')")
	exec("insert into back (title, link, type, data) values ('C','Page','cat','')")
	exec("insert into back (title, link, type, data) values ('page','B','cat','')")
	exec("insert into back (title, link, type, data) values ('D','page','cat','')")
	// External backlinks: title = normalized hostname, link = source doc,
	// type = external. The two case-siblings use distinct hosts so their rows
	// are unambiguous while the SQL must still target by source doc name.
	exec("insert into back (title, link, type, data) values ('example.com','Page','external','')")
	exec("insert into back (title, link, type, data) values ('page.example.com','page','external','')")

	return db
}

func Count_rows(db *sql.DB, t *testing.T, query string, args ...any) int {
	var c int
	if err := db.QueryRow(query, args...).Scan(&c); err != nil {
		t.Errorf("count query failed: %v", err)
	}
	return c
}

// Test_Move_document_protects_external verifies that a REAL document move
// (data rename + backlink/history re-attribution, i.e. the full doc move, not
// just a naming change) re-attributes only Page's rows and leaves the case-sibling
// page's rows untouched; the external URL title is preserved while the source doc
// name is updated.
func Test_Move_document_protects_external(t *testing.T) {
	db := Setup_external_backlink_db(t)
	defer tool.DB_close(db)

	route.Move_document_normal_data(db, "Page", "Page2")

	if Count_rows(db, t, "select count(*) from back where title='Page2' and link='A' and type='cat'") != 1 {
		t.Fatalf("A's internal backlink to Page2 missing")
	}
	if Count_rows(db, t, "select count(*) from back where title='C' and link='Page2' and type='cat'") != 1 {
		t.Fatalf("Page's internal link to C should now be sourced from Page2")
	}
	if Count_rows(db, t, "select count(*) from back where title='example.com' and link='Page2' and type='external'") != 1 {
		t.Fatalf("Page's external link should be re-sourced from Page2 with URL title preserved")
	}
	// The data row itself must be renamed (a full move renames the document).
	var renamed int
	if err := db.QueryRow("select count(*) from data where "+tool.Exact_binary_eq("title"), "Page2").Scan(&renamed); err != nil || renamed != 1 {
		t.Fatalf("Page data row not renamed to Page2 (count=%d)", renamed)
	}

	// page's rows must be unchanged.
	if c := Count_rows(db, t, "select count(*) from back where title='page' and link='B' and type='cat'"); c != 1 {
		t.Fatalf("page's internal backlink from B corrupted (count=%d)", c)
	}
	if c := Count_rows(db, t, "select count(*) from back where title='D' and link='page' and type='cat'"); c != 1 {
		t.Fatalf("page's internal link to D corrupted (count=%d)", c)
	}
	if c := Count_rows(db, t, "select count(*) from back where title='page.example.com' and link='page' and type='external'"); c != 1 {
		t.Fatalf("page's external link corrupted (count=%d)", c)
	}
	// Old internal title 'Page' should no longer exist.
	if c := Count_rows(db, t, "select count(*) from back where title='Page' and type='cat'"); c != 0 {
		t.Fatalf("internal title 'Page' not renamed to 'Page2' (count=%d)", c)
	}
}

// Test_Move_document_rotate_protects_external verifies the 3-way rotate only
// re-attributes Page's rows and never touches page's.
func Test_Move_document_rotate_protects_external(t *testing.T) {
	db := Setup_external_backlink_db(t)
	defer tool.DB_close(db)

	exec := func(query string, args ...any) {
		if _, err := db.Exec(query, args...); err != nil {
			t.Fatalf("seed: %v", err)
		}
	}
	exec("insert into data (title, data) values ('tmp_rotate','content tmp')")

	route.Move_document_rotate(db, "tmp_rotate", "Page", "Page2")

	if Count_rows(db, t, "select count(*) from back where title='Page2' and link='A' and type='cat'") != 1 {
		t.Fatalf("A's internal backlink to Page2 missing after rotate")
	}
	if Count_rows(db, t, "select count(*) from back where title='example.com' and link='Page2' and type='external'") != 1 {
		t.Fatalf("Page's external not re-sourced to Page2 (URL title must remain)")
	}
	if Count_rows(db, t, "select count(*) from back where title='page' and link='B' and type='cat'") != 1 {
		t.Fatalf("page's internal backlink from B corrupted by rotate")
	}
	if Count_rows(db, t, "select count(*) from back where title='page.example.com' and link='page' and type='external'") != 1 {
		t.Fatalf("page's external corrupted by rotate")
	}
}

// Test_Move_document_merge_protects_external verifies merging Page into
// target deletes only Page's backlinks and leaves page's untouched.
func Test_Move_document_merge_protects_external(t *testing.T) {
	db := Setup_external_backlink_db(t)
	defer tool.DB_close(db)

	exec := func(query string, args ...any) {
		if _, err := db.Exec(query, args...); err != nil {
			t.Fatalf("seed: %v", err)
		}
	}
	exec("insert into data (title, data) values ('target','content target')")

	route.Move_document_merge(db, 0, "Page", "target", nil)

	if Count_rows(db, t, "select count(*) from back where title='page' and link='B' and type='cat'") != 1 {
		t.Fatalf("page's internal backlink from B deleted by merge of Page")
	}
	if Count_rows(db, t, "select count(*) from back where title='D' and link='page' and type='cat'") != 1 {
		t.Fatalf("page's internal link to D deleted by merge of Page")
	}
	if Count_rows(db, t, "select count(*) from back where title='page.example.com' and link='page' and type='external'") != 1 {
		t.Fatalf("page's external deleted by merge of Page")
	}
	// Page's external now sourced from target, URL preserved.
	if Count_rows(db, t, "select count(*) from back where title='example.com' and link='target' and type='external'") != 1 {
		t.Fatalf("Page's external should be re-sourced from target with URL title preserved")
	}
}

// Test_Xref_excludes_external verifies the existing general xref queries
// (route.Api_w_xref) return only internal backlinks, never external rows.
func Test_Xref_excludes_external(t *testing.T) {
	db := Setup_external_backlink_db(t)
	defer tool.DB_close(db)

	// do_type=1 : documents linking TO the target (where title = doc_name).
	rows1, err := db.Query(
		"select distinct link, type from back where title = ? and not type = 'no' and not type = 'nothing' and not type = 'external' order by type asc, link asc",
		"Page",
	)
	if err != nil {
		t.Fatalf("xref query: %v", err)
	}
	n1 := 0
	var name, typ string
	for rows1.Next() {
		if err := rows1.Scan(&name, &typ); err != nil {
			rows1.Close()
			t.Fatalf("scan: %v", err)
		}
		n1++
		t.Logf("xref1: link=%s type=%s", name, typ)
		if typ == "external" {
			rows1.Close()
			t.Fatalf("external row leaked into do_type=1 xref: link=%s type=%s", name, typ)
		}
	}
	rows1.Close()
	if n1 != 1 {
		t.Fatalf("do_type=1 xref for Page should return exactly 1 internal row (A), got %d", n1)
	}

	// do_type=2 : documents linked FROM the target (where link = doc_name).
	rows2, err := db.Query(
		"select distinct title, type from back where link = ? and not type = 'no' and not type = 'nothing' and not type = 'external' order by type asc, title asc",
		"Page",
	)
	if err != nil {
		t.Fatalf("xref query: %v", err)
	}
	n2 := 0
	for rows2.Next() {
		if err := rows2.Scan(&name, &typ); err != nil {
			rows2.Close()
			t.Fatalf("scan: %v", err)
		}
		n2++
		t.Logf("xref2: title=%s type=%s", name, typ)
		if typ == "external" {
			rows2.Close()
			t.Fatalf("external row leaked into do_type=2 xref: title=%s type=%s", name, typ)
		}
	}
	rows2.Close()
	if n2 != 1 {
		t.Fatalf("do_type=2 xref for Page should return exactly 1 internal row (C), got %d", n2)
	}
}

// Test_Delete_preserves_other_document_external verifies the delete SQL
// removes only the deleted document's backlinks (internal + external) and
// leaves every other document's backlinks untouched.
func Test_Delete_preserves_other_document_external(t *testing.T) {
	db := Setup_external_backlink_db(t)
	defer tool.DB_close(db)

	exec := func(query string, args ...any) {
		if _, err := db.Exec(query, args...); err != nil {
			t.Fatalf("delete: %v", err)
		}
	}
	exec("delete from back where link = 'Page' and type != 'external'")
	exec("delete from back where link = 'Page' and type = 'external'")

	if Count_rows(db, t, "select count(*) from back where title='page' and link='B' and type='cat'") != 1 {
		t.Fatalf("page's internal backlink from B deleted by Page's delete")
	}
	if Count_rows(db, t, "select count(*) from back where title='D' and link='page' and type='cat'") != 1 {
		t.Fatalf("page's internal link to D deleted by Page's delete")
	}
	if Count_rows(db, t, "select count(*) from back where title='page.example.com' and link='page' and type='external'") != 1 {
		t.Fatalf("page's external deleted by Page's delete")
	}
	// Page's own backlinks must be gone.
	if c := Count_rows(db, t, "select count(*) from back where title='example.com' and type='external'"); c != 0 {
		t.Fatalf("Page's external link should have been deleted (count=%d)", c)
	}
	if c := Count_rows(db, t, "select count(*) from back where title='C' and link='Page' and type='cat'"); c != 0 {
		t.Fatalf("Page's internal link to C should have been deleted (count=%d)", c)
	}
}
