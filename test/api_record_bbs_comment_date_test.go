package test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

// record_test_seed_comment_only seeds a legacy comment (comment_user_id row only,
// NO comment_date row) to model a date-missing legacy comment.
func record_test_seed_comment_only(db *sql.DB, t *testing.T, set_id, set_code, user string) {
	if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values ('comment_user_id', ?, ?, ?)`,
		set_code, set_id, user); err != nil {
		t.Fatalf("seed legacy comment_user_id: %v", err)
	}
}

// buggy_record_bbs_comment_dates replays the ORIGINAL pre-fix logic of
// Api_record_bbs_comment: group the user's comment rows by board, take min(set_id)
// as the representative post, then look up comment_date by that post's set_id ONLY
// (no user filter, no set_code). It returns board -> date so the pre-fix failure
// can be documented without reverting the fixed code.
func buggy_record_bbs_comment_dates(db *sql.DB, t *testing.T, user_name string) map[string]string {
	rows, err := db.Query(
		`select case when set_id like '-1-%' then '-1' else substr(set_id, 1, instr(set_id, '-') - 1) end as bbs_id,
			min(set_id) as set_id
		 from bbs_data
		 where set_name = "comment_user_id" and set_data = ?
		 group by bbs_id order by bbs_id + 0 desc limit 0, 50`,
		user_name,
	)
	if err != nil {
		t.Fatalf("buggy list: %v", err)
	}
	defer rows.Close()

	res := make(map[string]string)
	for rows.Next() {
		var bbs_id string
		var set_id string
		if err := rows.Scan(&bbs_id, &set_id); err != nil {
			t.Fatalf("buggy list scan: %v", err)
		}
		date := ""
		// Replicates the per-post subquery (QueryRow_DB leaves the var "" when no row).
		q := db.QueryRow(`select set_data from bbs_data where set_name = "comment_date" and set_id = ? order by set_data desc limit 1`, set_id)
		q.Scan(&date)
		res[bbs_id] = date
	}
	return res
}

// record_bbs_comment_get returns the board -> date map from the fixed API.
func record_bbs_comment_get(db *sql.DB, user_name, page string) map[string]string {
	res := route.Api_record_bbs_comment(tool.Config{IP: user_name}, user_name, page)
	if res["response"] != "ok" {
		panic("record_bbs_comment response is not ok: " + fmt.Sprint(res))
	}
	data, ok := res["data"].([][]string)
	if !ok {
		panic("record_bbs_comment data is not [][]string")
	}
	out := make(map[string]string)
	for _, d := range data {
		out[d[0]] = d[1]
	}
	return out
}

// sc_user_isolation proves a user's board date is their own latest comment, not a
// later comment by someone else on the same post (the pre-fix leak).
func sc_user_isolation(db *sql.DB, t *testing.T) {
	// alice: 5-1/code1 @ 2026-01-01 ; bob: 5-1/code2 @ 2026-03-03 (later).
	record_test_seed_comment(db, t, "5-1", "1", "2026-01-01 00:00:00", "alice")
	record_test_seed_comment(db, t, "5-1", "2", "2026-03-03 00:00:00", "bob")

	buggy := buggy_record_bbs_comment_dates(db, t, "alice")
	if buggy["5"] != "2026-03-03 00:00:00" {
		t.Fatalf("pre-fix user-isolation not reproduced: alice board got %q, want bob's 2026-03-03", buggy["5"])
	}
	t.Logf("pre-fix alice board date (bob leaked) = %q", buggy["5"])

	got := record_bbs_comment_get(db, "alice", "1")
	if got["5"] != "2026-01-01 00:00:00" {
		t.Fatalf("alice board date = %q, want her own 2026-01-01", got["5"])
	}

	bob := record_bbs_comment_get(db, "bob", "1")
	if bob["5"] != "2026-03-03 00:00:00" {
		t.Fatalf("bob board date = %q, want his own 2026-03-03", bob["5"])
	}
}

// sc_latest_per_board proves the board shows the user's LATEST comment across all
// posts in that board (not the min(set_id) post's date).
func sc_latest_per_board(db *sql.DB, t *testing.T) {
	// alice: 5-2/code1 @ 2026-06-06 (her latest), 5-10/code1 @ 2026-05-05.
	// bob: 5-10/code2 @ 2026-01-01 (earlier, same post as one of alice's).
	record_test_seed_comment(db, t, "5-2", "1", "2026-06-06 00:00:00", "alice")
	record_test_seed_comment(db, t, "5-10", "1", "2026-05-05 00:00:00", "alice")
	record_test_seed_comment(db, t, "5-10", "2", "2026-01-01 00:00:00", "bob")

	buggy := buggy_record_bbs_comment_dates(db, t, "alice")
	// min("5-2","5-10") lex = "5-10"; that post's max date (any user) = 2026-05-05.
	if buggy["5"] != "2026-05-05 00:00:00" {
		t.Fatalf("pre-fix latest-per-board not reproduced: alice board got %q, want 2026-05-05", buggy["5"])
	}
	t.Logf("pre-fix alice board date (min-post, wrong) = %q", buggy["5"])

	got := record_bbs_comment_get(db, "alice", "1")
	if got["5"] != "2026-06-06 00:00:00" {
		t.Fatalf("alice board date = %q, want her latest 2026-06-06", got["5"])
	}
}

// sc_legacy_missing_date covers the date-missing legacy comment: a board whose
// comments have no comment_date row still appears (with empty date), and a board
// with mixed dated/undated comments shows the max dated one.
func sc_legacy_missing_date(db *sql.DB, t *testing.T) {
	// Board 5: alice has ONLY an undated legacy comment -> board must still appear.
	record_test_seed_comment_only(db, t, "5-8", "1", "alice")

	got := record_bbs_comment_get(db, "alice", "1")
	if _, ok := got["5"]; !ok {
		t.Fatalf("legacy board '5' missing from result; legacy undated comments must still list the board")
	}
	if got["5"] != "" {
		t.Fatalf("legacy board '5' date = %q, want empty (no date row)", got["5"])
	}
	t.Logf("legacy board 5 present with empty date (OK)")

	// Board 6: mixed undated + dated -> max dated wins.
	record_test_seed_comment_only(db, t, "6-1", "1", "alice")
	record_test_seed_comment(db, t, "6-1", "2", "2026-02-02 00:00:00", "alice")
	got6 := record_bbs_comment_get(db, "alice", "1")
	if got6["6"] != "2026-02-02 00:00:00" {
		t.Fatalf("mixed board '6' date = %q, want 2026-02-02", got6["6"])
	}
}

// sc_order_special preserves board classification (numeric + -1 special), numeric
// ordering, and the [][]string{bbs_id, date} response form.
func sc_order_special(db *sql.DB, t *testing.T) {
	for sid, date := range map[string]string{
		"2-1":  "2026-01-01 00:00:00",
		"9-1":  "2026-01-02 00:00:00",
		"10-1": "2026-01-03 00:00:00",
		"11-1": "2026-01-04 00:00:00",
		"0-1":  "2026-01-05 00:00:00",
		"-1-1": "2026-01-06 00:00:00",
	} {
		record_test_seed_comment(db, t, sid, "1", date, "alice")
	}

	res := route.Api_record_bbs_comment(tool.Config{IP: "alice"}, "alice", "1")
	if res["response"] != "ok" {
		t.Fatalf("response = %v, want ok", res["response"])
	}
	data, ok := res["data"].([][]string)
	if !ok {
		t.Fatalf("data is not [][]string: %v", res["data"])
	}
	want := []string{"11", "10", "9", "2", "0", "-1"}
	if len(data) != len(want) {
		t.Fatalf("len = %d, want %d", len(data), len(want))
	}
	for i, v := range want {
		if len(data[i]) != 2 {
			t.Fatalf("row %d not [bbs_id, date]: %v", i, data[i])
		}
		if data[i][0] != v {
			t.Fatalf("board order: row %d = %q, want %q", i, data[i][0], v)
		}
	}
	// -1 special board keeps its date.
	if data[5][1] != "2026-01-06 00:00:00" {
		t.Fatalf("-1 board date = %q, want 2026-01-06", data[5][1])
	}
}

// sc_pagination preserves 50-item paging across the 50/51 boundary with 60 boards.
func sc_pagination(db *sql.DB, t *testing.T) {
	for i := 1; i <= 60; i++ {
		record_test_seed_comment(db, t, fmt.Sprintf("%d-1", i), "1", "2026-01-01 00:00:00", "alice")
	}
	page1 := record_bbs_comment_get(db, "alice", "1")
	page2 := record_bbs_comment_get(db, "alice", "2")
	if len(page1) != 50 {
		t.Fatalf("page1 = %d, want 50", len(page1))
	}
	if len(page2) != 10 {
		t.Fatalf("page2 = %d, want 10", len(page2))
	}
	// page1 = boards 60..11, page2 = boards 10..1.
	for _, want := range []string{"60", "11"} {
		if _, ok := page1[want]; !ok {
			t.Fatalf("page1 missing board %q", want)
		}
	}
	if d, ok := page1["10"]; ok {
		t.Fatalf("page1 incorrectly contains board 10; page1 date=%q", d)
	}
	if d, ok := page1["1"]; ok {
		t.Fatalf("page1 incorrectly contains board 1; page1 date=%q", d)
	}
	for _, want := range []string{"10", "1"} {
		if _, ok := page2[want]; !ok {
			t.Fatalf("page2 missing board %q", want)
		}
	}
	if d, ok := page2["60"]; ok {
		t.Fatalf("page2 incorrectly contains board 60; page2 date=%q", d)
	}
	// no overlap between pages
	for k := range page1 {
		if _, ok := page2[k]; ok {
			t.Fatalf("pagination overlap: board %q on both pages", k)
		}
	}
}

// Test_record_bbs_comment_date runs every scenario against every backend.
func Test_record_bbs_comment_date(t *testing.T) {
	backends := []string{"sqlite", "turso"}
	if os.Getenv("MYSQL_AUDIT_DSN") != "" {
		backends = append(backends, "mysql")
	}

	scenarios := []struct {
		name string
		fn   func(db *sql.DB, t *testing.T)
	}{
		{"user_isolation", sc_user_isolation},
		{"latest_per_board", sc_latest_per_board},
		{"legacy_missing_date", sc_legacy_missing_date},
		{"order_special", sc_order_special},
		{"pagination", sc_pagination},
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
