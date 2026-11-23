package test

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
	"opennamu/route/tool"
)

// Create_sitemap_db builds a minimal SQLite schema containing only the tables
// Make_sitemap touches: `other` (settings) and `data` (document titles).
func Create_sitemap_db(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "sitemap.db"))
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	if _, err := db.Exec(`CREATE TABLE other (name TEXT, data TEXT, coverage TEXT)`); err != nil {
		t.Fatalf("create other table: %v", err)
	}
	if _, err := db.Exec(`CREATE TABLE data (test TEXT, title TEXT, data TEXT, type TEXT)`); err != nil {
		t.Fatalf("create data table: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO data (title, data, type) VALUES ('테스트문서', 'x', '')`); err != nil {
		t.Fatalf("insert document: %v", err)
	}

	return db
}

// Regression: with the "도메인 제외" option off and no domain setting,
// Make_sitemap used to write sitemap files with broken loc URLs
// ("http:///w/..."). It must now refuse with Err_sitemap_domain_require and
// leave no files behind, while the explicit exclude option keeps working.
func Test_make_sitemap_domain_require(t *testing.T) {
	// Make_sitemap writes sitemap files into the working directory.
	t.Chdir(t.TempDir())
	db := Create_sitemap_db(t)

	// no domain setting, exclude option off → refuse, no files written
	_, _, err := tool.Make_sitemap(db)
	if !errors.Is(err, tool.Err_sitemap_domain_require) {
		t.Fatalf("Make_sitemap err = %v, want Err_sitemap_domain_require", err)
	}
	if _, stat_err := os.Stat("sitemap.xml"); stat_err == nil {
		t.Errorf("sitemap.xml must not be created when the domain setting is missing")
	}

	// the explicit "도메인 제외" option keeps working without a domain setting
	if _, err := db.Exec(`INSERT INTO other (name, data, coverage) VALUES ('sitemap_auto_exclude_domain', 'on', '')`); err != nil {
		t.Fatalf("insert exclude setting: %v", err)
	}
	if _, _, err := tool.Make_sitemap(db); err != nil {
		t.Fatalf("Make_sitemap with exclude option: %v", err)
	}
	sitemap_zero, err := os.ReadFile("sitemap_0.xml")
	if err != nil {
		t.Fatalf("read sitemap_0.xml: %v", err)
	}
	if !strings.Contains(string(sitemap_zero), `<loc>/w/`) {
		t.Errorf("sitemap_0.xml with exclude option should use relative loc, got: %s", sitemap_zero)
	}
	if strings.Contains(string(sitemap_zero), `http:///`) {
		t.Errorf("sitemap_0.xml must not contain broken host URLs: %s", sitemap_zero)
	}

	// domain setting present → absolute URLs
	if _, err := db.Exec(`UPDATE other SET data = '' WHERE name = 'sitemap_auto_exclude_domain'`); err != nil {
		t.Fatalf("clear exclude setting: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO other (name, data, coverage) VALUES ('domain', 'wiki.example.org', '')`); err != nil {
		t.Fatalf("insert domain setting: %v", err)
	}
	if _, _, err := tool.Make_sitemap(db); err != nil {
		t.Fatalf("Make_sitemap with domain setting: %v", err)
	}
	sitemap_zero, err = os.ReadFile("sitemap_0.xml")
	if err != nil {
		t.Fatalf("read sitemap_0.xml: %v", err)
	}
	if !strings.Contains(string(sitemap_zero), `<loc>http://wiki.example.org/w/`) {
		t.Errorf("sitemap_0.xml should use the domain setting, got: %s", sitemap_zero)
	}
	index_data, err := os.ReadFile("sitemap.xml")
	if err != nil {
		t.Fatalf("read sitemap.xml: %v", err)
	}
	if !strings.Contains(string(index_data), `<loc>http://wiki.example.org/sitemap_0.xml</loc>`) {
		t.Errorf("sitemap.xml index should use the domain setting, got: %s", index_data)
	}
}
