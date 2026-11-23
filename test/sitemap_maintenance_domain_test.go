package test

import (
	"os"
	"path/filepath"
	"testing"

	"opennamu/route/tool"
)

// Regression: after the domain requirement was introduced, a missing domain
// setting made Maintenance_daily return an error. The scheduler treats that as
// a transient failure and retries every few minutes without ever advancing the
// daily last-run state, so every other daily task kept re-running.
// A missing domain must be a silent skip instead, keeping the daily cadence.
func Test_Maintenance_sitemap_auto_skips_without_domain(t *testing.T) {
	t.Chdir(t.TempDir())
	db := new_test_db(t, filepath.Join(t.TempDir(), "sitemap-maintenance.db"))
	defer tool.DB_close(db)

	if _, err := db.Exec(`insert into data (title, data, type) values ('Doc', 'content', 'd')`); err != nil {
		t.Fatalf("insert document: %v", err)
	}

	set_test_setting(db, t, "sitemap_auto_make", "1")
	set_test_setting(db, t, "ua_expiration_date", "")
	set_test_setting(db, t, "auth_history_expiration_date", "")

	if err := tool.Maintenance_daily(db); err != nil {
		t.Fatalf("daily maintenance must not fail on a missing domain setting: %v", err)
	}
	if _, err := os.Stat("sitemap.xml"); err == nil {
		t.Errorf("sitemap.xml must not be created when the domain setting is missing")
	}

	// with a domain configured the same setup must still generate the files
	if _, err := db.Exec(`INSERT INTO other (name, data, coverage) VALUES ('domain', 'wiki.example.org', '')`); err != nil {
		t.Fatalf("insert domain setting: %v", err)
	}
	if err := tool.Maintenance_daily(db); err != nil {
		t.Fatalf("daily maintenance with a domain setting: %v", err)
	}
	if _, err := os.Stat("sitemap.xml"); err != nil {
		t.Fatalf("sitemap.xml missing after maintenance with a domain setting: %v", err)
	}
	if _, err := os.Stat("sitemap_0.xml"); err != nil {
		t.Fatalf("sitemap_0.xml missing after maintenance with a domain setting: %v", err)
	}
}
