package test

import (
	"database/sql"
	"os"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

func Test_topic_reference_duplicate_notification(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso", "mysql"} {
		if backend == "mysql" && os.Getenv("MYSQL_AUDIT_DSN") == "" {
			continue
		}
		t.Run(backend, func(t *testing.T) {
			db := new_record_db(t, backend)
			defer tool.DB_close(db)

			t.Run("same_at_mention_multiple_times", func(t *testing.T) {
				reference_seed_history(db, t, "same_at_user")

				route.Topic_reference_notify(
					db,
					tool.Config{IP: "commenter", Session: &emptySession{}},
					"@same_at_user one @same_at_user two",
					"2",
					"1",
					"5",
					"board",
					"title",
					"bbs",
				)

				if count := notice_count(db, "same_at_user"); count != 1 {
					t.Fatalf("same_at_user notice count = %d, want 1", count)
				}
				if count := notice_count(db, "commenter"); count != 0 {
					t.Fatalf("commenter notice count = %d, want 0", count)
				}
			})

			t.Run("same_bbs_reference_multiple_times", func(t *testing.T) {
				reference_seed_bbs_comment_user(db, t, "100", "1", "1", "same_ref_user")

				route.Topic_reference_notify(
					db,
					tool.Config{IP: "commenter", Session: &emptySession{}},
					"#1 first #1 second",
					"2",
					"1",
					"100",
					"board",
					"title",
					"bbs",
				)

				if count := notice_count(db, "same_ref_user"); count != 1 {
					t.Fatalf("same_ref_user notice count = %d, want 1", count)
				}
			})

			t.Run("same_bbs_post_reference_multiple_times", func(t *testing.T) {
				reference_seed_bbs_post_user(db, t, "101", "1", "same_post_user")

				route.Topic_reference_notify(
					db,
					tool.Config{IP: "commenter", Session: &emptySession{}},
					"#0 first #0 second",
					"2",
					"1",
					"101",
					"board",
					"title",
					"bbs",
				)

				if count := notice_count(db, "same_post_user"); count != 1 {
					t.Fatalf("same_post_user notice count = %d, want 1", count)
				}
			})

			t.Run("reference_and_at_resolve_to_same_user", func(t *testing.T) {
				reference_seed_bbs_comment_user(db, t, "102", "1", "1", "mixed_user")
				reference_seed_history(db, t, "mixed_user")

				route.Topic_reference_notify(
					db,
					tool.Config{IP: "commenter", Session: &emptySession{}},
					"#1 @mixed_user",
					"2",
					"1",
					"102",
					"board",
					"title",
					"bbs",
				)

				if count := notice_count(db, "mixed_user"); count != 1 {
					t.Fatalf("mixed_user notice count = %d, want 1", count)
				}
			})

			t.Run("references_resolve_to_different_users", func(t *testing.T) {
				reference_seed_bbs_comment_user(db, t, "103", "1", "1", "ref_user_a")
				reference_seed_bbs_comment_user(db, t, "103", "1", "2", "ref_user_b")
				reference_seed_history(db, t, "at_user_c")

				route.Topic_reference_notify(
					db,
					tool.Config{IP: "commenter", Session: &emptySession{}},
					"#1 first #2 second @at_user_c",
					"2",
					"1",
					"103",
					"board",
					"title",
					"bbs",
				)

				for name, want := range map[string]int{
					"ref_user_a": 1,
					"ref_user_b": 1,
					"at_user_c":  1,
				} {
					if count := notice_count(db, name); count != want {
						t.Fatalf("%s notice count = %d, want %d", name, count, want)
					}
				}
			})

			t.Run("commenter_mentions_themselves", func(t *testing.T) {
				reference_seed_history(db, t, "self_user")
				reference_seed_bbs_comment_user(db, t, "104", "1", "1", "self_user")

				route.Topic_reference_notify(
					db,
					tool.Config{IP: "self_user", Session: &emptySession{}},
					"@self_user #1",
					"2",
					"1",
					"104",
					"board",
					"title",
					"bbs",
				)

				if count := notice_count(db, "self_user"); count != 0 {
					t.Fatalf("self_user notice count = %d, want 0", count)
				}
			})

			t.Run("bbs_reference_and_thread_reference_dont_mix", func(t *testing.T) {
				reference_seed_bbs_comment_user(db, t, "105", "1", "1", "bbs_user")
				reference_seed_thread_comment_user(db, t, "1", "2", "thread_user")

				route.Topic_reference_notify(
					db,
					tool.Config{IP: "commenter", Session: &emptySession{}},
					"#1",
					"2",
					"1",
					"105",
					"board",
					"title",
					"bbs",
				)
				route.Topic_reference_notify(
					db,
					tool.Config{IP: "commenter", Session: &emptySession{}},
					"#2",
					"2",
					"1",
					"-1",
					"thread",
					"title",
					"thread",
				)

				if count := notice_count(db, "bbs_user"); count != 1 {
					t.Fatalf("bbs_user notice count = %d, want 1", count)
				}
				if count := notice_count(db, "thread_user"); count != 1 {
					t.Fatalf("thread_user notice count = %d, want 1", count)
				}
			})

			t.Run("thread_same_reference_multiple_times", func(t *testing.T) {
				reference_seed_thread_comment_user(db, t, "2", "2", "same_thread_user")

				route.Topic_reference_notify(
					db,
					tool.Config{IP: "commenter", Session: &emptySession{}},
					"#2 first #2 second",
					"2",
					"2",
					"-1",
					"thread",
					"title",
					"thread",
				)

				if count := notice_count(db, "same_thread_user"); count != 1 {
					t.Fatalf("same_thread_user notice count = %d, want 1", count)
				}
			})

			t.Run("thread_post_author_reference_multiple_times", func(t *testing.T) {
				reference_seed_thread_post_user(db, t, "3", "same_thread_post_user")

				route.Topic_reference_notify(
					db,
					tool.Config{IP: "commenter", Session: &emptySession{}},
					"#0 first #0 second",
					"2",
					"3",
					"-1",
					"thread",
					"title",
					"thread",
				)

				if count := notice_count(db, "same_thread_post_user"); count != 1 {
					t.Fatalf("same_thread_post_user notice count = %d, want 1", count)
				}
			})

			t.Run("thread_reference_and_at_resolve_to_same_user", func(t *testing.T) {
				reference_seed_thread_post_user(db, t, "4", "same_thread_mixed_user")
				reference_seed_history(db, t, "same_thread_mixed_user")

				route.Topic_reference_notify(
					db,
					tool.Config{IP: "commenter", Session: &emptySession{}},
					"#0 @same_thread_mixed_user",
					"2",
					"4",
					"-1",
					"thread",
					"title",
					"thread",
				)

				if count := notice_count(db, "same_thread_mixed_user"); count != 1 {
					t.Fatalf("same_thread_mixed_user notice count = %d, want 1", count)
				}
			})

			t.Run("thread_commenter_mentions_themselves", func(t *testing.T) {
				reference_seed_thread_comment_user(db, t, "5", "2", "self_thread_user")

				route.Topic_reference_notify(
					db,
					tool.Config{IP: "self_thread_user", Session: &emptySession{}},
					"#2 @self_thread_user",
					"2",
					"5",
					"-1",
					"thread",
					"title",
					"thread",
				)

				if count := notice_count(db, "self_thread_user"); count != 0 {
					t.Fatalf("self_thread_user notice count = %d, want 0", count)
				}
			})
		})
	}
}

