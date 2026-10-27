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

func recent_change_exec(db *sql.DB, t *testing.T, q string, args ...any) {
	if _, err := db.Exec(q, args...); err != nil {
		t.Fatalf("recent_change seed %q: %v", q, err)
	}
}

func recent_change_seed_watch(db *sql.DB, t *testing.T, user string, titles []string) {
	for _, title := range titles {
		recent_change_exec(db, t,
			`insert into user_set (name, id, data) values ('watchlist', ?, ?)`,
			user, title,
		)
	}
}

func recent_change_seed_rc(db *sql.DB, t *testing.T, titles []string) {
	for i, title := range titles {
		id := fmt.Sprintf("rc_%d", i+1)
		date := fmt.Sprintf("2026-01-02 %02d:00:00", i+1)
		recent_change_exec(db, t,
			`insert into rc (id, title, date, type) values (?, ?, ?, 'normal')`,
			id, title, date,
		)
	}
}

func recent_change_seed_history(db *sql.DB, t *testing.T, titles []string, hidden int) {
	for i, title := range titles {
		id := fmt.Sprintf("his_%d", i+1)
		date := fmt.Sprintf("2026-01-01 %02d:00:00", i+1)
		hide := ""
		if hidden > 0 && i == hidden-1 {
			hide = "hide"
		}
		recent_change_exec(db, t,
			`insert into history (id, title, data, date, ip, send, leng, hide, type) values (?, ?, ?, ?, ?, ?, ?, ?, 'normal')`,
			id, title, "", date, "1.1.1.1", "", "0", hide,
		)
	}
}

func recent_change_backends() []string {
	backends := []string{"sqlite", "turso"}
	if os.Getenv("MYSQL_AUDIT_DSN") != "" {
		backends = append(backends, "mysql")
	}
	return backends
}

