package test

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

func seed_record_access_auth(db *sql.DB, t *testing.T) {
	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed %q: %v", q, err)
		}
	}
	exec(`insert into alist (name, acl) values ('hidel_group', 'hidel')`)
	exec(`insert into alist (name, acl) values ('view_group', 'view')`)
	exec(`insert into user_set (name, id, data) values ('acl', 'hidel_user', 'hidel_group')`)
	exec(`insert into user_set (name, id, data) values ('acl', 'hidel_user', 'view_group')`)
	exec(`insert into user_set (name, id, data) values ('acl', 'regular', 'view_group')`)
	exec(`insert into user_set (name, id, data) values ('acl', 'alice', 'view_group')`)
}

func record_test_seed_post(db *sql.DB, t *testing.T, set_id, set_code, user, title, date string) {
	record_test_seed_user_id(db, t, set_code, set_id, user)
	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed %q: %v", q, err)
		}
	}
	exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values ('title', ?, ?, ?)`, set_code, set_id, title)
	exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values ('date', ?, ?, ?)`, set_code, set_id, date)
}

func seed_record_access_posts(db *sql.DB, t *testing.T) {
	// 60 normal boards for Api_record_bbs page 2.
	for i := 1; i <= 60; i++ {
		set_id := fmt.Sprintf("%d", i)
		record_test_seed_post(db, t, set_id, "3", "alice", fmt.Sprintf("Post %d", i), "2026-01-01 00:00:00")
	}

	// Board 200 with 60 posts for Api_record_bbs_in page 2.
	for i := 1; i <= 60; i++ {
		set_code := fmt.Sprintf("%d", i)
		record_test_seed_post(db, t, "200", set_code, "alice", fmt.Sprintf("Board 200 Post %d", i), "2026-01-02 00:00:00")
	}

	// 60 normal boards with comments for Api_record_bbs_comment page 2.
	for i := 1; i <= 60; i++ {
		set_id := fmt.Sprintf("%d-3", i)
		record_test_seed_comment(db, t, set_id, "1", "2026-02-01 00:00:00", "alice")
	}

	// Board 200 with 60 comments for Api_record_bbs_comment_in page 2.
	for i := 1; i <= 60; i++ {
		set_code := fmt.Sprintf("%d", i)
		record_test_seed_comment(db, t, "200-3", set_code, "2026-02-02 00:00:00", "alice")
	}

	// Unique post for View_record_bbs / View_record_bbs_in assertions.
	record_test_seed_post(db, t, "999", "3", "alice", "Access Post", "2026-03-01 00:00:00")
}

func assert_api_denied(t *testing.T, name string, res map[string]any) {
	if response := res["response"]; response != "require auth" {
		t.Fatalf("%s response = %v, want require auth", name, res)
	}
}

func assert_api_ok(t *testing.T, name string, res map[string]any, data_len int) {
	if response := res["response"]; response != "ok" {
		t.Fatalf("%s response = %v, want ok", name, res)
	}
	switch v := res["data"].(type) {
	case [][]string:
		if len(v) < data_len {
			t.Fatalf("%s data len = %d, want >= %d", name, len(v), data_len)
		}
	case []string:
		if len(v) < data_len {
			t.Fatalf("%s data len = %d, want >= %d", name, len(v), data_len)
		}
	default:
		t.Fatalf("%s data type = %T, want [][]string or []string", name, res["data"])
	}
}

