package test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

// new_record_db returns a fresh, fully-initialized DB with an owner auth for the
// given backend ("sqlite"|"turso"|"mysql"). MySQL requires MYSQL_AUDIT_DSN and
// skips otherwise (the caller must filter it out).
func new_record_db(t *testing.T, backend string) *sql.DB {
	switch backend {
	case "mysql":
		db, _ := mysql_audit_setup(t)
		return db
	default:
		db := audit_setup(t, backend)
		audit_seed_auth(db, t)
		return db
	}
}

// record_test_seed_user_id seeds one post (user_id row).
func record_test_seed_user_id(db *sql.DB, t *testing.T, set_code, set_id, set_data string) {
	if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`,
		"user_id", set_code, set_id, set_data); err != nil {
		t.Fatalf("seed user_id: %v", err)
	}
}

// record_test_seed_comment seeds one comment (comment_user_id + comment_date).
func record_test_seed_comment(db *sql.DB, t *testing.T, set_id, set_code, date, user string) {
	if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values ('comment_user_id', ?, ?, ?)`,
		set_code, set_id, user); err != nil {
		t.Fatalf("seed comment_user_id: %v", err)
	}
	if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values ('comment_date', ?, ?, ?)`,
		set_code, set_id, date); err != nil {
		t.Fatalf("seed comment_date: %v", err)
	}
}

// sc_bbs_in_numeric verifies the per-board post record is sorted by post number
// (numeric) descending, not lexicographically.
func sc_bbs_in_numeric(db *sql.DB, t *testing.T) {
	for _, code := range []string{"2", "9", "10", "11"} {
		record_test_seed_user_id(db, t, code, "5", "admin_ip")
	}

	// Reproduce the pre-fix bug directly to document the lexicographic failure.
	var buggy []string
	brows, berr := db.Query(`select set_code from bbs_data where set_name = 'user_id' and set_id = '5' and set_data = 'admin_ip' order by set_code desc limit 0, 50`)
	defer brows.Close()
	if berr != nil {
		t.Fatalf("buggy probe: %v", berr)
	}
	for brows.Next() {
		var c string
		if err := brows.Scan(&c); err != nil {
			t.Fatalf("buggy probe scan: %v", err)
		}
		buggy = append(buggy, c)
	}
	if len(buggy) == 0 || buggy[0] != "9" {
		t.Fatalf("pre-fix lex order not reproduced; got %v", buggy)
	}

	res := route.Api_record_bbs_in(tool.Config{IP: "admin_ip"}, "admin_ip", "5", "1")
	if res["response"] != "ok" {
		t.Fatalf("response = %v, want ok", res["response"])
	}
	got, ok := res["data"].([]string)
	if !ok {
		t.Fatalf("data is not []string: %v", res["data"])
	}
	want := []string{"11", "10", "9", "2"}
	for i, v := range got {
		if v != want[i] {
			t.Fatalf("record_bbs_in order: got %v, want %v", got, want)
		}
	}
	t.Logf("pre-fix lex=%v post-fix numeric=%v", buggy, got)
}

// sc_bbs_in_boundary verifies 50-item pagination across the 50/51 boundary with
// 60 posts: page 1 = 60..11, page 2 = 10..1, no gap/overlap.
func sc_bbs_in_boundary(db *sql.DB, t *testing.T) {
	for i := 1; i <= 60; i++ {
		record_test_seed_user_id(db, t, fmt.Sprintf("%d", i), "5", "admin_ip")
	}

	page1 := func() []string {
		res := route.Api_record_bbs_in(tool.Config{IP: "admin_ip"}, "admin_ip", "5", "1")
		got, _ := res["data"].([]string)
		return got
	}
	page2 := func() []string {
		res := route.Api_record_bbs_in(tool.Config{IP: "admin_ip"}, "admin_ip", "5", "2")
		got, _ := res["data"].([]string)
		return got
	}

	p1 := page1()
	p2 := page2()
	t.Logf("page1 (first/last): %q ... %q", p1[0], p1[len(p1)-1])
	for _, v := range p1 {
		t.Logf("page1: %q", v)
	}
	t.Logf("page2 (first/last): %q ... %q", p2[0], p2[len(p2)-1])
	for _, v := range p2 {
		t.Logf("page2: %q", v)
	}
	if len(p1) != 50 {
		t.Fatalf("page1 len = %d, want 50", len(p1))
	}
	if len(p2) != 10 {
		t.Fatalf("page2 len = %d, want 10", len(p2))
	}
	if p1[0] != "60" || p1[len(p1)-1] != "11" {
		t.Fatalf("page1 boundary wrong: first=%q last=%q", p1[0], p1[len(p1)-1])
	}
	if p2[0] != "10" || p2[len(p2)-1] != "1" {
		t.Fatalf("page2 boundary wrong: first=%q last=%q", p2[0], p2[len(p2)-1])
	}

	// No overlap and no gap between the two pages.
	merged := append(append([]string{}, p1...), p2...)
	seen := make(map[string]int)
	for _, v := range merged {
		seen[v]++
	}
	if len(seen) != 60 {
		t.Fatalf("union of pages = %d distinct, want 60 (gap/overlap)", len(seen))
	}
	for _, v := range merged {
		if seen[v] != 1 {
			t.Fatalf("duplicate or missing in pages: %q count=%d", v, seen[v])
		}
	}
}

// sc_bbs_numeric verifies the post board list is sorted by board number (numeric)
// descending.
func sc_bbs_numeric(db *sql.DB, t *testing.T) {
	for _, set_id := range []string{"2", "9", "10", "11"} {
		record_test_seed_user_id(db, t, "1", set_id, "admin_ip")
	}

	res := route.Api_record_bbs(tool.Config{IP: "admin_ip"}, "admin_ip", "1")
	if res["response"] != "ok" {
		t.Fatalf("response = %v, want ok", res["response"])
	}
	got, ok := res["data"].([][]string)
	if !ok {
		t.Fatalf("data is not [][]string: %v", res["data"])
	}
	if len(got) != 4 {
		t.Fatalf("len = %d, want 4", len(got))
	}
	want := []string{"11", "10", "9", "2"}
	for i, v := range got {
		if v[0] != want[i] {
			t.Fatalf("record_bbs order: got %v, want %v", got, want)
		}
	}
}

// sc_bbs_comment_numeric verifies the comment board list is sorted by board
// number (numeric) descending.
func sc_bbs_comment_numeric(db *sql.DB, t *testing.T) {
	for _, set_id := range []string{"2-1", "9-1", "10-1", "11-1"} {
		record_test_seed_comment(db, t, set_id, "1", "2026-01-01 00:00:00", "admin_ip")
	}

	res := route.Api_record_bbs_comment(tool.Config{IP: "admin_ip"}, "admin_ip", "1")
	if res["response"] != "ok" {
		t.Fatalf("response = %v, want ok", res["response"])
	}
	got, ok := res["data"].([][]string)
	if !ok {
		t.Fatalf("data is not [][]string: %v", res["data"])
	}
	if len(got) != 4 {
		t.Fatalf("len = %d, want 4", len(got))
	}
	want := []string{"11", "10", "9", "2"}
	for i, v := range got {
		if v[0] != want[i] {
			t.Fatalf("record_bbs_comment order: got %v, want %v", got, want)
		}
	}
}

// sc_bbs_comment_in_date verifies the per-board comment record is sorted by
// comment DATE (chronological), not by set_code. set_code is assigned
// independently per post, so across posts in one board set_code != time.
func sc_bbs_comment_in_date(db *sql.DB, t *testing.T) {
	// Board 5, two posts (5-3, 5-7) with interleaved dates vs per-post code.
	//  5-3/code2 -> old date
	//  5-7/code1 -> new date
	record_test_seed_comment(db, t, "5-3", "2", "2026-10-01 00:00:00", "admin_ip")
	record_test_seed_comment(db, t, "5-7", "1", "2026-10-05 00:00:00", "admin_ip")

	// Reproduce the pre-fix bug (set_code desc across posts gives wrong order).
	var buggy []string
	brows, berr := db.Query(`select set_id, set_code from bbs_data where set_name = 'comment_user_id' and set_id like '5-%' and set_data = 'admin_ip' order by set_code desc limit 0, 50`)
	defer brows.Close()
	if berr != nil {
		t.Fatalf("buggy probe: %v", berr)
	}
	for brows.Next() {
		var sid, code string
		if err := brows.Scan(&sid, &code); err != nil {
			t.Fatalf("buggy probe scan: %v", err)
		}
		buggy = append(buggy, sid+"["+code+"]")
	}
	if len(buggy) == 0 || buggy[0] != "5-3[2]" {
		t.Fatalf("pre-fix set_code order not reproduced; got %v", buggy)
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
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[0][2] != "5-7" || got[1][2] != "5-3" {
		t.Fatalf("record_bbs_comment_in order: got set_ids %v, want [5-7, 5-3]",
			[]string{got[0][2], got[1][2]})
	}
	if got[0][1] != "7" || got[1][1] != "3" {
		t.Fatalf("record_bbs_comment_in post codes: got %v, want [7, 3]",
			[]string{got[0][1], got[1][1]})
	}
	t.Logf("pre-fix set_code=%v post-fix date order=[%s %s]", buggy, got[0][2], got[1][2])
}

// sc_special_boards ensures special board identifiers ('0', '-1') are still
// returned and ordered after numeric boards without breaking.
func sc_special_boards(db *sql.DB, t *testing.T) {
	for _, set_id := range []string{"2", "10", "0", "-1"} {
		record_test_seed_user_id(db, t, "1", set_id, "admin_ip")
	}

	res := route.Api_record_bbs(tool.Config{IP: "admin_ip"}, "admin_ip", "1")
	if res["response"] != "ok" {
		t.Fatalf("response = %v, want ok", res["response"])
	}
	got, ok := res["data"].([][]string)
	if !ok {
		t.Fatalf("data is not [][]string: %v", res["data"])
	}
	if len(got) != 4 {
		t.Fatalf("len = %d, want 4", len(got))
	}
	want := []string{"10", "2", "0", "-1"}
	for i, v := range got {
		if v[0] != want[i] {
			t.Fatalf("special boards order: got %v, want %v", got, want)
		}
	}
}

// Test_record_bbs_sort runs every sort scenario against every available backend.
func Test_record_bbs_sort(t *testing.T) {
	backends := []string{"sqlite", "turso"}
	if os.Getenv("MYSQL_AUDIT_DSN") != "" {
		backends = append(backends, "mysql")
	}

	scenarios := []struct {
		name string
		fn   func(db *sql.DB, t *testing.T)
	}{
		{"post_in_numeric", sc_bbs_in_numeric},
		{"post_in_boundary_60", sc_bbs_in_boundary},
		{"board_numeric", sc_bbs_numeric},
		{"comment_board_numeric", sc_bbs_comment_numeric},
		{"comment_in_date", sc_bbs_comment_in_date},
		{"special_boards", sc_special_boards},
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