func Test_recent_change_pagination(t *testing.T) {
	for _, backend := range recent_change_backends() {
		t.Run(backend, func(t *testing.T) {
			t.Run("50_items_no_next", func(t *testing.T) {
				db := new_record_db(t, backend)
				defer tool.DB_close(db)

				audit_seed_auth(db, t)
				recent_change_seed_history(db, t, make_titles(50), 0)

				config_admin := tool.Config{IP: "admin_ip", Session: &emptySession{}}
				out := route.View_list_recent_change(config_admin, "normal", "50", "1")

				if strings.Contains(out, `href="/recent_change/2/normal"`) {
					t.Fatalf("50 items page1 next link present")
				}
			})

			t.Run("49_items_no_next", func(t *testing.T) {
				db := new_record_db(t, backend)
				defer tool.DB_close(db)

				audit_seed_auth(db, t)
				recent_change_seed_history(db, t, make_titles(49), 0)

				config_admin := tool.Config{IP: "admin_ip", Session: &emptySession{}}
				out := route.View_list_recent_change(config_admin, "normal", "50", "1")

				if strings.Contains(out, `href="/recent_change/2/normal"`) {
					t.Fatalf("49 items page1 next link present")
				}
			})

			t.Run("51_items_next", func(t *testing.T) {
				db := new_record_db(t, backend)
				defer tool.DB_close(db)

				audit_seed_auth(db, t)
				recent_change_seed_history(db, t, make_titles(51), 0)

				config_admin := tool.Config{IP: "admin_ip", Session: &emptySession{}}
				page1 := route.View_list_recent_change(config_admin, "normal", "50", "1")
				page2 := route.View_list_recent_change(config_admin, "normal", "50", "2")

				if !strings.Contains(page1, `href="/recent_change/2/normal"`) {
					t.Fatalf("51 items page1 next link missing")
				}
				if !strings.Contains(page2, `href="/recent_change/1/normal"`) {
					t.Fatalf("51 items page2 previous link missing")
				}
				if strings.Contains(page2, `href="/recent_change/3/normal"`) {
					t.Fatalf("51 items page2 next link present")
				}
			})

			t.Run("rc_history_mixed_boundary", func(t *testing.T) {
				db := new_record_db(t, backend)
				defer tool.DB_close(db)

				audit_seed_auth(db, t)
				rc_titles := make_titles(30)
				his_titles := make_titles(30)
				for i := range his_titles {
					his_titles[i] = "his_" + strings.TrimPrefix(his_titles[i], "doc_")
				}
				recent_change_seed_rc(db, t, rc_titles)
				recent_change_seed_history(db, t, his_titles, 0)

				config_admin := tool.Config{IP: "admin_ip", Session: &emptySession{}}
				page1 := route.Api_list_recent_change(config_admin, "normal", "50", "1")
				page2 := route.Api_list_recent_change(config_admin, "normal", "50", "2")

				page1_data := page1["data"].([][]string)
				page2_data := page2["data"].([][]string)
				if len(page1_data) != 50 {
					t.Fatalf("mixed page1 len = %d, want 50", len(page1_data))
				}
				if len(page2_data) != 10 {
					t.Fatalf("mixed page2 len = %d, want 10", len(page2_data))
				}

				seen := make(map[string]bool)
				for _, row := range append(page1_data, page2_data...) {
					key := row[0] + "|" + row[1]
					if seen[key] {
						t.Fatalf("duplicate rc/history row: %q", key)
					}
					seen[key] = true
				}
				if len(seen) != 60 {
					t.Fatalf("mixed distinct rows = %d, want 60", len(seen))
				}
			})

			t.Run("hidden_records_maintained", func(t *testing.T) {
				db := new_record_db(t, backend)
				defer tool.DB_close(db)

				audit_seed_auth(db, t)
				recent_change_seed_history(db, t, make_titles(50), 25)

				config_admin := tool.Config{IP: "admin_ip", Session: &emptySession{}}
				config_ip := tool.Config{IP: "1.1.1.1", Session: &emptySession{}}

				admin_api := route.Api_list_recent_change(config_admin, "normal", "50", "1")
				ip_api := route.Api_list_recent_change(config_ip, "normal", "50", "1")

				admin_data := admin_api["data"].([][]string)
				ip_data := ip_api["data"].([][]string)

				if len(admin_data) != 50 {
					t.Fatalf("admin hidden len = %d, want 50", len(admin_data))
				}
				if len(ip_data) != 50 {
					t.Fatalf("ip hidden len = %d, want 50", len(ip_data))
				}

				hidden_count := 0
				for _, row := range ip_data {
					if row[1] == "" && row[6] != "" {
						hidden_count++
					}
				}
				if hidden_count != 1 {
					t.Fatalf("ip hidden placeholder count = %d, want 1", hidden_count)
				}

				admin_out := route.View_list_recent_change(config_admin, "normal", "50", "1")
				if strings.Contains(admin_out, `href="/recent_change/2/normal"`) {
					t.Fatalf("admin hidden 50 items page1 next link present")
				}

				ip_out := route.View_list_recent_change(config_ip, "normal", "50", "1")
				if strings.Contains(ip_out, `href="/recent_change/2/normal"`) {
					t.Fatalf("ip hidden 50 items page1 next link present")
				}
			})

			t.Run("no_history_view_excludes_history", func(t *testing.T) {
				db := new_record_db(t, backend)
				defer tool.DB_close(db)

				audit_seed_auth(db, t)
				recent_change_exec(db, t,
					`insert into alist (name, acl) values ('nothing_user_group','nothing')`)
				recent_change_exec(db, t,
					`insert into user_set (name, id, data) values ('acl','nothing_user','nothing_user_group')`)

				recent_change_seed_history(db, t, make_titles(51), 0)

				config_admin := tool.Config{IP: "admin_ip", Session: &emptySession{}}
				config_ip := tool.Config{IP: "1.1.1.1", Session: &emptySession{}}
				config_none := tool.Config{IP: "nothing_user", Session: &emptySession{}}

				admin_api := route.Api_list_recent_change(config_admin, "normal", "50", "1")
				ip_api := route.Api_list_recent_change(config_ip, "normal", "50", "1")
				none_api := route.Api_list_recent_change(config_none, "normal", "50", "1")

				if len(admin_api["data"].([][]string)) != 50 {
					t.Fatalf("admin len = %d, want 50", len(admin_api["data"].([][]string)))
				}
				if len(ip_api["data"].([][]string)) != 50 {
					t.Fatalf("ip len = %d, want 50", len(ip_api["data"].([][]string)))
				}
				if len(none_api["data"].([][]string)) != 0 {
					t.Fatalf("none history_view len = %d, want 0", len(none_api["data"].([][]string)))
				}
			})

			t.Run("watch_list_boundary", func(t *testing.T) {
				db := new_record_db(t, backend)
				defer tool.DB_close(db)

				audit_seed_auth(db, t)

				watched := make_titles(50)
				unwatched := make_titles(50)
				for i := range unwatched {
					title := unwatched[i]
					unwatched[i] = "unwatched_" + strings.TrimPrefix(title, "doc_")
				}

				recent_change_seed_watch(db, t, "watch_user", watched)
				recent_change_seed_rc(db, t, watched)
				recent_change_seed_rc(db, t, unwatched)

				config_watch := tool.Config{IP: "watch_user", Session: &emptySession{}}

				watch_page1 := route.Api_list_recent_change(config_watch, "watch", "50", "1")
				normal_page1 := route.Api_list_recent_change(config_watch, "normal", "50", "1")
				normal_page2 := route.Api_list_recent_change(config_watch, "normal", "50", "2")

				watch_data := watch_page1["data"].([][]string)
				normal1 := normal_page1["data"].([][]string)
				normal2 := normal_page2["data"].([][]string)

				if len(watch_data) != 50 {
					t.Fatalf("watch len = %d, want 50", len(watch_data))
				}
				if len(normal1) != 50 {
					t.Fatalf("normal page1 len = %d, want 50", len(normal1))
				}
				if len(normal2) == 0 {
					t.Fatalf("normal page2 len = %d, want >0", len(normal2))
				}

				watch_out := route.View_list_recent_change(config_watch, "watch", "50", "1")
				if strings.Contains(watch_out, `href="/recent_change/2/watch"`) {
					t.Fatalf("watch 50 items page1 next link present")
				}

				normal_out1 := route.View_list_recent_change(config_watch, "normal", "50", "1")
				if !strings.Contains(normal_out1, `href="/recent_change/2/normal"`) {
					t.Fatalf("normal 100 items page1 next link missing")
				}
			})
		})
	}
}

func make_titles(n int) []string {
	titles := []string{}
	for i := 1; i <= n; i++ {
		titles = append(titles, fmt.Sprintf("doc_%d", i))
	}
	return titles
}
