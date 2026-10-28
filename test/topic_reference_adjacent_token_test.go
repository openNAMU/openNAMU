package test

import (
	"database/sql"
	"os"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

func Test_topic_reference_adjacent_token_parsing(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso", "mysql"} {
		if backend == "mysql" && os.Getenv("MYSQL_AUDIT_DSN") == "" {
			continue
		}
		t.Run(backend, func(t *testing.T) {
			bbsRef1 := `<a href="/bbs/w/1/100#1">#1</a>`
			bbsRef2 := `<a href="/bbs/w/1/100#2">#2</a>`
			bbsRef3 := `<a href="/bbs/w/1/100#3">#3</a>`
			bbsRef1_2 := `<a href="/bbs/w/1/2#1">#1-2</a>`
			bbsRef1_3 := `<a href="/bbs/w/1/3#1">#1-3</a>`

			threadRef0 := `<a href="/bbs/w/-1/10#0">#0</a>`
			threadRef1 := `<a href="/bbs/w/-1/10#1">#1</a>`
			threadRef2 := `<a href="/bbs/w/-1/10#2">#2</a>`
			threadRef3 := `<a href="/bbs/w/-1/10#3">#3</a>`
			threadRef1_2 := `<a href="/bbs/w/-1/2#1">#1-2</a>`

			alice := `<a href="/w/user:alice">@alice</a>`
			bob := `<a href="/w/user:bob">@bob</a>`
			charlie := `<a href="/w/user:charlie">@charlie</a>`

			t.Run("render_reference_adjacent_two", func(t *testing.T) {
				got := route.Render_topic_reference("#1 #2", "1", "1", "100", "bbs")
				want := bbsRef1 + " " + bbsRef2
				if got != want {
					t.Fatalf("render_reference_adjacent_two got=%q want=%q", got, want)
				}
			})

			t.Run("render_reference_adjacent_three", func(t *testing.T) {
				got := route.Render_topic_reference("#1 #2 #3", "1", "1", "100", "bbs")
				want := bbsRef1 + " " + bbsRef2 + " " + bbsRef3
				if got != want {
					t.Fatalf("render_reference_adjacent_three got=%q want=%q", got, want)
				}
			})

			t.Run("render_reference_two_spaces", func(t *testing.T) {
				got := route.Render_topic_reference("#1  #2", "1", "1", "100", "bbs")
				want := bbsRef1 + "  " + bbsRef2
				if got != want {
					t.Fatalf("render_reference_two_spaces got=%q want=%q", got, want)
				}
			})

			t.Run("render_reference_newline", func(t *testing.T) {
				got := route.Render_topic_reference("#1\n#2", "1", "1", "100", "bbs")
				want := bbsRef1 + "\n" + bbsRef2
				if got != want {
					t.Fatalf("render_reference_newline got=%q want=%q", got, want)
				}
			})

			t.Run("render_reference_start_end", func(t *testing.T) {
				for _, case_ := range []struct {
					data string
					want string
				}{
					{"#1", bbsRef1},
					{" #1", " " + bbsRef1},
					{"#1 ", bbsRef1 + " "},
					{"#1 #2", bbsRef1 + " " + bbsRef2},
					{"#1 #2 ", bbsRef1 + " " + bbsRef2 + " "},
					{" #1 #2", " " + bbsRef1 + " " + bbsRef2},
					{" #1 #2 ", " " + bbsRef1 + " " + bbsRef2 + " "},
				} {
					got := route.Render_topic_reference(case_.data, "1", "1", "100", "bbs")
					if got != case_.want {
						t.Fatalf("render_reference_start_end data=%q got=%q want=%q", case_.data, got, case_.want)
					}
				}
			})

			t.Run("render_reference_non_reference_unchanged", func(t *testing.T) {
				for _, case_ := range []struct {
					data string
					want string
				}{
					{"#1a", "#1a"},
					{"a#1", "a#1"},
					{"#1#2", "#1#2"},
					{"#1-", "#1-"},
					{"#1a #2", "#1a " + bbsRef2},
					{"#1-2abc #3", "#1-2abc " + bbsRef3},
					{"foo bar", "foo bar"},
					{"text with no tokens", "text with no tokens"},
				} {
					got := route.Render_topic_reference(case_.data, "1", "1", "100", "bbs")
					if got != case_.want {
						t.Fatalf("render_reference_non_reference_unchanged data=%q got=%q want=%q", case_.data, got, case_.want)
					}
				}
			})

			t.Run("render_reference_hyphen_code_adjacent", func(t *testing.T) {
				got := route.Render_topic_reference("#1-2 #1-3", "1", "1", "100", "bbs")
				want := bbsRef1_2 + " " + bbsRef1_3
				if got != want {
					t.Fatalf("render_reference_hyphen_code_adjacent got=%q want=%q", got, want)
				}
			})

			t.Run("render_thread_reference_adjacent", func(t *testing.T) {
				got := route.Render_topic_reference("#1 #2 #0", "10", "", "", "thread")
				want := threadRef1 + " " + threadRef2 + " " + threadRef0
				if got != want {
					t.Fatalf("render_thread_reference_adjacent got=%q want=%q", got, want)
				}
			})

			t.Run("render_thread_reference_three", func(t *testing.T) {
				got := route.Render_topic_reference("#1 #2 #3", "10", "", "", "thread")
				want := threadRef1 + " " + threadRef2 + " " + threadRef3
				if got != want {
					t.Fatalf("render_thread_reference_three got=%q want=%q", got, want)
				}
			})

			t.Run("render_thread_reference_post_author", func(t *testing.T) {
				got := route.Render_topic_reference("#0 #1 #2", "10", "", "", "thread")
				want := threadRef0 + " " + threadRef1 + " " + threadRef2
				if got != want {
					t.Fatalf("render_thread_reference_post_author got=%q want=%q", got, want)
				}
			})

			t.Run("render_thread_reference_hyphen_code", func(t *testing.T) {
				got := route.Render_topic_reference("#1-2", "10", "", "", "thread")
				if got != threadRef1_2 {
					t.Fatalf("render_thread_reference_hyphen_code got=%q want=%q", got, threadRef1_2)
				}
			})

			t.Run("render_mention_adjacent_two", func(t *testing.T) {
				got := route.Render_topic_reference("@alice @bob", "1", "1", "100", "bbs")
				want := alice + " " + bob
				if got != want {
					t.Fatalf("render_mention_adjacent_two got=%q want=%q", got, want)
				}
			})

			t.Run("render_mention_adjacent_three", func(t *testing.T) {
				got := route.Render_topic_reference("@alice @bob @charlie", "1", "1", "100", "bbs")
				want := alice + " " + bob + " " + charlie
				if got != want {
					t.Fatalf("render_mention_adjacent_three got=%q want=%q", got, want)
				}
			})

			t.Run("render_mention_newline", func(t *testing.T) {
				got := route.Render_topic_reference("@alice\n@bob", "1", "1", "100", "bbs")
				want := alice + "\n" + bob
				if got != want {
					t.Fatalf("render_mention_newline got=%q want=%q", got, want)
				}
			})

			t.Run("render_mention_start_end", func(t *testing.T) {
				for _, case_ := range []struct {
					data string
					want string
				}{
					{"@alice", alice},
					{" @alice", " " + alice},
					{"@alice ", alice + " "},
					{"@alice @bob", alice + " " + bob},
					{"@alice @bob ", alice + " " + bob + " "},
					{" @alice @bob", " " + alice + " " + bob},
				} {
					got := route.Render_topic_reference(case_.data, "1", "1", "100", "bbs")
					if got != case_.want {
						t.Fatalf("render_mention_start_end data=%q got=%q want=%q", case_.data, got, case_.want)
					}
				}
			})

			t.Run("render_mention_non_reference_unchanged", func(t *testing.T) {
				got := route.Render_topic_reference("a@alice", "1", "1", "100", "bbs")
				if got != "a@alice" {
					t.Fatalf("render_mention_non_reference_unchanged got=%q want=%q", got, "a@alice")
				}
			})

			t.Run("render_legacy_thread_tag_then_reference", func(t *testing.T) {
				got := route.Render_topic_reference("<topic_a>#2</topic_a> #3", "1", "1", "100", "bbs")
				want := `<a href="#2">#2</a> ` + bbsRef3
				if got != want {
					t.Fatalf("render_legacy_thread_tag_then_reference got=%q want=%q", got, want)
				}
			})

			notify_case := func(seed func(db *sql.DB, t *testing.T), data, topic_num, set_id, do_type string, wants map[string]int) {
				db := new_record_db(t, backend)
				defer tool.DB_close(db)

				seed(db, t)
				route.Topic_reference_notify(
					db,
					tool.Config{IP: "commenter", Session: &emptySession{}},
					data,
					"2",
					topic_num,
					set_id,
					"title",
					"bbs",
					do_type,
				)

				assert_notice_counts(db, t, wants)
			}

			notify_case(func(db *sql.DB, tt *testing.T) {
				reference_seed_bbs_comment_user(db, tt, "100", "1", "1", "ref_user_1")
				reference_seed_bbs_comment_user(db, tt, "100", "1", "2", "ref_user_2")
			}, "#1 #2", "1", "100", "bbs", map[string]int{
				"ref_user_1": 1,
				"ref_user_2": 1,
				"commenter":  0,
			})

			notify_case(func(db *sql.DB, tt *testing.T) {
				reference_seed_bbs_comment_user(db, tt, "100", "1", "1", "ref_user_1")
				reference_seed_bbs_comment_user(db, tt, "100", "1", "2", "ref_user_2")
				reference_seed_bbs_comment_user(db, tt, "100", "1", "3", "ref_user_3")
			}, "#1 #2 #3", "1", "100", "bbs", map[string]int{
				"ref_user_1": 1,
				"ref_user_2": 1,
				"ref_user_3": 1,
			})

			notify_case(func(db *sql.DB, tt *testing.T) {
				reference_seed_bbs_comment_user(db, tt, "100", "1", "1", "ref_user_1")
				reference_seed_bbs_comment_user(db, tt, "100", "1", "2", "ref_user_2")
			}, "#1  #2", "1", "100", "bbs", map[string]int{
				"ref_user_1": 1,
				"ref_user_2": 1,
			})

			notify_case(func(db *sql.DB, tt *testing.T) {
				reference_seed_bbs_comment_user(db, tt, "100", "1", "1", "ref_user_1")
				reference_seed_bbs_comment_user(db, tt, "100", "1", "2", "ref_user_2")
			}, "#1\n#2", "1", "100", "bbs", map[string]int{
				"ref_user_1": 1,
				"ref_user_2": 1,
			})

			notify_case(func(db *sql.DB, tt *testing.T) {
				reference_seed_bbs_comment_user(db, tt, "100", "1", "1", "ref_user_1")
				reference_seed_bbs_comment_user(db, tt, "100", "1", "2", "ref_user_2")
			}, " #1 #2 ", "1", "100", "bbs", map[string]int{
				"ref_user_1": 1,
				"ref_user_2": 1,
			})

			notify_case(func(db *sql.DB, tt *testing.T) {
				reference_seed_bbs_comment_user(db, tt, "100", "1", "1", "same_ref_user")
			}, "#1 first #1 second", "1", "100", "bbs", map[string]int{
				"same_ref_user": 1,
			})

			notify_case(func(db *sql.DB, tt *testing.T) {
				reference_seed_bbs_comment_user(db, tt, "100", "1", "1", "ref_user_a")
				reference_seed_bbs_comment_user(db, tt, "100", "1", "2", "ref_user_b")
				reference_seed_history(db, tt, "at_user_c")
			}, "#1 #2 @at_user_c", "1", "100", "bbs", map[string]int{
				"ref_user_a": 1,
				"ref_user_b": 1,
				"at_user_c":  1,
			})

			notify_case(func(db *sql.DB, tt *testing.T) {
				reference_seed_bbs_comment_user(db, tt, "100", "1", "1", "mixed_user")
				reference_seed_history(db, tt, "mixed_user")
			}, "#1 @mixed_user", "1", "100", "bbs", map[string]int{
				"mixed_user": 1,
			})

			notify_case(func(db *sql.DB, tt *testing.T) {
				reference_seed_history(db, tt, "alice")
				reference_seed_history(db, tt, "bob")
			}, "@alice @bob", "1", "100", "bbs", map[string]int{
				"alice": 1,
				"bob":   1,
			})

			notify_case(func(db *sql.DB, tt *testing.T) {
				reference_seed_history(db, tt, "alice")
				reference_seed_history(db, tt, "bob")
				reference_seed_history(db, tt, "charlie")
			}, "@alice @bob @charlie", "1", "100", "bbs", map[string]int{
				"alice":   1,
				"bob":     1,
				"charlie": 1,
			})

			notify_case(func(db *sql.DB, tt *testing.T) {
				reference_seed_history(db, tt, "alice")
				reference_seed_history(db, tt, "bob")
			}, "@alice\n@bob", "1", "100", "bbs", map[string]int{
				"alice": 1,
				"bob":   1,
			})

			notify_case(func(db *sql.DB, tt *testing.T) {
				reference_seed_bbs_comment_user(db, tt, "100", "2", "1", "ref_code_2")
				reference_seed_bbs_comment_user(db, tt, "100", "3", "1", "ref_code_3")
			}, "#1-2 #1-3", "1", "100", "bbs", map[string]int{
				"ref_code_2": 1,
				"ref_code_3": 1,
			})

			notify_case(func(db *sql.DB, tt *testing.T) {
				reference_seed_bbs_comment_user(db, tt, "100", "1", "2", "ref_user")
				reference_seed_history(db, tt, "alice")
			}, "#1a #2 @unknown @alice", "1", "100", "bbs", map[string]int{
				"ref_user": 1,
				"alice":    1,
				"unknown":  0,
			})

			notify_case(func(db *sql.DB, tt *testing.T) {
				reference_seed_thread_comment_user(db, tt, "10", "1", "thread_ref_1")
				reference_seed_thread_comment_user(db, tt, "10", "2", "thread_ref_2")
				reference_seed_thread_comment_user(db, tt, "10", "3", "thread_ref_3")
			}, "#1 #2 #3", "10", "-1", "thread", map[string]int{
				"thread_ref_1": 1,
				"thread_ref_2": 1,
				"thread_ref_3": 1,
			})

			notify_case(func(db *sql.DB, tt *testing.T) {
				reference_seed_thread_comment_user(db, tt, "10", "1", "thread_ref_1")
				reference_seed_thread_comment_user(db, tt, "10", "2", "thread_ref_2")
			}, " #1 #2 ", "10", "-1", "thread", map[string]int{
				"thread_ref_1": 1,
				"thread_ref_2": 1,
			})

			notify_case(func(db *sql.DB, tt *testing.T) {
				reference_seed_thread_post_user(db, tt, "10", "thread_post_user")
				reference_seed_thread_comment_user(db, tt, "10", "1", "thread_comment_1")
				reference_seed_thread_comment_user(db, tt, "10", "2", "thread_comment_2")
			}, "#0 #1 #2", "10", "-1", "thread", map[string]int{
				"thread_post_user":   1,
				"thread_comment_1":   1,
				"thread_comment_2":   1,
			})
		})
	}
}

func assert_notice_counts(db *sql.DB, t *testing.T, counts map[string]int) {
	for name, want := range counts {
		if got := notice_count(db, name); got != want {
			t.Fatalf("%s notice count = %d, want %d", name, got, want)
		}
	}
}
