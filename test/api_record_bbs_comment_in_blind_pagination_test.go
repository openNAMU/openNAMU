package test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

func record_blind_test_seed_post(db *sql.DB, t *testing.T, bbs_id, post_code, author, date, blind string) {
	exec := func(query string, args ...any) {
		if _, err := db.Exec(query, args...); err != nil {
			t.Fatalf("seed post %q/%q: %v", bbs_id, post_code, err)
		}
	}

	exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values ('user_id', ?, ?, ?)`,
		post_code, bbs_id, author)
	exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values ('title', ?, ?, ?)`,
		post_code, bbs_id, "Blind Test Post")
	exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values ('date', ?, ?, ?)`,
		post_code, bbs_id, date)
	if blind == "O" {
		exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values ('blind', ?, ?, ?)`,
			post_code, bbs_id, "O")
	}
}

func record_bbs_comment_in_data(t *testing.T, config tool.Config, user_name, bbs_id, page string) [][]string {
	res := route.Api_record_bbs_comment_in(config, user_name, bbs_id, page)
	if res["response"] != "ok" {
		t.Fatalf("response = %v, want ok", res)
	}
	data, ok := res["data"].([][]string)
	if !ok {
		t.Fatalf("data is not [][]string: %v", res["data"])
	}
	return data
}

func record_bbs_comment_in_old_visible_count(db *sql.DB, config tool.Config, user_name, bbs_id, page string) int {
	page_int := tool.Str_to_int(page)
	num := 0
	if page_int*50 > 0 {
		num = page_int*50 - 50
	}

	rows := tool.Query_DB(
		db,
		`select c.set_id, c.set_code, coalesce(d.set_data, '') as date
		 from bbs_data c
		 left join bbs_data d on d.set_name = "comment_date" and d.set_id = c.set_id and d.set_code = c.set_code
		 where c.set_name = "comment_user_id" and c.set_id like ? and c.set_data = ?
		 order by coalesce(d.set_data, '') desc, c.set_id desc limit ?, 50`,
		bbs_id+"-%",
		user_name,
		num,
	)
	defer rows.Close()

	var auth_info map[string]bool
	count := 0
	for rows.Next() {
		var set_id string
		var set_code string
		var date string

		if err := rows.Scan(&set_id, &set_code, &date); err != nil {
			panic(err)
		}

		post_set_id, post_set_code, _, exists := route.Bbs_comment_storage_location(set_id, set_code)
		if exists {
			if auth_info == nil {
				auth_info = tool.Get_auth_info(db, config.IP)
			}
			if !route.Bbs_post_blind_allowed(db, post_set_id, post_set_code, config.IP, auth_info) {
				continue
			}
		}
		count++
	}

	return count
}

func record_blind_test_scenario(db *sql.DB, t *testing.T) {
	record_blind_test_seed_post(db, t, "5", "3", "author", "2026-10-15 00:00:00", "O")
	for i := 1; i <= 50; i++ {
		record_test_seed_comment(db, t, "5-3", fmt.Sprintf("%d", i), "2026-10-15 12:00:00", "alice")
	}
	record_blind_test_seed_post(db, t, "5", "4", "author", "2026-10-15 00:00:00", "")
	record_test_seed_comment(db, t, "5-4", "1", "2026-10-15 11:00:00", "alice")
}

func sc_pre_fix_visible_missing(db *sql.DB, t *testing.T) {
	record_blind_test_scenario(db, t)

	regular_old := record_bbs_comment_in_old_visible_count(db, 	tool.Config{IP: "hidel_user"}, "alice", "5", "1")
	if regular_old != 0 {
		t.Fatalf("pre-fix regular page 1 visible count = %d, want 0 (old code must have hidden the visible comment)", regular_old)
	}

	admin_old := record_bbs_comment_in_old_visible_count(db, tool.Config{IP: "admin_ip"}, "alice", "5", "1")
	if admin_old != 50 {
		t.Fatalf("pre-fix admin page 1 visible count = %d, want 50 (visible comment missing from admin page 1)", admin_old)
	}
}

