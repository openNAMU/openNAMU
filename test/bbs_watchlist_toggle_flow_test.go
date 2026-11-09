package test

import (
	"database/sql"
	"os"
	"sync"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

func Test_Api_bbs_watch_post_non_thread_toggle_and_authorization(t *testing.T) {
	backends := []string{"sqlite", "turso"}
	if os.Getenv("MYSQL_AUDIT_DSN") != "" {
		backends = append(backends, "mysql")
	}

	for _, backend := range backends {
		t.Run(backend, func(t *testing.T) {
			db := new_record_db(t, backend)
			Seed_watchlist_test_board(db, t)
			owner_config := tool.Config{IP: "admin_ip", Session: &emptySession{}}

			for _, expected_count := range []int{1, 0} {
				result := route.Api_bbs_watch_post(owner_config, "2", "1")
				if result["response"] != "ok" {
					t.Fatalf("BBS watch toggle = %#v, want ok", result)
				}
				if count := Bbs_watchlist_row_count(db, t, "admin_ip", "2-1"); count != expected_count {
					t.Fatalf("BBS watch row count = %d, want %d", count, expected_count)
				}
			}

			missing_post := route.Api_bbs_watch_post(owner_config, "2", "missing")
			if missing_post["response"] != "not exist" || missing_post["data"] != "post" {
				t.Fatalf("missing-post response = %#v, want not exist/post", missing_post)
			}
			for _, query := range []string{
				`insert into bbs_data (set_name, set_code, set_id, set_data) values ('title', '1', '4', 'Orphan board post')`,
				`insert into bbs_data (set_name, set_code, set_id, set_data) values ('user_id', '1', '4', 'post_owner')`,
			} {
				if _, err := db.Exec(query); err != nil {
					t.Fatalf("seed orphan board post: %v", err)
				}
			}
			missing_board := route.Api_bbs_watch_post(owner_config, "4", "1")
			if missing_board["response"] != "not exist" || missing_board["data"] != "bbs" {
				t.Fatalf("missing-board response = %#v, want not exist/bbs", missing_board)
			}

			private_viewer := tool.Config{IP: "192.0.2.88", Session: &emptySession{}}
			denied := route.Api_bbs_watch_post(private_viewer, "2", "1")
			if denied["response"] != "require auth" {
				t.Fatalf("private-board watch response = %#v, want require auth", denied)
			}
			if count := Bbs_watchlist_row_count(db, t, "192.0.2.88", "2-1"); count != 0 {
				t.Fatalf("denied viewer acquired a watch row: %d", count)
			}

			for _, expected_count := range []int{1, 0} {
				result := route.Api_w_watch_list_post(owner_config, "2-1", "bbs_watchlist")
				if result["response"] != "ok" {
					t.Fatalf("standalone watch-list API toggle = %#v, want ok", result)
				}
				if count := Bbs_watchlist_row_count(db, t, "admin_ip", "2-1"); count != expected_count {
					t.Fatalf("standalone BBS watch row count = %d, want %d", count, expected_count)
				}
			}
		})
	}
}

func Test_Api_bbs_watch_post_concurrent_toggles(t *testing.T) {
	backends := []string{"sqlite", "turso"}
	if os.Getenv("MYSQL_AUDIT_DSN") != "" {
		backends = append(backends, "mysql")
	}

	for _, backend := range backends {
		t.Run(backend, func(t *testing.T) {
			db := new_record_db(t, backend)
			Seed_watchlist_test_board(db, t)
			config := tool.Config{IP: "admin_ip", Session: &emptySession{}}

			const toggle_count = 8
			results := make(chan map[string]any, toggle_count)
			var wait_group sync.WaitGroup
			for range toggle_count {
				wait_group.Add(1)
				go func() {
					defer wait_group.Done()
					results <- route.Api_bbs_watch_post(config, "2", "1")
				}()
			}
			wait_group.Wait()
			close(results)
			for result := range results {
				if result["response"] != "ok" {
					t.Errorf("concurrent BBS watch toggle = %#v, want ok", result)
				}
			}
			if count := Bbs_watchlist_row_count(db, t, "admin_ip", "2-1"); count != 0 {
				t.Fatalf("watch row after %d serialized toggles = %d, want 0", toggle_count, count)
			}
		})
	}
}

func Seed_watchlist_test_board(db *sql.DB, t *testing.T) {
	t.Helper()
	for _, query := range []string{
		`insert into bbs_set (set_name, set_code, set_id, set_data) values ('bbs_name', '', '2', 'Watch board')`,
		`insert into bbs_set (set_name, set_code, set_id, set_data) values ('bbs_view_acl', '', '2', 'owner')`,
		`insert into bbs_data (set_name, set_code, set_id, set_data) values ('title', '1', '2', 'Watched post')`,
		`insert into bbs_data (set_name, set_code, set_id, set_data) values ('user_id', '1', '2', 'post_owner')`,
	} {
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("seed watch board: %v", err)
		}
	}
}

func Bbs_watchlist_row_count(db *sql.DB, t *testing.T, user_name string, key string) int {
	t.Helper()
	var count int
	if err := db.QueryRow(`select count(*) from user_set where id = ? and name = 'bbs_watchlist' and data = ?`, user_name, key).Scan(&count); err != nil {
		t.Fatalf("count BBS watch rows: %v", err)
	}
	return count
}
