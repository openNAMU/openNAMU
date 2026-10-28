package test

import (
	"database/sql"
	"os"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

func Test_bbs_comment_duplicate_notification(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso", "mysql"} {
		if backend == "mysql" && os.Getenv("MYSQL_AUDIT_DSN") == "" {
			continue
		}
		t.Run(backend, func(t *testing.T) {
			db := new_record_db(t, backend)
			defer tool.DB_close(db)

			t.Run("same_author_and_parent_comment_once", func(t *testing.T) {
				bbs_comment_seed_board(db, t, "7", "1", "post 1", "owner_parent")
				bbs_comment_seed_parent(db, t, "7", "1", "2", "owner_parent")
				seed_user_auth(db, t, "owner_parent")
				seed_user_auth(db, t, "reply_author")
				seed_user_watch(db, t, "bbs_watchlist", "owner_parent", "7-1")
				seed_user_watch(db, t, "bbs_watchlist", "watcher1", "7-1")

				res := route.Api_bbs_w_comment_post(
					tool.Config{IP: "reply_author", Session: &emptySession{}},
					"7",
					"1",
					"2",
					"reply data",
				)
				if res["response"] != "ok" {
					t.Fatalf("response = %v, want ok", res["response"])
				}

				if count := notice_count(db, "owner_parent"); count != 1 {
					t.Fatalf("owner_parent notice count = %d, want 1 (post_user and parent_user are the same)", count)
				}
				if count := notice_count(db, "watcher1"); count != 1 {
					t.Fatalf("watcher1 notice count = %d, want 1", count)
				}
				if count := notice_count(db, "reply_author"); count != 0 {
					t.Fatalf("reply_author notice count = %d, want 0", count)
				}
			})

			t.Run("different_author_and_parent_no_duplicate", func(t *testing.T) {
				bbs_comment_seed_board(db, t, "8", "1", "post 2", "author_one")
				bbs_comment_seed_parent(db, t, "8", "1", "2", "parent_two")
				seed_user_auth(db, t, "author_one")
				seed_user_auth(db, t, "parent_two")
				seed_user_auth(db, t, "commenter")

				res := route.Api_bbs_w_comment_post(
					tool.Config{IP: "commenter", Session: &emptySession{}},
					"8",
					"1",
					"2",
					"reply data",
				)
				if res["response"] != "ok" {
					t.Fatalf("response = %v, want ok", res["response"])
				}

				if count := notice_count(db, "author_one"); count != 1 {
					t.Fatalf("author_one notice count = %d, want 1", count)
				}
				if count := notice_count(db, "parent_two"); count != 1 {
					t.Fatalf("parent_two notice count = %d, want 1", count)
				}
				if count := notice_count(db, "commenter"); count != 0 {
					t.Fatalf("commenter notice count = %d, want 0", count)
				}
			})

			t.Run("reply_author_is_post_user", func(t *testing.T) {
				bbs_comment_seed_board(db, t, "9", "1", "post 3", "reply_owner")
				seed_user_auth(db, t, "reply_owner")
				seed_user_watch(db, t, "bbs_watchlist", "reply_owner", "9-1")

				res := route.Api_bbs_w_comment_post(
					tool.Config{IP: "reply_owner", Session: &emptySession{}},
					"9",
					"1",
					"0",
					"top comment by owner",
				)
				if res["response"] != "ok" {
					t.Fatalf("response = %v, want ok", res["response"])
				}

				if count := notice_count(db, "reply_owner"); count != 0 {
					t.Fatalf("reply_owner notice count = %d, want 0 (self notification excluded)", count)
				}
			})

			t.Run("reply_author_is_parent_user", func(t *testing.T) {
				bbs_comment_seed_board(db, t, "10", "1", "post 4", "author")
				bbs_comment_seed_parent(db, t, "10", "1", "2", "parent")
				seed_user_auth(db, t, "author")
				seed_user_auth(db, t, "parent")
				seed_user_watch(db, t, "bbs_watchlist", "author", "10-1")
				seed_user_watch(db, t, "bbs_watchlist", "parent", "10-1")

				res := route.Api_bbs_w_comment_post(
					tool.Config{IP: "parent", Session: &emptySession{}},
					"10",
					"1",
					"2",
					"reply by parent",
				)
				if res["response"] != "ok" {
					t.Fatalf("response = %v, want ok", res["response"])
				}

				if count := notice_count(db, "author"); count != 1 {
					t.Fatalf("author notice count = %d, want 1", count)
				}
				if count := notice_count(db, "parent"); count != 0 {
					t.Fatalf("parent notice count = %d, want 0 (self notification excluded)", count)
				}
			})

			t.Run("top_level_comment", func(t *testing.T) {
				bbs_comment_seed_board(db, t, "11", "1", "post 5", "top_owner")
				seed_user_auth(db, t, "top_owner")
				seed_user_auth(db, t, "top_commenter")
				seed_user_watch(db, t, "bbs_watchlist", "top_owner", "11-1")
				seed_user_watch(db, t, "bbs_watchlist", "top_commenter", "11-1")

				res := route.Api_bbs_w_comment_post(
					tool.Config{IP: "top_commenter", Session: &emptySession{}},
					"11",
					"1",
					"0",
					"top-level comment",
				)
				if res["response"] != "ok" {
					t.Fatalf("response = %v, want ok", res["response"])
				}

				if count := notice_count(db, "top_owner"); count != 1 {
					t.Fatalf("top_owner notice count = %d, want 1", count)
				}
				if count := notice_count(db, "top_commenter"); count != 0 {
					t.Fatalf("top_commenter notice count = %d, want 0 (self notification excluded)", count)
				}
			})

			t.Run("watcher_overlaps_direct_recipients", func(t *testing.T) {
				bbs_comment_seed_board(db, t, "12", "1", "post 6", "direct_owner")
				bbs_comment_seed_parent(db, t, "12", "1", "2", "direct_parent")
				seed_user_auth(db, t, "direct_owner")
				seed_user_auth(db, t, "direct_parent")
				seed_user_auth(db, t, "commenter")
				seed_user_watch(db, t, "bbs_watchlist", "direct_owner", "12-1")
				seed_user_watch(db, t, "bbs_watchlist", "direct_parent", "12-1")
				seed_user_watch(db, t, "bbs_watchlist", "watcher_overlap", "12-1")

				res := route.Api_bbs_w_comment_post(
					tool.Config{IP: "commenter", Session: &emptySession{}},
					"12",
					"1",
					"2",
					"reply data",
				)
				if res["response"] != "ok" {
					t.Fatalf("response = %v, want ok", res["response"])
				}

				if count := notice_count(db, "direct_owner"); count != 1 {
					t.Fatalf("direct_owner notice count = %d, want 1 (direct recipient watcher should not get a second watch notification)", count)
				}
				if count := notice_count(db, "direct_parent"); count != 1 {
					t.Fatalf("direct_parent notice count = %d, want 1 (direct recipient watcher should not get a second watch notification)", count)
				}
				if count := notice_count(db, "watcher_overlap"); count != 1 {
					t.Fatalf("watcher_overlap notice count = %d, want 1", count)
				}
				if count := notice_count(db, "commenter"); count != 0 {
					t.Fatalf("commenter notice count = %d, want 0", count)
				}
			})

			t.Run("watcher_is_commenter", func(t *testing.T) {
				bbs_comment_seed_board(db, t, "13", "1", "post 7", "owner")
				bbs_comment_seed_parent(db, t, "13", "1", "2", "parent")
				seed_user_auth(db, t, "owner")
				seed_user_auth(db, t, "parent")
				seed_user_auth(db, t, "commenter")
				seed_user_watch(db, t, "bbs_watchlist", "owner", "13-1")
				seed_user_watch(db, t, "bbs_watchlist", "parent", "13-1")
				seed_user_watch(db, t, "bbs_watchlist", "commenter", "13-1")

				res := route.Api_bbs_w_comment_post(
					tool.Config{IP: "commenter", Session: &emptySession{}},
					"13",
					"1",
					"2",
					"reply data",
				)
				if res["response"] != "ok" {
					t.Fatalf("response = %v, want ok", res["response"])
				}

				if count := notice_count(db, "owner"); count != 1 {
					t.Fatalf("owner notice count = %d, want 1", count)
				}
				if count := notice_count(db, "parent"); count != 1 {
					t.Fatalf("parent notice count = %d, want 1", count)
				}
				if count := notice_count(db, "commenter"); count != 0 {
					t.Fatalf("commenter notice count = %d, want 0 (watcher and reply author are the same)", count)
				}
			})
		})
	}
}

func bbs_comment_seed_board(db *sql.DB, t *testing.T, set_id string, set_code string, title string, post_user string) {
	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed board %q: %v", q, err)
		}
	}
	exec(`insert into bbs_set (set_name, set_code, set_id, set_data) values (?, '', ?, ?)`, "bbs_name", set_id, "comment board")
	exec(`insert into bbs_set (set_name, set_code, set_id, set_data) values (?, '', ?, ?)`, "bbs_type", set_id, "comment")
	exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`, "title", set_code, set_id, title)
	exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`, "date", set_code, set_id, "2026-01-01 00:00:00")
	exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`, "user_id", set_code, set_id, post_user)
}

func bbs_comment_seed_parent(db *sql.DB, t *testing.T, set_id string, set_code string, comment_code string, user string) {
	base_id := set_id + "-" + set_code
	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed parent comment %q: %v", q, err)
		}
	}
	exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`, "comment", comment_code, base_id, "parent comment")
	exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`, "comment_date", comment_code, base_id, "2026-01-01 00:00:00")
	exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`, "comment_user_id", comment_code, base_id, user)
}