func sc_regular_visible_after_blind(db *sql.DB, t *testing.T) {
	record_blind_test_scenario(db, t)

	regular := record_bbs_comment_in_data(t, 	tool.Config{IP: "hidel_user"}, "alice", "5", "1")
	if len(regular) != 1 {
		t.Fatalf("regular page 1 len = %d, want 1 visible comment", len(regular))
	}
	if regular[0][1] != "4" {
		t.Fatalf("regular page 1 post code = %v, want 4 (non-blind post)", regular[0][1])
	}

	admin_page1 := record_bbs_comment_in_data(t, tool.Config{IP: "admin_ip"}, "alice", "5", "1")
	if len(admin_page1) != 50 {
		t.Fatalf("admin page 1 len = %d, want 50 blind comments", len(admin_page1))
	}
	for _, row := range admin_page1 {
		if row[1] == "4" {
			t.Fatalf("admin page 1 incorrectly contains visible comment on post 4")
		}
	}

	admin_page2 := record_bbs_comment_in_data(t, tool.Config{IP: "admin_ip"}, "alice", "5", "2")
	if len(admin_page2) != 1 {
		t.Fatalf("admin page 2 len = %d, want 1 visible comment", len(admin_page2))
	}
	if admin_page2[0][1] != "4" {
		t.Fatalf("admin page 2 post code = %v, want 4", admin_page2[0][1])
	}
}

func sc_visible_pagination_boundary_after_blind(db *sql.DB, t *testing.T) {
	record_blind_test_seed_post(db, t, "5", "3", "author", "2026-10-15 00:00:00", "O")
	for i := 1; i <= 50; i++ {
		record_test_seed_comment(db, t, "5-3", fmt.Sprintf("%d", i), "2026-10-15 12:00:00", "alice")
	}
	record_blind_test_seed_post(db, t, "5", "4", "author", "2026-10-15 00:00:00", "")
	for i := 1; i <= 100; i++ {
		record_test_seed_comment(db, t, "5-4", fmt.Sprintf("%d", i), "2026-10-10 00:00:00", "alice")
	}

	page1 := record_bbs_comment_in_data(t, 	tool.Config{IP: "hidel_user"}, "alice", "5", "1")
	page2 := record_bbs_comment_in_data(t, 	tool.Config{IP: "hidel_user"}, "alice", "5", "2")
	page3 := record_bbs_comment_in_data(t, 	tool.Config{IP: "hidel_user"}, "alice", "5", "3")

	if len(page1) != 50 {
		t.Fatalf("regular page 1 len = %d, want 50", len(page1))
	}
	if len(page2) != 50 {
		t.Fatalf("regular page 2 len = %d, want 50", len(page2))
	}
	if len(page3) != 0 {
		t.Fatalf("regular page 3 len = %d, want 0", len(page3))
	}

	check_no_blind := func(rows [][]string, page int) {
		for _, row := range rows {
			if row[1] == "3" {
				t.Fatalf("regular page %d contains blind post 3: %v", page, row)
			}
		}
	}
	check_no_blind(page1, 1)
	check_no_blind(page2, 2)

	keys := func(rows [][]string) map[string]bool {
		out := map[string]bool{}
		for _, row := range rows {
			out[row[2]+"-"+row[3]] = true
		}
		return out
	}
	page1_keys := keys(page1)
	page2_keys := keys(page2)
	if len(page1_keys) != 50 {
		t.Fatalf("regular page 1 distinct comments = %d, want 50", len(page1_keys))
	}
	if len(page2_keys) != 50 {
		t.Fatalf("regular page 2 distinct comments = %d, want 50", len(page2_keys))
	}
	for key := range page1_keys {
		if _, ok := page2_keys[key]; ok {
			t.Fatalf("regular pages 1 and 2 overlap on %q", key)
		}
	}
}

