package route

import (
	"database/sql"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func create_test_report_db(t *testing.T) *sql.DB {
	db_path := filepath.Join(t.TempDir(), "report.db")
	db, err := sql.Open("sqlite", db_path)
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}
	if _, err := db.Exec(`CREATE TABLE bbs_data (set_name TEXT, set_code TEXT, set_id TEXT, set_data TEXT)`); err != nil {
		db.Close()
		t.Fatalf("create bbs_data table: %v", err)
	}
	return db
}

func Test_report_create_success(t *testing.T) {
	db := create_test_report_db(t)
	defer db.Close()

	set_code_new, err := Report_create(db, "127.0.0.1", "Document title", "/w/document", "Report", "spam reason text")
	if err != nil {
		t.Fatalf("Report_create success: %v", err)
	}
	if set_code_new == "" {
		t.Fatalf("Report_create returned empty set_code")
	}

	var stored_code string
	err = db.QueryRow(
		`select set_code from bbs_data where set_name = 'title' and set_id = '-2' order by set_code + 0 desc limit 1`,
	).Scan(&stored_code)
	if err != nil {
		t.Fatalf("read stored report title: %v", err)
	}
	if stored_code != set_code_new {
		t.Errorf("stored report set_code = %q, expected %q", stored_code, set_code_new)
	}

	var title string
	if err := db.QueryRow(`select set_data from bbs_data where set_name = 'title' and set_id = '-2' and set_code = ?`, set_code_new).Scan(&title); err != nil {
		t.Fatalf("read stored report title value: %v", err)
	}
	if title != "Document title" {
		t.Errorf("stored report title = %q, expected %q", title, "Document title")
	}
}

func Test_report_create_error_no_panic(t *testing.T) {
	db := create_test_report_db(t)
	db.Close()

	_, err := Report_create(db, "127.0.0.1", "Document title", "/w/document", "Report", "spam reason text")
	if err == nil {
		t.Errorf("Report_create on a closed db should return an error, got nil")
	}
}