func Test_record_bbs_access(t *testing.T) {
	backends := []string{"sqlite", "turso"}
	if os.Getenv("MYSQL_AUDIT_DSN") != "" {
		backends = append(backends, "mysql")
	}

	for _, backend := range backends {
		t.Run(backend, func(t *testing.T) {
			db := new_record_db(t, backend)
			defer tool.DB_close(db)

			seed_record_access_auth(db, t)
			seed_record_access_posts(db, t)

			t.Run("Api_record_bbs_denied_regular", func(t *testing.T) {
				for _, page := range []string{"1", "2"} {
					res := route.Api_record_bbs(tool.Config{IP: "regular"}, "alice", page)
					assert_api_denied(t, "Api_record_bbs regular "+page, res)
				}
			})
			t.Run("Api_record_bbs_allowed_self", func(t *testing.T) {
				for page, min_len := range map[string]int{"1": 1, "2": 1} {
					res := route.Api_record_bbs(tool.Config{IP: "alice"}, "alice", page)
					assert_api_ok(t, "Api_record_bbs self "+page, res, min_len)
				}
			})
			t.Run("Api_record_bbs_allowed_hidel", func(t *testing.T) {
				for page, min_len := range map[string]int{"1": 1, "2": 1} {
					res := route.Api_record_bbs(tool.Config{IP: "hidel_user"}, "alice", page)
					assert_api_ok(t, "Api_record_bbs hidel "+page, res, min_len)
				}
			})

			t.Run("Api_record_bbs_in_denied_regular", func(t *testing.T) {
				for _, page := range []string{"1", "2"} {
					res := route.Api_record_bbs_in(tool.Config{IP: "regular"}, "alice", "200", page)
					assert_api_denied(t, "Api_record_bbs_in regular "+page, res)
				}
			})
			t.Run("Api_record_bbs_in_allowed_self", func(t *testing.T) {
				for page, min_len := range map[string]int{"1": 1, "2": 1} {
					res := route.Api_record_bbs_in(tool.Config{IP: "alice"}, "alice", "200", page)
					assert_api_ok(t, "Api_record_bbs_in self "+page, res, min_len)
				}
			})
			t.Run("Api_record_bbs_in_allowed_hidel", func(t *testing.T) {
				for page, min_len := range map[string]int{"1": 1, "2": 1} {
					res := route.Api_record_bbs_in(tool.Config{IP: "hidel_user"}, "alice", "200", page)
					assert_api_ok(t, "Api_record_bbs_in hidel "+page, res, min_len)
				}
			})
			t.Run("Api_record_bbs_in_special_report_denied_even_hidel", func(t *testing.T) {
				res := route.Api_record_bbs_in(tool.Config{IP: "hidel_user"}, "alice", "-2", "1")
				assert_api_denied(t, "Api_record_bbs_in -2 hidel", res)
				res = route.Api_record_bbs_in(tool.Config{IP: "alice"}, "alice", "-2", "1")
				assert_api_denied(t, "Api_record_bbs_in -2 self", res)
				res = route.Api_record_bbs_in(tool.Config{IP: "regular"}, "alice", "-2", "1")
				assert_api_denied(t, "Api_record_bbs_in -2 regular", res)
			})
			t.Run("Api_record_bbs_in_special_thread_allowed_ok_empty", func(t *testing.T) {
				res := route.Api_record_bbs_in(tool.Config{IP: "alice"}, "alice", "-1", "1")
				if res["response"] != "ok" {
					t.Fatalf("Api_record_bbs_in -1 self response = %v, want ok", res)
				}
				res = route.Api_record_bbs_in(tool.Config{IP: "hidel_user"}, "alice", "-1", "1")
				if res["response"] != "ok" {
					t.Fatalf("Api_record_bbs_in -1 hidel response = %v, want ok", res)
				}
			})
			t.Run("Api_record_bbs_in_special_wiki_allowed_ok_empty", func(t *testing.T) {
				res := route.Api_record_bbs_in(tool.Config{IP: "alice"}, "alice", "0", "1")
				if res["response"] != "ok" {
					t.Fatalf("Api_record_bbs_in 0 self response = %v, want ok", res)
				}
				res = route.Api_record_bbs_in(tool.Config{IP: "hidel_user"}, "alice", "0", "1")
				if res["response"] != "ok" {
					t.Fatalf("Api_record_bbs_in 0 hidel response = %v, want ok", res)
				}
			})

			t.Run("Api_record_bbs_comment_denied_regular", func(t *testing.T) {
				for _, page := range []string{"1", "2"} {
					res := route.Api_record_bbs_comment(tool.Config{IP: "regular"}, "alice", page)
					assert_api_denied(t, "Api_record_bbs_comment regular "+page, res)
				}
			})
			t.Run("Api_record_bbs_comment_allowed_self", func(t *testing.T) {
				for page, min_len := range map[string]int{"1": 1, "2": 1} {
					res := route.Api_record_bbs_comment(tool.Config{IP: "alice"}, "alice", page)
					assert_api_ok(t, "Api_record_bbs_comment self "+page, res, min_len)
				}
			})
			t.Run("Api_record_bbs_comment_allowed_hidel", func(t *testing.T) {
				for page, min_len := range map[string]int{"1": 1, "2": 1} {
					res := route.Api_record_bbs_comment(tool.Config{IP: "hidel_user"}, "alice", page)
					assert_api_ok(t, "Api_record_bbs_comment hidel "+page, res, min_len)
				}
			})

			t.Run("Api_record_bbs_comment_in_denied_regular", func(t *testing.T) {
				for _, page := range []string{"1", "2"} {
					res := route.Api_record_bbs_comment_in(tool.Config{IP: "regular"}, "alice", "200", page)
					assert_api_denied(t, "Api_record_bbs_comment_in regular "+page, res)
				}
			})
			t.Run("Api_record_bbs_comment_in_allowed_self", func(t *testing.T) {
				for page, min_len := range map[string]int{"1": 1, "2": 1} {
					res := route.Api_record_bbs_comment_in(tool.Config{IP: "alice"}, "alice", "200", page)
					assert_api_ok(t, "Api_record_bbs_comment_in self "+page, res, min_len)
				}
			})
			t.Run("Api_record_bbs_comment_in_allowed_hidel", func(t *testing.T) {
				for page, min_len := range map[string]int{"1": 1, "2": 1} {
					res := route.Api_record_bbs_comment_in(tool.Config{IP: "hidel_user"}, "alice", "200", page)
					assert_api_ok(t, "Api_record_bbs_comment_in hidel "+page, res, min_len)
				}
			})

			t.Run("View_record_bbs_denied_regular", func(t *testing.T) {
				out := route.View_record_bbs(tool.Config{IP: "regular", Session: &emptySession{}}, "alice", "1")
				if strings.Contains(out, "Access Post") || strings.Contains(out, "2026-03-01") {
					t.Fatalf("View_record_bbs regular exposed Access Post")
				}
			})
			t.Run("View_record_bbs_allowed_self", func(t *testing.T) {
				out := route.View_record_bbs(tool.Config{IP: "alice", Session: &emptySession{}}, "alice", "1")
				if !strings.Contains(out, "Access Post") && !strings.Contains(out, "2026-03-01") {
					t.Fatalf("View_record_bbs self missing Access Post")
				}
			})
			t.Run("View_record_bbs_allowed_hidel", func(t *testing.T) {
				out := route.View_record_bbs(tool.Config{IP: "hidel_user", Session: &emptySession{}}, "alice", "1")
				if !strings.Contains(out, "Access Post") && !strings.Contains(out, "2026-03-01") {
					t.Fatalf("View_record_bbs hidel missing Access Post")
				}
			})
			t.Run("View_record_bbs_in_denied_regular", func(t *testing.T) {
				out := route.View_record_bbs_in(tool.Config{IP: "regular", Session: &emptySession{}}, "alice", "999", "1")
				if strings.Contains(out, "Access Post") || strings.Contains(out, "2026-03-01") || strings.Contains(out, "/bbs/w/999/3") {
					t.Fatalf("View_record_bbs_in regular exposed Access Post")
				}
			})
			t.Run("View_record_bbs_in_allowed_self", func(t *testing.T) {
				out := route.View_record_bbs_in(tool.Config{IP: "alice", Session: &emptySession{}}, "alice", "999", "1")
				if !strings.Contains(out, "Access Post") && !strings.Contains(out, "2026-03-01") && !strings.Contains(out, "/bbs/w/999/3") {
					t.Fatalf("View_record_bbs_in self missing Access Post")
				}
			})
			t.Run("View_record_bbs_in_allowed_hidel", func(t *testing.T) {
				out := route.View_record_bbs_in(tool.Config{IP: "hidel_user", Session: &emptySession{}}, "alice", "999", "1")
				if !strings.Contains(out, "Access Post") && !strings.Contains(out, "2026-03-01") && !strings.Contains(out, "/bbs/w/999/3") {
					t.Fatalf("View_record_bbs_in hidel missing Access Post")
				}
			})
		})
	}
}


