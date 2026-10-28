package test

import (
	"database/sql"
	"os"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

func Test_thread_watchlist_legacy_notification(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso", "mysql"} {
		if backend == "mysql" && os.Getenv("MYSQL_AUDIT_DSN") == "" {
			continue
		}
		t.Run(backend, func(t *testing.T) {
			db := new_record_db(t, backend)
			defer tool.DB_close(db)

			thread_seed_bbs(db, t, "100", "thread 100")
			thread_seed_bbs(db, t, "200", "thread 200")

			t.Run("legacy_raw_thread_only_new_post", func(t *testing.T) {
				seed_user_watch(db, t, "thread_watchlist", "legacy_raw", "100")
				check_notification(t, db, "legacy_raw", 1, thread_notify_new_post(db))
			})

			t.Run("new_bbs_thread_new_post", func(t *testing.T) {
				seed_user_watch(db, t, "bbs_watchlist", "new_bbs", "-1-100")
				check_notification(t, db, "new_bbs", 1, thread_notify_new_post(db))
			})

			t.Run("legacy_prefix_thread_only_new_post", func(t *testing.T) {
				seed_user_watch(db, t, "thread_watchlist", "legacy_prefix", "-1-100")
				check_notification(t, db, "legacy_prefix", 1, thread_notify_new_post(db))
			})

			t.Run("duplicate_legacy_and_bbs_new_post", func(t *testing.T) {
				seed_user_watch(db, t, "bbs_watchlist", "dup", "-1-100")
				seed_user_watch(db, t, "thread_watchlist", "dup", "100")
				check_notification(t, db, "dup", 1, thread_notify_new_post(db))
			})

			t.Run("other_thread_watcher_not_notified", func(t *testing.T) {
				seed_user_watch(db, t, "thread_watchlist", "other_thread", "200")
				check_notification(t, db, "other_thread", 0, thread_notify_new_post(db))
			})

			t.Run("other_bbs_watcher_not_notified", func(t *testing.T) {
				seed_user_watch(db, t, "bbs_watchlist", "other_bbs", "5-7")
				check_notification(t, db, "other_bbs", 0, thread_notify_new_post(db))
			})

			t.Run("author_not_notified", func(t *testing.T) {
				seed_user_watch(db, t, "bbs_watchlist", "author", "-1-100")
				check_notification(t, db, "author", 0, thread_notify_author_new_post(db))
			})

			t.Run("parent_commenter_not_notified", func(t *testing.T) {
				seed_user_watch(db, t, "thread_watchlist", "parent", "100")
				check_notification(t, db, "parent", 0, thread_notify_parent_comment(db, "parent"))
			})

			t.Run("existing_comment_legacy_raw", func(t *testing.T) {
				seed_user_watch(db, t, "thread_watchlist", "legacy_comment", "100")
				check_notification(t, db, "legacy_comment", 1, thread_notify_comment(db))
			})

			t.Run("nested_comment_legacy_raw", func(t *testing.T) {
				seed_user_watch(db, t, "thread_watchlist", "legacy_nested", "100")
				check_notification(t, db, "legacy_nested", 1, thread_notify_nested_comment(db))
			})

			t.Run("multiple_watchers_exact_once", func(t *testing.T) {
				seed_user_watch(db, t, "thread_watchlist", "m1", "100")
				seed_user_watch(db, t, "bbs_watchlist", "m2", "-1-100")
				seed_user_watch(db, t, "bbs_watchlist", "m3", "-1-100")
				seed_user_watch(db, t, "thread_watchlist", "m3", "100")
				seed_user_watch(db, t, "thread_watchlist", "m4", "200")
				seed_user_watch(db, t, "bbs_watchlist", "m5", "5-7")
				seed_user_watch(db, t, "bbs_watchlist", "m6", "-1-100")

				thread_notify_new_post(db)()

				expectations := map[string]int{
					"m1": 1,
					"m2": 1,
					"m3": 1,
					"m4": 0,
					"m5": 0,
					"m6": 1,
				}
				for name, want := range expectations {
					got := notice_count(db, name)
					if got != want {
						t.Fatalf("notice count for %s = %d, want %d", name, got, want)
					}
				}
			})
		})
	}
}

func thread_notify_new_post(db *sql.DB) func() {
	return func() {
		route.Bbs_watch_notify(
			db,
			tool.Config{IP: "author", Session: &emptySession{}},
			"-1",
			"100",
			"1",
			"thread",
			"thread 100",
			"author",
			"",
		)
	}
}

func thread_notify_author_new_post(db *sql.DB) func() {
	return thread_notify_new_post(db)
}

func thread_notify_comment(db *sql.DB) func() {
	return func() {
		route.Bbs_watch_notify(
			db,
			tool.Config{IP: "commenter", Session: &emptySession{}},
			"-1",
			"100",
			"2",
			"thread",
			"thread 100",
			"commenter",
			"",
		)
	}
}

func thread_notify_nested_comment(db *sql.DB) func() {
	return func() {
		route.Bbs_watch_notify(
			db,
			tool.Config{IP: "nested_commenter", Session: &emptySession{}},
			"-1",
			"100",
			"3",
			"thread",
			"thread 100",
			"nested_commenter",
			"parent_commenter",
		)
	}
}

func thread_notify_parent_comment(db *sql.DB, parent_user string) func() {
	return func() {
		route.Bbs_watch_notify(
			db,
			tool.Config{IP: "commenter", Session: &emptySession{}},
			"-1",
			"100",
			"2",
			"thread",
			"thread 100",
			"commenter",
			parent_user,
		)
	}
}

func check_notification(t *testing.T, db *sql.DB, name string, want_delta int, action func()) {
	t.Helper()
	before := notice_count(db, name)
	action()
	after := notice_count(db, name)
	if after != before+want_delta {
		t.Fatalf("notice count for %s = %d, want delta %d from %d", name, after, want_delta, before)
	}
}

func notice_count(db *sql.DB, name string) int {
	var count int
	err := db.QueryRow("select count(*) from user_notice where name = ?", name).Scan(&count)
	if err != nil {
		return 0
	}
	return count
}