func reference_seed_history(db *sql.DB, t *testing.T, user string) {
	if _, err := db.Exec(
		"insert into history (title, data, date, ip, send, leng, hide, type) values (?, ?, ?, ?, '', '', '', 'document')",
		user,
		"history seed",
		"2026-01-01 00:00:00",
		user,
	); err != nil {
		t.Fatalf("seed history %s: %v", user, err)
	}
}

func reference_seed_bbs_post_user(db *sql.DB, t *testing.T, set_id string, topic_num string, user string) {
	if _, err := db.Exec(
		"insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)",
		"user_id",
		topic_num,
		set_id,
		user,
	); err != nil {
		t.Fatalf("seed bbs post user %s: %v", user, err)
	}
}

func reference_seed_bbs_comment_user(db *sql.DB, t *testing.T, set_id string, topic_num string, comment_code string, user string) {
	comment_set_id := set_id + "-" + topic_num
	if _, err := db.Exec(
		"insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)",
		"comment_user_id",
		comment_code,
		comment_set_id,
		user,
	); err != nil {
		t.Fatalf("seed bbs comment user %s: %v", user, err)
	}
}

func reference_seed_thread_post_user(db *sql.DB, t *testing.T, topic_num string, user string) {
	if _, err := db.Exec(
		"insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)",
		"user_id",
		topic_num,
		"-1",
		user,
	); err != nil {
		t.Fatalf("seed thread post user %s: %v", user, err)
	}
}

func reference_seed_thread_comment_user(db *sql.DB, t *testing.T, topic_num string, comment_code string, user string) {
	comment_set_id := "-1-" + topic_num
	if _, err := db.Exec(
		"insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)",
		"comment_user_id",
		comment_code,
		comment_set_id,
		user,
	); err != nil {
		t.Fatalf("seed thread comment user %s: %v", user, err)
	}
}
