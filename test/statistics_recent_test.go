package test

import (
	"database/sql"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"opennamu/route"
	"opennamu/route/tool"
)

// Setup_statistics_db builds the full schema and seeds history rows spanning
// the last 24h/7d/30d so the recent edit-count items can be verified.
func Setup_statistics_db(t *testing.T) *sql.DB {
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

	// default permission groups so view helpers that touch acl resolve cleanly.
	exec("insert into alist (name, acl) values ('owner','owner')")
	exec("insert into alist (name, acl) values ('admin','admin')")
	exec("insert into alist (name, acl) values ('user','user')")
	exec("insert into alist (name, acl) values ('ip','ip')")

	// Force ko-KR so the rendered labels are the known Korean strings.
	exec("insert into other (name, data, coverage) values ('language','ko-KR','')")

	now := time.Now()
	fmtDate := func(d time.Time) string { return d.Format("2006-01-02 15:04:05") }
	insertHistory := func(date string) {
		exec(`insert into history (id,title,data,date,ip,send,leng,hide,type) values (?,?,?,?,?,?,?,?,?)`,
			"stat_test_doc", "stat doc", "content", date, "1.2.3.4", "edit", "+1", "", "")
	}
	// 1 row in 24h (and also in 7d and 30d)
	insertHistory(fmtDate(now.Add(-2 * time.Hour)))
	// 1 row in 7d but not 24h
	insertHistory(fmtDate(now.Add(-72 * time.Hour)))
	// 1 row in 30d but not 7d
	insertHistory(fmtDate(now.Add(-360 * time.Hour)))
	// 1 row older than 30d (excluded from all three)
	insertHistory(fmtDate(now.Add(-1440 * time.Hour)))

	return db
}

// count_for_label returns the number shown next to the translated label in the
// rendered page. It anchors on the label + ":" so it is robust to the surrounding
// template markup (e.g. whether <li> is escaped).
func count_for_label(html string, db *sql.DB, key string, t *testing.T) int {
	label := tool.Get_language(db, key, true)
	re := regexp.MustCompile(label + `[[:space:]]*:[[:space:]]*([0-9]+)`)
	m := re.FindStringSubmatch(html)
	if m == nil {
		if err := os.WriteFile("/home/spdv/dev/openNAMU/tmp/statistics_output.html", []byte(html), 0644); err == nil {
			t.Logf("dumped rendered page to ./tmp/statistics_output.html")
		}
		t.Fatalf("could not find %q in rendered statistics page", label)
	}
	got, err := strconv.Atoi(m[1])
	if err != nil {
		t.Fatalf("could not parse count %q for %q", m[1], key)
	}
	return got
}

func Test_view_statistics_recent_edit_counts(t *testing.T) {
	db := Setup_statistics_db(t)
	defer tool.DB_close(db)

	config := tool.Config{IP: "1.2.3.4", Session: &emptySession{}}
	html := route.View_statistics(config)

	// New recent edit-count items: exact counts (seeded 1/2/3).
	if got := count_for_label(html, db, "statistics_recent_24h_edit_count", t); got != 1 {
		t.Errorf("recent 24h edit count = %d, want 1", got)
	}
	if got := count_for_label(html, db, "statistics_recent_7d_edit_count", t); got != 2 {
		t.Errorf("recent 7d edit count = %d, want 2", got)
	}
	if got := count_for_label(html, db, "statistics_recent_30d_edit_count", t); got != 3 {
		t.Errorf("recent 30d edit count = %d, want 3", got)
	}

	// Existing items are still rendered (regression guard).
	existing_keys := []string{
		"statistics_document_count",
		"statistics_user_count",
		"statistics_edit_count",
		"statistics_bbs_post_count",
		"statistics_bbs_comment_count",
		"statistics_monthly_edit_count",
		"statistics_monthly_top_contributor",
	}
	for _, key := range existing_keys {
		label := tool.Get_language(db, key, true)
		if !strings.Contains(html, label) {
			t.Errorf("missing existing statistics item: %q", label)
		}
	}
}
