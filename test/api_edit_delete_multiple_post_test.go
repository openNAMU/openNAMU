package test

import (
	"database/sql"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

// Test_api_edit_delete_multiple_post_preserves_original_name is a regression
// test for the bug where Api_edit_delete_multiple_post applied
// strings.TrimSpace to each document name before deletion. That rewrote a
// document name that differs only by surrounding whitespace into the identity
// of a coexisting look-alike, so the wrong document was deleted.
//
// The fix must pass each line's original (untrimmed) text to
// Api_edit_delete_post, while still skipping blank lines (empty or
// whitespace-only). The scenario seeds coexisting pairs that differ only by a
// leading space and requests only the space-bearing names for deletion; the
// space-less names must survive.
//
// Leading-space variants are used (rather than trailing) because SQLite/Turso
// (BINARY collation) and MySQL both treat a leading space as significant,
// whereas MySQL's collation collapses trailing spaces (a trailing-space pair
// is not addressable distinctly on MySQL).
func Test_api_edit_delete_multiple_post_preserves_original_name(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso", "mysql"} {
		backend := backend
		t.Run(backend, func(t *testing.T) {
			db := setup_delete_db(t, backend)
			defer tool.DB_close(db)

			exec := func(q string, args ...any) {
				if _, err := db.Exec(q, args...); err != nil {
					t.Fatalf("seed %v: %v", q, err)
				}
			}
			// Coexisting docs that differ only by a leading space.
			exec("insert into data (title, data) values ('DocA','DocA body')")
			exec("insert into data (title, data) values (' DocA','DocA leading body')")
			exec("insert into data (title, data) values ('Bar','Bar body')")
			exec("insert into data (title, data) values (' Bar','Bar leading body')")

			config := tool.Config{IP: "admin_ip"}

			// Only the leading-space names are requested. An empty line and a
			// whitespace-only line are present and must be skipped.
			content := " DocA\n\n   \n Bar"
			out := route.Api_edit_delete_multiple_post(config, content, "ok", "yes")
			response, _ := out["response"].(string)
			if response != "ok" {
				t.Fatalf("multiple delete response %q, want ok; full %v", response, out)
			}

			if doc_exists(db, " DocA", t) {
				t.Errorf("' DocA' should be deleted but still exists")
			}
			if !doc_exists(db, "DocA", t) {
				t.Errorf("'DocA' was deleted by mistake; only ' DocA' (leading space) was requested (TrimSpace bug)")
			}
			if doc_exists(db, " Bar", t) {
				t.Errorf("' Bar' should be deleted but still exists")
			}
			if !doc_exists(db, "Bar", t) {
				t.Errorf("'Bar' was deleted by mistake; only ' Bar' (leading space) was requested (TrimSpace bug)")
			}
		})
	}
}

// setup_delete_db boots a fresh DB for the named backend with the owner admin
// seeded, ready for document deletion.
func setup_delete_db(t *testing.T, backend string) *sql.DB {
	t.Helper()
	switch backend {
	case "sqlite":
		db := audit_setup(t, "sqlite")
		audit_seed_auth(db, t)
		return db
	case "turso":
		if !tool.Turso_platform_supported() {
			t.Skip("turso unsupported platform")
		}
		db := audit_setup(t, "turso")
		audit_seed_auth(db, t)
		return db
	case "mysql":
		// mysql_audit_setup seeds the owner admin itself and uses a real
		// MariaDB from MYSQL_AUDIT_DSN (skips when unset).
		db, _ := mysql_audit_setup(t)
		return db
	default:
		t.Fatalf("unknown backend: %q", backend)
		return nil
	}
}

// doc_exists reports whether a document with the exact title exists.
func doc_exists(db *sql.DB, title string, t *testing.T) bool {
	var c int
	if err := db.QueryRow("select count(*) from data where title = ?", title).Scan(&c); err != nil {
		t.Fatalf("doc_exists %q: %v", title, err)
	}
	return c > 0
}
