package test

import (
	"database/sql"
	"os"
	"strings"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

func Test_document_watchlist_notification_direction(t *testing.T) {
	backends := []string{"sqlite", "turso"}
	if os.Getenv("MYSQL_AUDIT_DSN") != "" {
		backends = append(backends, "mysql")
	}

	for _, backend := range backends {
		t.Run(backend, func(t *testing.T) {
			db := new_record_db(t, backend)
			defer tool.DB_close(db)

			exec := func(q string, args ...any) {
				if _, err := db.Exec(q, args...); err != nil {
					t.Fatalf("seed %q: %v", q, err)
				}
			}

			config := tool.Config{IP: "admin_ip", Session: &emptySession{}}

			t.Run("helper_multiple_watchers", func(t *testing.T) {
				doc := "WatchDoc1"
				bob := "Bob1"
				charlie := "Charlie1"
				unrelated := "Dave1"

				exec("insert into data (title, data) values (?, ?)", doc, "initial content")
				exec(`insert into user_set (name, id, data) values ('watchlist', ?, ?)`, bob, doc)
				exec(`insert into user_set (name, id, data) values ('watchlist', ?, ?)`, charlie, doc)
				exec(`insert into user_set (name, id, data) values ('watchlist', ?, ?)`, unrelated, "OtherDoc1")

				before := map[string]int{
					"admin_ip": 0,
					bob:        0,
					charlie:    0,
					unrelated:  0,
				}
				for name := range before {
					if count := notice_count(db, name); count != 0 {
						t.Fatalf("helper_multiple_watchers: %s already has %d notices", name, count)
					}

				}

				tool.Do_watchlist_alarm_send(db, config, doc)

				after := map[string]int{}
				for name := range before {
					after[name] = notice_count(db, name)
				}

				t.Logf(
					"backend=%s helper_multiple_watchers before=%v after=%v",
					backend,
					before,
					after,
				)

				if after["admin_ip"] != 0 {
					t.Errorf("helper_multiple_watchers: admin_ip received %d notices, want 0", after["admin_ip"])
				}
				if after[bob] != 1 {
					t.Errorf("helper_multiple_watchers: %s received %d notices, want 1", bob, after[bob])
				}
				if after[charlie] != 1 {
					t.Errorf("helper_multiple_watchers: %s received %d notices, want 1", charlie, after[charlie])
				}
				if after[unrelated] != 0 {
					t.Errorf(
						"helper_multiple_watchers: unrelated %s received %d notices, want 0",
						unrelated,
						after[unrelated],
					)
				}

				bob_data := latest_notice_data(db, bob)
				if !strings.Contains(bob_data, "admin_ip |") {
					t.Errorf("helper_multiple_watchers: Bob notification missing sender, got %q", bob_data)
				}
				if !strings.Contains(bob_data, doc) {
					t.Errorf("helper_multiple_watchers: Bob notification missing document %q", doc)
				}
			})

			t.Run("integration_multiple_watchers", func(t *testing.T) {
				doc := "WatchDoc2"
				bob := "Bob2"
				charlie := "Charlie2"
				unrelated := "Dave2"

				exec("insert into data (title, data) values (?, ?)", doc, "initial content")
				exec(`insert into user_set (name, id, data) values ('watchlist', ?, ?)`, bob, doc)
				exec(`insert into user_set (name, id, data) values ('watchlist', ?, ?)`, charlie, doc)
				exec(`insert into user_set (name, id, data) values ('watchlist', ?, ?)`, unrelated, "OtherDoc2")

				for _, name := range []string{"admin_ip", bob, charlie, unrelated} {
					if count := notice_count(db, name); count != 0 {
						t.Fatalf("integration_multiple_watchers: %s already has %d notices", name, count)
					}
				}

				out := route.Api_edit_post(config, doc, "edited content", "", "", "")

				if resp, ok := out["response"].(string); !ok || resp != "ok" {
					t.Fatalf("integration_multiple_watchers: response = %q, want ok", out)
				}

				after := map[string]int{}
				for _, name := range []string{"admin_ip", bob, charlie, unrelated} {
					after[name] = notice_count(db, name)
				}

				t.Logf(
					"backend=%s integration_multiple_watchers editor=%d bob=%d charlie=%d unrelated=%d response=%q",
					backend,
					after["admin_ip"],
					after[bob],
					after[charlie],
					after[unrelated],
					out["response"],
				)

				if after["admin_ip"] != 0 {
					t.Errorf(
						"integration_multiple_watchers: editor received %d notices, want 0",
						after["admin_ip"],
					)
				}
				if after[bob] != 1 {
					t.Errorf("integration_multiple_watchers: %s received %d notices, want 1", bob, after[bob])
				}
				if after[charlie] != 1 {
					t.Errorf("integration_multiple_watchers: %s received %d notices, want 1", charlie, after[charlie])
				}
				if after[unrelated] != 0 {
					t.Errorf(
						"integration_multiple_watchers: unrelated %s received %d notices, want 0",
						unrelated,
						after[unrelated],
					)
				}
			})

			t.Run("self_watch_no_self_notification", func(t *testing.T) {
				doc := "SelfDoc1"
				exec("insert into data (title, data) values (?, ?)", doc, "initial content")
				exec(`insert into user_set (name, id, data) values ('watchlist', ?, ?)`, "admin_ip", doc)

				if count := notice_count(db, "admin_ip"); count != 0 {
					t.Fatalf("self_watch_no_self_notification: admin_ip already has %d notices", count)
				}

				tool.Do_watchlist_alarm_send(db, config, doc)

				if count := notice_count(db, "admin_ip"); count != 0 {
					t.Errorf("self_watch_no_self_notification: admin_ip received %d self-notices, want 0", count)
				}
				t.Logf(
					"backend=%s self_watch_no_self_notification editor=%d",
					backend,
					notice_count(db, "admin_ip"),
				)
			})

			t.Run("other_document_no_notification", func(t *testing.T) {
				doc := "MainDoc1"
				other_doc := "OtherDoc3"
				other_watcher := "OtherWatcher"

				exec("insert into data (title, data) values (?, ?)", doc, "initial main content")
				exec(`insert into user_set (name, id, data) values ('watchlist', ?, ?)`, other_watcher, other_doc)

				for _, name := range []string{"admin_ip", other_watcher} {
					if count := notice_count(db, name); count != 0 {
						t.Fatalf("other_document_no_notification: %s already has %d notices", name, count)
					}
				}

				out := route.Api_edit_post(config, doc, "edited main content", "", "", "")

				if resp, ok := out["response"].(string); !ok || resp != "ok" {
					t.Fatalf("other_document_no_notification: response = %q, want ok", out)
				}

				t.Logf(
					"backend=%s other_document_no_notification editor=%d other_watcher=%d response=%q",
					backend,
					notice_count(db, "admin_ip"),
					notice_count(db, other_watcher),
					out["response"],
				)

				if count := notice_count(db, "admin_ip"); count != 0 {
					t.Errorf("other_document_no_notification: editor received %d notices, want 0", count)
				}
				if count := notice_count(db, other_watcher); count != 0 {
					t.Errorf(
						"other_document_no_notification: %s received %d notices for unrelated document, want 0",
						other_watcher,
						count,
					)
				}
			})
		})
	}
}

func latest_notice_data(db *sql.DB, name string) string {
	var data string
	exists := tool.QueryRow_DB(
		db,
		"select data from user_notice where name = ? order by id + 0 desc limit 1",
		[]any{&data},
		name,
	)
	if !exists {
		return ""
	}
	return data
}
