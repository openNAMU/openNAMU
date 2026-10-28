package test

import (
	"database/sql"
	"os"
	"strings"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

func Test_thread_watchlist_legacy_visibility(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso", "mysql"} {
		if backend == "mysql" && os.Getenv("MYSQL_AUDIT_DSN") == "" {
			continue
		}
		t.Run(backend, func(t *testing.T) {
			db := new_record_db(t, backend)
			defer tool.DB_close(db)

			thread_seed_bbs(db, t, "1", "thread 1")
			thread_seed_bbs(db, t, "2", "thread 2")
			thread_seed_bbs(db, t, "3", "thread 3")
			seed_user_auth(db, t, "legacy_user")
			seed_user_auth(db, t, "new_user")
			seed_user_auth(db, t, "dup_user")

			t.Run("legacy_thread_only_is_visible_in_bbs_watch_list", func(t *testing.T) {
				seed_user_watch(db, t, "thread_watchlist", "legacy_user", "1")

				html := route.View_user_watch_list(tool.Config{IP: "legacy_user", Session: &emptySession{}}, "1", "bbs_watchlist")
				if !strings.Contains(html, `href="/bbs/w/-1/1">`) {
					t.Fatalf("legacy thread 1 is missing from /bbs_watch_list")
				}
				if !strings.Contains(html, `/bbs_watch/-1/1`) {
					t.Fatalf("unwatch link for legacy thread 1 is missing")
				}

				res := route.Api_user_watch_list(tool.Config{IP: "legacy_user", Session: &emptySession{}}, "legacy_user", "1", "bbs_watchlist")
				items, ok := res["data"].([]route.User_watch_item)
				if !ok {
					t.Fatalf("api data is not []route.User_watch_item")
				}
				if len(items) != 1 {
					t.Fatalf("legacy only items = %d, want 1", len(items))
				}
				if items[0].Set_id != "-1" || items[0].Set_code != "1" {
					t.Fatalf("legacy item set = %v/%v, want -1/1", items[0].Set_id, items[0].Set_code)
				}
				if items[0].Display != "thread 1" {
					t.Fatalf("legacy display = %q, want thread 1", items[0].Display)
				}

				legacy_post := route.Api_bbs_watch_post(tool.Config{IP: "legacy_user", Session: &emptySession{}}, "-1", "1")
				if legacy_post["response"] != "ok" {
					t.Fatalf("unwatch legacy: %v", legacy_post)
				}
				if count_user_watch_rows(db, "thread_watchlist", "legacy_user", "1") != 0 {
					t.Fatalf("legacy thread row remains after unwatch")
				}
				if count_user_watch_rows(db, "bbs_watchlist", "legacy_user", "-1-1") != 0 {
					t.Fatalf("legacy unwatch created bbs_watchlist row")
				}
			})

			t.Run("new_bbs_only_still_visible", func(t *testing.T) {
				seed_user_watch(db, t, "bbs_watchlist", "new_user", "-1-2")

				html := route.View_user_watch_list(tool.Config{IP: "new_user", Session: &emptySession{}}, "1", "bbs_watchlist")
				if !strings.Contains(html, `href="/bbs/w/-1/2">`) {
					t.Fatalf("new bbs thread 2 is missing from /bbs_watch_list")
				}

				res := route.Api_user_watch_list(tool.Config{IP: "new_user", Session: &emptySession{}}, "new_user", "1", "bbs_watchlist")
				items, ok := res["data"].([]route.User_watch_item)
				if !ok {
					t.Fatalf("api data is not []route.User_watch_item")
				}
				if len(items) != 1 {
					t.Fatalf("new bbs items = %d, want 1", len(items))
				}
				if items[0].Set_id != "-1" || items[0].Set_code != "2" {
					t.Fatalf("new bbs item set = %v/%v, want -1/2", items[0].Set_id, items[0].Set_code)
				}

				new_bbs_post := route.Api_bbs_watch_post(tool.Config{IP: "new_user", Session: &emptySession{}}, "-1", "2")
				if new_bbs_post["response"] != "ok" {
					t.Fatalf("unwatch new bbs: %v", new_bbs_post)
				}
				if count_user_watch_rows(db, "bbs_watchlist", "new_user", "-1-2") != 0 {
					t.Fatalf("new bbs row remains after unwatch")
				}
			})

			t.Run("duplicate_legacy_and_bbs_shown_once", func(t *testing.T) {
				seed_user_watch(db, t, "bbs_watchlist", "dup_user", "-1-3")
				seed_user_watch(db, t, "thread_watchlist", "dup_user", "3")

				html := route.View_user_watch_list(tool.Config{IP: "dup_user", Session: &emptySession{}}, "1", "bbs_watchlist")
				if strings.Count(html, `href="/bbs/w/-1/3">`) != 1 {
					t.Fatalf("duplicate thread 3 count = %d, want 1", strings.Count(html, `href="/bbs/w/-1/3">`))
				}

				res := route.Api_user_watch_list(tool.Config{IP: "dup_user", Session: &emptySession{}}, "dup_user", "1", "bbs_watchlist")
				items, ok := res["data"].([]route.User_watch_item)
				if !ok {
					t.Fatalf("api data is not []route.User_watch_item")
				}
				if len(items) != 1 {
					t.Fatalf("duplicate items = %d, want 1", len(items))
				}
				if items[0].Set_id != "-1" || items[0].Set_code != "3" {
					t.Fatalf("duplicate item set = %v/%v, want -1/3", items[0].Set_id, items[0].Set_code)
				}

				dup_post := route.Api_bbs_watch_post(tool.Config{IP: "dup_user", Session: &emptySession{}}, "-1", "3")
				if dup_post["response"] != "ok" {
					t.Fatalf("unwatch duplicate: %v", dup_post)
				}
				if count_user_watch_rows(db, "thread_watchlist", "dup_user", "3") != 0 {
					t.Fatalf("duplicate thread row remains after unwatch")
				}
				if count_user_watch_rows(db, "bbs_watchlist", "dup_user", "-1-3") != 0 {
					t.Fatalf("duplicate bbs row remains after unwatch")
				}
			})

			t.Run("adding_missing_thread_inserts_bbs_watchlist", func(t *testing.T) {
				html := route.View_user_watch_list(tool.Config{IP: "new_user", Session: &emptySession{}}, "1", "bbs_watchlist")
				if !strings.Contains(html, `href="/bbs/w/-1/3">`) {
					// thread 3 is not watched by new_user yet; do not assert in html.
				}

				add_thread_post := route.Api_bbs_watch_post(tool.Config{IP: "new_user", Session: &emptySession{}}, "-1", "3")
				if add_thread_post["response"] != "ok" {
					t.Fatalf("add thread 3: %v", add_thread_post)
				}
				if count_user_watch_rows(db, "bbs_watchlist", "new_user", "-1-3") != 1 {
					t.Fatalf("added bbs_watchlist row count = %d, want 1", count_user_watch_rows(db, "bbs_watchlist", "new_user", "-1-3"))
				}
				if count_user_watch_rows(db, "thread_watchlist", "new_user", "3") != 0 {
					t.Fatalf("add created legacy thread_watchlist row")
				}

				html = route.View_user_watch_list(tool.Config{IP: "new_user", Session: &emptySession{}}, "1", "bbs_watchlist")
				if !strings.Contains(html, `href="/bbs/w/-1/3">`) {
					t.Fatalf("added thread 3 is missing from /bbs_watch_list")
				}
			})
		})
	}
}

func thread_seed_bbs(db *sql.DB, t *testing.T, set_code, title string) {
	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed bbs %q: %v", q, err)
		}
	}
	exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`, "title", set_code, "-1", title)
	exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`, "date", set_code, "-1", "2026-01-01 00:00:00")
	exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`, "user_id", set_code, "-1", "bbs_author")
}

func seed_user_auth(db *sql.DB, t *testing.T, user string) {
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl', ?, 'admin')`, user); err != nil {
		t.Fatalf("seed user auth %s: %v", user, err)
	}
}

func seed_user_watch(db *sql.DB, t *testing.T, watch_type, user, data string) {
	if _, err := db.Exec(`insert into user_set (name, id, data) values (?, ?, ?)`, watch_type, user, data); err != nil {
		t.Fatalf("seed user watch %s/%s/%s: %v", watch_type, user, data, err)
	}
}

func count_user_watch_rows(db *sql.DB, watch_type, user, data string) int {
	var count int
	err := db.QueryRow(`select count(*) from user_set where name = ? and id = ? and data = ?`, watch_type, user, data).Scan(&count)
	if err != nil {
		return 0
	}
	return count
}
