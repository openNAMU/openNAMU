package test

import (
	"database/sql"
	"os"
	"strings"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

// sc_date_correctness verifies Api_record_bbs_comment_in returns the COMMENT's
// own comment_date (data[4]) instead of the parent post's date, and keeps the
// previous 4-field layout by appending the date as the 5th field.
func sc_date_correctness(db *sql.DB, t *testing.T) {
	// Parent post 5/3 with its own, older date.
	post_date := "2020-01-01 00:00:00"
	if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values ('user_id', ?, ?, ?)`,
		"3", "5", "admin_ip"); err != nil {
		t.Fatalf("seed post user_id: %v", err)
	}
	if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values ('title', ?, ?, ?)`,
		"3", "5", "Post Title"); err != nil {
		t.Fatalf("seed post title: %v", err)
	}
	if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values ('date', ?, ?, ?)`,
		"3", "5", post_date); err != nil {
		t.Fatalf("seed post date: %v", err)
	}

	// Comment 5-3/code1 with a NEWER comment_date (different from the post date).
	comment_date := "2026-10-15 12:00:00"
	record_test_seed_comment(db, t, "5-3", "1", comment_date, "admin_ip")

	res := route.Api_record_bbs_comment_in(tool.Config{IP: "admin_ip"}, "admin_ip", "5", "1")
	if res["response"] != "ok" {
		t.Fatalf("response = %v, want ok", res["response"])
	}
	got, ok := res["data"].([][]string)
	if !ok {
		t.Fatalf("data is not [][]string: %v", res["data"])
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}

	// Response shape: first 4 fields intact, date appended as the 5th.
	if got[0][0] != "5" {
		t.Fatalf("bbs_id = %v, want 5", got[0][0])
	}
	if got[0][1] != "3" {
		t.Fatalf("post code = %v, want 3", got[0][1])
	}
	if got[0][2] != "5-3" {
		t.Fatalf("set_id = %v, want 5-3", got[0][2])
	}
	if got[0][3] != "1" {
		t.Fatalf("set_code = %v, want 1", got[0][3])
	}
	// The comment's own date, not the post's (which is different).
	if got[0][4] != comment_date {
		t.Fatalf("comment date = %v, want %s (post date is %s)", got[0][4], comment_date, post_date)
	}
}

// sc_comment_in_legacy_missing_date ensures a legacy comment that only has a
// comment_user_id row (no comment_date row) is still returned, with an empty
// date, instead of being dropped by the join. It also documents the pre-fix
// INNER JOIN failure (legacy comment absent).
func sc_comment_in_legacy_missing_date(db *sql.DB, t *testing.T) {
	// A dated comment.
	record_test_seed_comment(db, t, "5-7", "1", "2026-10-10 00:00:00", "admin_ip")
	// A legacy comment: only the comment_user_id row, no comment_date.
	record_test_seed_comment_only(db, t, "5-3", "2", "admin_ip")

	// Buggy (pre-fix) INNER JOIN probe: the legacy 5-3/2 has no comment_date row,
	// so it is dropped; only the dated 5-7/1 remains.
	var buggy []string
	brows, berr := db.Query(
		`select c.set_id
		 from bbs_data c
		 inner join bbs_data d on d.set_name = "comment_date" and d.set_id = c.set_id and d.set_code = c.set_code
		 where c.set_name = "comment_user_id" and c.set_id like '5-%' and c.set_data = 'admin_ip'
		 order by c.set_id desc limit 0, 50`)
	defer brows.Close()
	if berr != nil {
		t.Fatalf("buggy probe: %v", berr)
	}
	for brows.Next() {
		var sid string
		if err := brows.Scan(&sid); err != nil {
			t.Fatalf("buggy probe scan: %v", err)
		}
		buggy = append(buggy, sid)
	}
	if len(buggy) != 1 || buggy[0] != "5-7" {
		t.Errorf("pre-fix INNER JOIN probe = %v, want [5-7] (legacy 5-3/2 dropped)", buggy)
	}

	res := route.Api_record_bbs_comment_in(tool.Config{IP: "admin_ip"}, "admin_ip", "5", "1")
	if res["response"] != "ok" {
		t.Fatalf("response = %v, want ok", res["response"])
	}
	got, ok := res["data"].([][]string)
	if !ok {
		t.Fatalf("data is not [][]string: %v", res["data"])
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2 (dated + legacy)", len(got))
	}

	// Dated comment first (date desc), legacy (empty date) last.
	if got[0][2] != "5-7" || got[0][3] != "1" {
		t.Fatalf("first = set_id %v set_code %v, want 5-7/1", got[0][2], got[0][3])
	}
	if got[0][4] != "2026-10-10 00:00:00" {
		t.Fatalf("first date = %v, want 2026-10-10 00:00:00", got[0][4])
	}
	if got[1][2] != "5-3" || got[1][3] != "2" {
		t.Fatalf("second = set_id %v set_code %v, want 5-3/2 (legacy)", got[1][2], got[1][3])
	}
	if got[1][4] != "" {
		t.Fatalf("legacy date = %v, want empty", got[1][4])
	}
}

// sc_anchor_helper verifies the anchor comment_code produced by
// Bbs_comment_storage_location for general, nested, and -1 comments. This is
// the value the View must put after "#" so the link jumps to the comment (the
// comment page renders each comment with id="<code_id>", and the legacy record
// view used the same comment_code for its anchor).
func sc_anchor_helper(db *sql.DB, t *testing.T) {
	cases := []struct {
		name       string
		set_id     string
		set_code   string
		wantCode   string
		wantValid  bool
	}{
		{"general", "5-3", "1", "1", true},
		{"nested", "5-3-1", "2", "1-2", true},
		{"board_-1", "-1-2025", "1", "1", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, _, comment_code, valid := route.Bbs_comment_storage_location(c.set_id, c.set_code)
			if comment_code != c.wantCode || valid != c.wantValid {
				t.Fatalf("Bbs_comment_storage_location(%s, %s) = code %q valid %v, want code %q valid %v",
					c.set_id, c.set_code, comment_code, valid, c.wantCode, c.wantValid)
			}
		})
	}
}

// sc_view_render renders the actual record comment page and verifies the
// post-fix output: the comment's own comment_date is shown (not the parent post
// date), and the link jumps to the comment via a #comment_code anchor (the same
// value the comment page uses as each comment's id).
func sc_view_render(db *sql.DB, t *testing.T) {
	post_date := "2020-01-01 00:00:00"
	comment_date := "2026-10-15 12:00:00"
	if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values ('user_id', ?, ?, ?)`,
		"3", "5", "admin_ip"); err != nil {
		t.Fatalf("seed post user_id: %v", err)
	}
	if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values ('title', ?, ?, ?)`,
		"3", "5", "Post Title"); err != nil {
		t.Fatalf("seed post title: %v", err)
	}
	if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values ('date', ?, ?, ?)`,
		"3", "5", post_date); err != nil {
		t.Fatalf("seed post date: %v", err)
	}
	// Alice's comment, top-level under post 5/3 (set_id 5-3, code 1).
	record_test_seed_comment(db, t, "5-3", "1", comment_date, "alice")

	out := route.View_record_bbs_comment_in(
		tool.Config{IP: "admin_ip", Session: &emptySession{}},
		"alice", "5", "1")
	if out == "" {
		t.Fatal("view output is empty")
	}
	// The comment's own date is displayed.
	if !strings.Contains(out, comment_date) {
		t.Fatalf("view does not display comment date %s", comment_date)
	}
	// The parent post date must NOT be shown in the comment list.
	if strings.Contains(out, post_date) {
		t.Fatalf("view incorrectly shows the parent post date %s instead of the comment date", post_date)
	}
	// The link must jump to the comment (top-level comment code "1").
	if !strings.Contains(out, `href="/bbs/w/5/3#1"`) {
		t.Fatalf("view link missing comment anchor #1; output snippet: %s", out)
	}
	t.Logf("view rendered: comment date present, post date absent, anchor /bbs/w/5/3#1 present")
}

func Test_record_bbs_comment_view(t *testing.T) {
	backends := []string{"sqlite", "turso"}
	if os.Getenv("MYSQL_AUDIT_DSN") != "" {
		backends = append(backends, "mysql")
	}

	scenarios := []struct {
		name string
		fn   func(db *sql.DB, t *testing.T)
	}{
		{"date_correctness", sc_date_correctness},
		{"comment_in_legacy_missing_date", sc_comment_in_legacy_missing_date},
		{"anchor_helper", sc_anchor_helper},
		{"view_render", sc_view_render},
	}

	for _, b := range backends {
		for _, sc := range scenarios {
			t.Run(sc.name+"_"+b, func(t *testing.T) {
				db := new_record_db(t, b)
				sc.fn(db, t)
			})
		}
	}
}