func sc_nested_blind_hidden_for_regular(db *sql.DB, t *testing.T) {
	record_blind_test_seed_post(db, t, "5", "7", "author", "2026-10-15 00:00:00", "O")
	record_blind_test_seed_post(db, t, "5", "8", "author", "2026-10-15 00:00:00", "")
	record_test_seed_comment(db, t, "5-7", "1", "2026-10-15 12:00:00", "alice")
	record_test_seed_comment(db, t, "5-7-1", "2", "2026-10-15 11:00:00", "alice")
	record_test_seed_comment(db, t, "5-8", "1", "2026-10-15 10:00:00", "alice")

	regular := record_bbs_comment_in_data(t, 	tool.Config{IP: "hidel_user"}, "alice", "5", "1")
	if len(regular) != 1 {
		t.Fatalf("regular nested page 1 len = %d, want 1 visible comment", len(regular))
	}
	if regular[0][1] != "8" {
		t.Fatalf("regular nested page 1 post code = %v, want 8", regular[0][1])
	}

	admin := record_bbs_comment_in_data(t, tool.Config{IP: "admin_ip"}, "alice", "5", "1")
	if len(admin) != 3 {
		t.Fatalf("admin nested page 1 len = %d, want 3 (2 blind + 1 visible)", len(admin))
	}
}

func sc_minus1_blind_pagination(db *sql.DB, t *testing.T) {
	record_blind_test_seed_post(db, t, "-1", "2025", "author", "2026-10-15 00:00:00", "O")
	for i := 1; i <= 50; i++ {
		record_test_seed_comment(db, t, "-1-2025", fmt.Sprintf("%d", i), "2026-10-15 12:00:00", "alice")
	}
	record_blind_test_seed_post(db, t, "-1", "2026", "author", "2026-10-15 00:00:00", "")
	record_test_seed_comment(db, t, "-1-2026", "1", "2026-10-15 11:00:00", "alice")

	regular := record_bbs_comment_in_data(t, 	tool.Config{IP: "hidel_user"}, "alice", "-1", "1")
	if len(regular) != 1 {
		t.Fatalf("minus1 regular page 1 len = %d, want 1 visible comment", len(regular))
	}
	if regular[0][1] != "2026" {
		t.Fatalf("minus1 regular page 1 post code = %v, want 2026", regular[0][1])
	}

	admin_page1 := record_bbs_comment_in_data(t, tool.Config{IP: "admin_ip"}, "alice", "-1", "1")
	if len(admin_page1) != 50 {
		t.Fatalf("minus1 admin page 1 len = %d, want 50 blind comments", len(admin_page1))
	}
	for _, row := range admin_page1 {
		if row[1] == "2026" {
			t.Fatalf("minus1 admin page 1 incorrectly contains visible comment on post 2026")
		}
	}

	admin_page2 := record_bbs_comment_in_data(t, tool.Config{IP: "admin_ip"}, "alice", "-1", "2")
	if len(admin_page2) != 1 {
		t.Fatalf("minus1 admin page 2 len = %d, want 1 visible comment", len(admin_page2))
	}
	if admin_page2[0][1] != "2026" {
		t.Fatalf("minus1 admin page 2 post code = %v, want 2026", admin_page2[0][1])
	}
}

func Test_record_bbs_comment_in_blind_pagination(t *testing.T) {
	backends := []string{"sqlite", "turso"}
	if os.Getenv("MYSQL_AUDIT_DSN") != "" {
		backends = append(backends, "mysql")
	}

	scenarios := []struct {
		name string
		fn   func(db *sql.DB, t *testing.T)
	}{
		{"pre_fix_visible_missing", sc_pre_fix_visible_missing},
		{"regular_visible_after_blind", sc_regular_visible_after_blind},
		{"visible_pagination_boundary_after_blind", sc_visible_pagination_boundary_after_blind},
		{"nested_blind_hidden_for_regular", sc_nested_blind_hidden_for_regular},
		{"minus1_blind_pagination", sc_minus1_blind_pagination},
	}

	for _, b := range backends {
		for _, sc := range scenarios {
			t.Run(sc.name+"_"+b, func(t *testing.T) {
				db := new_record_db(t, b)
				seed_record_access_auth(db, t)
				sc.fn(db, t)
			})
		}
	}
}
