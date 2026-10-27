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

type record_pagination_case struct {
	n    int
	name string
}

func seed_record_pagination_n(db *sql.DB, t *testing.T, n int) {
	for i := 1; i <= n; i++ {
		set_id := fmt.Sprintf("%d", i)
		record_test_seed_post(db, t, set_id, "3", "alice", fmt.Sprintf("Post %d", i), "2026-06-01 00:00:00")
		record_test_seed_comment(db, t, fmt.Sprintf("%d-3", i), "1", "2026-06-02 00:00:00", "alice")
	}
}

func seed_record_board_posts(db *sql.DB, t *testing.T, set_id string, n int) {
	for i := 1; i <= n; i++ {
		set_code := fmt.Sprintf("%d", i)
		record_test_seed_post(db, t, set_id, set_code, "alice", fmt.Sprintf("Board Post %d", i), "2026-06-03 00:00:00")
	}
}

func seed_record_board_comments(db *sql.DB, t *testing.T, set_id string, n int) {
	for i := 1; i <= n; i++ {
		set_code := fmt.Sprintf("%d", i)
		record_test_seed_comment(db, t, set_id+"-3", set_code, "2026-06-03 00:00:00", "alice")
	}
}

func Test_record_bbs_pagination_route(t *testing.T) {
	backends := []string{"sqlite", "turso"}
	if os.Getenv("MYSQL_AUDIT_DSN") != "" {
		backends = append(backends, "mysql")
	}

	for _, backend := range backends {
		t.Run(backend, func(t *testing.T) {
			cases := []record_pagination_case{
				{51, "page1_to_page2_51"},
				{50, "no_next_50"},
			}
			for _, case_item := range cases {
				t.Run("record_bbs_" + case_item.name, func(t *testing.T) {
					db := new_record_db(t, backend)
					defer tool.DB_close(db)

					seed_record_access_auth(db, t)
					seed_record_pagination_n(db, t, case_item.n)

					config_self := tool.Config{IP: "alice", Session: &emptySession{}}
					config_hidel := tool.Config{IP: "hidel_user", Session: &emptySession{}}
					config_regular := tool.Config{IP: "regular", Session: &emptySession{}}

					page1 := route.View_record_bbs(config_self, "alice", "1")
					page2 := route.View_record_bbs(config_self, "alice", "2")

					if case_item.n == 51 {
						if !strings.Contains(page1, `href="/record_bbs/alice/2"`) {
							t.Fatalf("bbs page1 missing next /record_bbs/alice/2")
						}
						if strings.Contains(page1, `href="/record/bbs/alice/2"`) {
							t.Fatalf("bbs page1 legacy next /record/bbs/alice/2 remains")
						}
						if !strings.Contains(page1, "/record_bbs_in/51/alice") {
							t.Fatalf("bbs page1 missing first board 51")
						}
						if !strings.Contains(page1, "/record_bbs_in/2/alice") {
							t.Fatalf("bbs page1 missing last board 2")
						}
						if strings.Contains(page1, "/record_bbs_in/1/alice") {
							t.Fatalf("bbs page1 incorrectly contains board 1")
						}

						if !strings.Contains(page2, `href="/record_bbs/alice/1"`) {
							t.Fatalf("bbs page2 missing previous /record_bbs/alice/1")
						}
						if strings.Contains(page2, `href="/record_bbs/alice/3"`) {
							t.Fatalf("bbs page2 incorrectly contains next /record_bbs/alice/3")
						}
						if !strings.Contains(page2, "/record_bbs_in/1/alice") {
							t.Fatalf("bbs page2 missing board 1")
						}
					} else {
						if !strings.Contains(page1, "/record_bbs_in/50/alice") {
							t.Fatalf("bbs exact50 page1 missing board 50")
						}
						if !strings.Contains(page1, "/record_bbs_in/1/alice") {
							t.Fatalf("bbs exact50 page1 missing board 1")
						}
						if strings.Contains(page1, `href="/record_bbs/alice/2"`) {
							t.Fatalf("bbs exact50 page1 shows unnecessary next /record_bbs/alice/2")
						}
						if strings.Contains(page1, `href="/record/bbs/alice/2"`) {
							t.Fatalf("bbs exact50 page1 legacy next /record/bbs/alice/2 remains")
						}
					}

					hidel_page1 := route.View_record_bbs(config_hidel, "alice", "1")
					if case_item.n == 51 {
						if !strings.Contains(hidel_page1, `href="/record_bbs/alice/2"`) {
							t.Fatalf("hidel bbs page1 missing next")
						}
					} else {
						if strings.Contains(hidel_page1, `href="/record_bbs/alice/2"`) {
							t.Fatalf("hidel bbs exact50 page1 shows unnecessary next")
						}
					}

					regular_page1 := route.View_record_bbs(config_regular, "alice", "1")
					if case_item.n == 51 {
						if strings.Contains(regular_page1, `href="/record_bbs/alice/2"`) {
							t.Fatalf("regular bbs page1 exposes next link")
						}
						if strings.Contains(regular_page1, "/record_bbs_in/51/alice") {
							t.Fatalf("regular bbs page1 exposes board 51")
						}
					}
					if strings.Contains(regular_page1, "/record_bbs_in/50/alice") {
						t.Fatalf("regular bbs exact50 page1 exposes board 50")
					}
				})
			}

			for _, case_item := range cases {
				t.Run("record_bbs_comment_" + case_item.name, func(t *testing.T) {
					db := new_record_db(t, backend)
					defer tool.DB_close(db)

					seed_record_access_auth(db, t)
					seed_record_pagination_n(db, t, case_item.n)

					config_self := tool.Config{IP: "alice", Session: &emptySession{}}
					config_hidel := tool.Config{IP: "hidel_user", Session: &emptySession{}}
					config_regular := tool.Config{IP: "regular", Session: &emptySession{}}

					page1 := route.View_record_bbs_comment(config_self, "alice", "1")
					page2 := route.View_record_bbs_comment(config_self, "alice", "2")

					if case_item.n == 51 {
						if !strings.Contains(page1, `href="/record_bbs_comment/alice/2"`) {
							t.Fatalf("bbs comment page1 missing next /record_bbs_comment/alice/2")
						}
						if strings.Contains(page1, `href="/record/bbs_comment/alice/2"`) {
							t.Fatalf("bbs comment page1 legacy next /record/bbs_comment/alice/2 remains")
						}
						if !strings.Contains(page1, "/record_bbs_comment_in/51/alice") {
							t.Fatalf("bbs comment page1 missing first board 51")
						}
						if !strings.Contains(page1, "/record_bbs_comment_in/2/alice") {
							t.Fatalf("bbs comment page1 missing last board 2")
						}
						if strings.Contains(page1, "/record_bbs_comment_in/1/alice") {
							t.Fatalf("bbs comment page1 incorrectly contains board 1")
						}

						if !strings.Contains(page2, `href="/record_bbs_comment/alice/1"`) {
							t.Fatalf("bbs comment page2 missing previous /record_bbs_comment/alice/1")
						}
						if strings.Contains(page2, `href="/record_bbs_comment/alice/3"`) {
							t.Fatalf("bbs comment page2 incorrectly contains next /record_bbs_comment/alice/3")
						}
						if !strings.Contains(page2, "/record_bbs_comment_in/1/alice") {
							t.Fatalf("bbs comment page2 missing board 1")
						}
					} else {
						if !strings.Contains(page1, "/record_bbs_comment_in/50/alice") {
							t.Fatalf("bbs comment exact50 page1 missing board 50")
						}
						if !strings.Contains(page1, "/record_bbs_comment_in/1/alice") {
							t.Fatalf("bbs comment exact50 page1 missing board 1")
						}
						if strings.Contains(page1, `href="/record_bbs_comment/alice/2"`) {
							t.Fatalf("bbs comment exact50 page1 shows unnecessary next /record_bbs_comment/alice/2")
						}
						if strings.Contains(page1, `href="/record/bbs_comment/alice/2"`) {
							t.Fatalf("bbs comment exact50 page1 legacy next /record/bbs_comment/alice/2 remains")
						}
					}

					hidel_page1 := route.View_record_bbs_comment(config_hidel, "alice", "1")
					if case_item.n == 51 {
						if !strings.Contains(hidel_page1, `href="/record_bbs_comment/alice/2"`) {
							t.Fatalf("hidel bbs comment page1 missing next")
						}
					} else {
						if strings.Contains(hidel_page1, `href="/record_bbs_comment/alice/2"`) {
							t.Fatalf("hidel bbs comment exact50 page1 shows unnecessary next")
						}
					}

					regular_page1 := route.View_record_bbs_comment(config_regular, "alice", "1")
					if case_item.n == 51 {
						if strings.Contains(regular_page1, `href="/record_bbs_comment/alice/2"`) {
							t.Fatalf("regular bbs comment page1 exposes next link")
						}
						if strings.Contains(regular_page1, "/record_bbs_comment_in/51/alice") {
							t.Fatalf("regular bbs comment page1 exposes board 51")
						}
					}
					if strings.Contains(regular_page1, "/record_bbs_comment_in/50/alice") {
						t.Fatalf("regular bbs comment exact50 page1 exposes board 50")
					}
				})
			}

			for _, case_item := range cases {
				t.Run("record_bbs_in_" + case_item.name, func(t *testing.T) {
					db := new_record_db(t, backend)
					defer tool.DB_close(db)

					seed_record_access_auth(db, t)
					seed_record_board_posts(db, t, "5", case_item.n)

					config_self := tool.Config{IP: "alice", Session: &emptySession{}}
					config_hidel := tool.Config{IP: "hidel_user", Session: &emptySession{}}
					config_regular := tool.Config{IP: "regular", Session: &emptySession{}}

					page1 := route.View_record_bbs_in(config_self, "alice", "5", "1")
					page2 := route.View_record_bbs_in(config_self, "alice", "5", "2")

					if case_item.n == 51 {
						if !strings.Contains(page1, `href="/record_bbs_in/5/alice/2"`) {
							t.Fatalf("bbs_in page1 missing next /record_bbs_in/5/alice/2")
						}
						if !strings.Contains(page1, "/bbs/w/5/51") {
							t.Fatalf("bbs_in page1 missing post 51")
						}
						if !strings.Contains(page1, "/bbs/w/5/2") {
							t.Fatalf("bbs_in page1 missing post 2")
						}
						if strings.Contains(page1, `"/bbs/w/5/1">`) {
							t.Fatalf("bbs_in page1 incorrectly contains post 1")
						}

						if !strings.Contains(page2, `href="/record_bbs_in/5/alice/1"`) {
							t.Fatalf("bbs_in page2 missing previous /record_bbs_in/5/alice/1")
						}
						if strings.Contains(page2, `href="/record_bbs_in/5/alice/3"`) {
							t.Fatalf("bbs_in page2 incorrectly contains next /record_bbs_in/5/alice/3")
						}
						if !strings.Contains(page2, `"/bbs/w/5/1">`) {
							t.Fatalf("bbs_in page2 missing post 1")
						}
					} else {
						if !strings.Contains(page1, `"/bbs/w/5/50">`) {
							t.Fatalf("bbs_in exact50 page1 missing post 50")
						}
						if !strings.Contains(page1, `"/bbs/w/5/1">`) {
							t.Fatalf("bbs_in exact50 page1 missing post 1")
						}
						if strings.Contains(page1, `href="/record_bbs_in/5/alice/2"`) {
							t.Fatalf("bbs_in exact50 page1 shows unnecessary next")
						}
					}

					hidel_page1 := route.View_record_bbs_in(config_hidel, "alice", "5", "1")
					if case_item.n == 51 {
						if !strings.Contains(hidel_page1, `href="/record_bbs_in/5/alice/2"`) {
							t.Fatalf("hidel bbs_in page1 missing next")
						}
					} else {
						if strings.Contains(hidel_page1, `href="/record_bbs_in/5/alice/2"`) {
							t.Fatalf("hidel bbs_in exact50 page1 shows unnecessary next")
						}
					}

					regular_page1 := route.View_record_bbs_in(config_regular, "alice", "5", "1")
					if case_item.n == 51 {
						if strings.Contains(regular_page1, `href="/record_bbs_in/5/alice/2"`) {
							t.Fatalf("regular bbs_in page1 exposes next link")
						}
						if strings.Contains(regular_page1, "/bbs/w/5/51") {
							t.Fatalf("regular bbs_in page1 exposes post 51")
						}
					}
					if strings.Contains(regular_page1, `"/bbs/w/5/50">`) {
						t.Fatalf("regular bbs_in exact50 page1 exposes post 50")
					}
				})
			}

			for _, case_item := range cases {
				t.Run("record_bbs_comment_in_" + case_item.name, func(t *testing.T) {
					db := new_record_db(t, backend)
					defer tool.DB_close(db)

					seed_record_access_auth(db, t)
					seed_record_board_comments(db, t, "5", case_item.n)

					config_self := tool.Config{IP: "alice", Session: &emptySession{}}
					config_hidel := tool.Config{IP: "hidel_user", Session: &emptySession{}}
					config_regular := tool.Config{IP: "regular", Session: &emptySession{}}

					page1 := route.View_record_bbs_comment_in(config_self, "alice", "5", "1")
					page2 := route.View_record_bbs_comment_in(config_self, "alice", "5", "2")

					if case_item.n == 51 {
						if !strings.Contains(page1, `href="/record_bbs_comment_in/5/alice/2"`) {
							t.Fatalf("comment_in page1 missing next /record_bbs_comment_in/5/alice/2")
						}
						if !strings.Contains(page2, `href="/record_bbs_comment_in/5/alice/1"`) {
							t.Fatalf("comment_in page2 missing previous /record_bbs_comment_in/5/alice/1")
						}
						if strings.Contains(page2, `href="/record_bbs_comment_in/5/alice/3"`) {
							t.Fatalf("comment_in page2 incorrectly contains next /record_bbs_comment_in/5/alice/3")
						}
					} else {
						if strings.Contains(page1, `href="/record_bbs_comment_in/5/alice/2"`) {
							t.Fatalf("comment_in exact50 page1 shows unnecessary next")
						}
					}

					hidel_page1 := route.View_record_bbs_comment_in(config_hidel, "alice", "5", "1")
					if case_item.n == 51 {
						if !strings.Contains(hidel_page1, `href="/record_bbs_comment_in/5/alice/2"`) {
							t.Fatalf("hidel comment_in page1 missing next")
						}
					} else {
						if strings.Contains(hidel_page1, `href="/record_bbs_comment_in/5/alice/2"`) {
							t.Fatalf("hidel comment_in exact50 page1 shows unnecessary next")
						}
					}

					regular_page1 := route.View_record_bbs_comment_in(config_regular, "alice", "5", "1")
					if case_item.n == 51 {
						if strings.Contains(regular_page1, `href="/record_bbs_comment_in/5/alice/2"`) {
							t.Fatalf("regular comment_in page1 exposes next link")
						}
						if strings.Contains(regular_page1, "/bbs/w/5/3#1") {
							t.Fatalf("regular comment_in page1 exposes comment anchor")
						}
					}
					if strings.Contains(regular_page1, "/bbs/w/5/3#50") {
						t.Fatalf("regular comment_in exact50 page1 exposes comment anchor")
					}
				})
			}
		})
	}
}
