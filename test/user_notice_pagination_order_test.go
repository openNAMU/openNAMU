package test

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

type user_notice_row struct {
	id     string
	data   string
	date   string
	readme string
}

func Test_user_notice_pagination_order(t *testing.T) {
	backends := []string{"sqlite", "turso"}
	if os.Getenv("MYSQL_AUDIT_DSN") != "" {
		backends = append(backends, "mysql")
	}

	for _, backend := range backends {
		t.Run(backend, func(t *testing.T) {
			db := new_record_db(t, backend)
			defer tool.DB_close(db)

			t.Run("counts_and_order", func(t *testing.T) {
				user := "order_user"
				for i := 1; i <= 101; i++ {
					exec_user_notice(t, db, user, fmt.Sprint(i), "order_data_"+fmt.Sprint(i), "2026-01-01 12:00:00", "")
				}

				page1 := user_notice_ids(db, user, 0)
				page2 := user_notice_ids(db, user, 50)
				page3 := user_notice_ids(db, user, 100)

				t.Run("page1_order", func(t *testing.T) {
					if len(page1) != 50 {
						t.Fatalf("page1 count = %d, want 50", len(page1))
					}
					expected := ids_desc(101, 52)
					if !slice_equal(page1, expected) {
						t.Errorf("page1 order = %v, want %v", page1, expected)
					}
				})

				t.Run("page2_order", func(t *testing.T) {
					if len(page2) != 50 {
						t.Fatalf("page2 count = %d, want 50", len(page2))
					}
					expected := ids_desc(51, 2)
					if !slice_equal(page2, expected) {
						t.Errorf("page2 order = %v, want %v", page2, expected)
					}
				})

				t.Run("page3_order", func(t *testing.T) {
					if len(page3) != 1 {
						t.Fatalf("page3 count = %d, want 1", len(page3))
					}
					if page3[0] != "1" {
						t.Errorf("page3 id = %s, want 1", page3[0])
					}
				})

				all := append(append(page1, page2...), page3...)
				if len(all) != 101 {
					t.Fatalf("total collected = %d, want 101", len(all))
				}
				if !sorted_slice_equal(all, ids_asc(1, 101)) {
					t.Errorf("all pages missing or duplicated IDs; want sorted 1..101")
				}
			})

			t.Run("same_timestamp", func(t *testing.T) {
				user := "same_ts_user"
				for i := 1; i <= 55; i++ {
					exec_user_notice(t, db, user, fmt.Sprint(i), "same_ts_data_"+fmt.Sprint(i), "2026-01-01 12:00:00", "")
				}

				page1 := user_notice_ids(db, user, 0)
				page2 := user_notice_ids(db, user, 50)

				if len(page1) != 50 || len(page2) != 5 {
					t.Fatalf("same_timestamp page counts = %d, %d, want 50, 5", len(page1), len(page2))
				}
				if !slice_equal(page1, ids_desc(55, 6)) {
					t.Errorf("same_timestamp page1 = %v, want %v", page1, ids_desc(55, 6))
				}
				if !slice_equal(page2, ids_desc(5, 1)) {
					t.Errorf("same_timestamp page2 = %v, want %v", page2, ids_desc(5, 1))
				}

				repeat1 := user_notice_ids(db, user, 0)
				repeat2 := user_notice_ids(db, user, 50)
				if !slice_equal(page1, repeat1) || !slice_equal(page2, repeat2) {
					t.Errorf("same_timestamp ordering is not stable across repeated queries")
				}
			})

			t.Run("counts_boundary", func(t *testing.T) {
				for _, total := range []int{49, 50, 51, 100, 101} {
					t.Run(fmt.Sprintf("total=%d", total), func(t *testing.T) {
						user := fmt.Sprintf("boundary_user_%d", total)
						for i := 1; i <= total; i++ {
							exec_user_notice(t, db, user, fmt.Sprint(i), "boundary_data_"+fmt.Sprint(i), "2026-01-01 12:00:00", "")
						}

						page1 := user_notice_ids(db, user, 0)
						page2 := user_notice_ids(db, user, 50)
						page3 := user_notice_ids(db, user, 100)

						want1 := total
						if want1 > 50 {
							want1 = 50
						}
						want2 := 0
						if total > 50 {
							want2 = total - 50
							if want2 > 50 {
								want2 = 50
							}
						}
						want3 := 0
						if total > 100 {
							want3 = total - 100
						}

						if len(page1) != want1 || len(page2) != want2 || len(page3) != want3 {
							t.Fatalf("boundary page counts = %d, %d, %d, want %d, %d, %d", len(page1), len(page2), len(page3), want1, want2, want3)
						}
					})
				}
			})

			t.Run("read_unread_states", func(t *testing.T) {
				user := "readme_user"
				exec_user_notice(t, db, user, "1", "readme_unread", "2026-01-01 12:00:00", "")
				exec_user_notice(t, db, user, "2", "readme_read", "2026-01-01 12:00:00", "1")
				exec_user_notice(t, db, user, "3", "readme_unread_second", "2026-01-01 12:00:00", "")

				rows := user_notice_rows(db, user, 0)
				if len(rows) != 3 {
					t.Fatalf("readme rows = %d, want 3", len(rows))
				}
				if rows[0].id != "3" || rows[1].id != "2" || rows[2].id != "1" {
					t.Fatalf("readme order = %v, want 3,2,1", ids_of_rows(rows))
				}
				if rows[1].readme != "1" {
					t.Errorf("id 2 readme = %q, want 1", rows[1].readme)
				}

				html := view_page( tool.Config{IP: user, Session: &emptySession{}}, user, 1)
				if !strings.Contains(html, "style=\"opacity: 0.75;\"") {
					t.Error("view missing read notification opacity style")
				}
			})

			t.Run("separate_inboxes", func(t *testing.T) {
				exec_user_notice(t, db, "user_a", "1", "a_data_1", "2026-01-01 12:00:00", "")
				exec_user_notice(t, db, "user_a", "2", "a_data_2", "2026-01-01 12:00:00", "")
				exec_user_notice(t, db, "user_b", "1", "b_data_1", "2026-01-01 12:00:00", "")
				exec_user_notice(t, db, "user_b", "2", "b_data_2", "2026-01-01 12:00:00", "")
				exec_user_notice(t, db, "user_b", "3", "b_data_3", "2026-01-01 12:00:00", "")

				a := user_notice_ids(db, "user_a", 0)
				b := user_notice_ids(db, "user_b", 0)

				if len(a) != 2 || len(b) != 3 {
					t.Fatalf("inbox counts = %d, %d, want 2, 3", len(a), len(b))
				}
				if !slice_equal(a, ids_desc(2, 1)) {
					t.Errorf("user_a IDs = %v, want %v", a, ids_desc(2, 1))
				}
				if !slice_equal(b, ids_desc(3, 1)) {
					t.Errorf("user_b IDs = %v, want %v", b, ids_desc(3, 1))
				}
			})

			t.Run("owner_viewing_another_user", func(t *testing.T) {
				exec_user_notice(t, db, "bob1", "1", "bob_data_1", "2026-01-01 12:00:00", "")
				exec_user_notice(t, db, "bob1", "2", "bob_data_2", "2026-01-01 12:00:00", "")
				exec_user_notice(t, db, "alice", "1", "alice_data_1", "2026-01-01 12:00:00", "")

				html := view_page( tool.Config{IP: "admin_ip", Session: &emptySession{}}, "bob1", 1)
				if strings.Contains(html, "alice_data_1") {
					t.Errorf("owner viewing another user leaked alice's notification into bob's inbox")
				}
				if !strings.Contains(html, "bob_data_2") {
					t.Errorf("owner viewing another user missing bob_data_2")
				}
				if !strings.Contains(html, `/alarm_user/bob1/delete/2`) {
					t.Errorf("owner viewing another user missing delete link for bob1 id 2")
				}
			})

			t.Run("delete_single_and_pagination", func(t *testing.T) {
				user := "delete_user"
				for i := 1; i <= 101; i++ {
					exec_user_notice(t, db, user, fmt.Sprint(i), "delete_data_"+fmt.Sprint(i), "2026-01-01 12:00:00", "")
				}

				delete_id := "75"
				out := route.Api_alarm_delete_post(tool.Config{IP: "admin_ip", Session: &emptySession{}}, user, delete_id, false)
				if out["response"] != "ok" {
					t.Fatalf("delete response = %v, want ok", out)
				}

				page1 := user_notice_ids(db, user, 0)
				page2 := user_notice_ids(db, user, 50)
				if len(page1) != 50 || len(page2) != 50 {
					t.Fatalf("delete page counts = %d, %d, want 50, 50", len(page1), len(page2))
				}

				all := append(append(page1, page2...), page3_ids(db, user, 100)...)
				if len(all) != 100 {
					t.Fatalf("post-delete collected = %d, want 100", len(all))
				}
				for _, id := range all {
					if id == delete_id {
						t.Fatalf("deleted notification id %s still present", id)
					}
				}
				if !sorted_slice_equal(all, ids_except(1, 101, delete_id)) {
					t.Errorf("post-delete pagination missing or duplicated IDs")
				}

				html1 := view_page( tool.Config{IP: user, Session: &emptySession{}}, user, 1)
				html2 := view_page( tool.Config{IP: user, Session: &emptySession{}}, user, 2)
				if !has_next(html1, db) {
					t.Errorf("page 1 after deletion should have next link because 100 records remain")
				}
				if !has_prev(html2, db) {
					t.Errorf("page 2 after deletion missing previous link")
				}
				if has_next(html2, db) {
					t.Errorf("page 2 after deletion still has next link, want none")
				}
			})

			t.Run("view_links_read_status", func(t *testing.T) {
				user := "link_user"
				exec_user_notice(t, db, user, "1", "one", "2026-01-01 12:00:00", "")
				exec_user_notice(t, db, user, "2", "two", "2026-01-01 12:00:00", "1")
				exec_user_notice(t, db, user, "3", "three", "2026-01-01 12:00:00", "")

				html := view_page( tool.Config{IP: user, Session: &emptySession{}}, user, 1)
				for _, id := range []string{"3", "2", "1"} {
					if !strings.Contains(html, `/alarm/delete/`+id) {
						t.Errorf("missing delete link for notification %s", id)
					}
				}
				if !strings.Contains(html, "style=\"opacity: 0.75;\"") {
					t.Error("view missing read notification opacity style")
				}
				if !strings.Contains(html, "three") && !strings.Contains(html, "two") && !strings.Contains(html, "one") {
					t.Error("view missing notification data")
				}
			})
		})
	}
}

func has_prev(html string, db *sql.DB) bool {
	label := tool.Get_language(db, "previous", true)
	return strings.Contains(html, "("+label+")")
}

func has_next(html string, db *sql.DB) bool {
	label := tool.Get_language(db, "next", true)
	return strings.Contains(html, "("+label+")")
}

func view_page(config tool.Config, user_name string, page int) string {
	values := to_url_values(map[string]string{"num": fmt.Sprint(page)})
	return route.View_alarm(config, user_name, values)
}

func to_url_values(values map[string]string) url.Values {
	v := url.Values{}
	for key, value := range values {
		v.Set(key, value)
	}
	return v
}

func exec_user_notice(t *testing.T, db *sql.DB, name string, id, data, date, readme string) {
	if _, err := db.Exec(
		"insert into user_notice (id, name, data, date, readme) values (?, ?, ?, ?, ?)",
		id, name, data, date, readme,
	); err != nil {
		t.Fatalf("exec_user_notice %q: %v", name, err)
	}
}

func user_notice_ids(db *sql.DB, name string, offset int) []string {
	rows := tool.Get_user_notice_rows(db, name, offset)
	defer rows.Close()

	ids := []string{}
	for rows.Next() {
		var id, data, date, readme string
		if rows.Scan(&id, &data, &date, &readme) != nil {
			continue
		}
		ids = append(ids, id)
	}
	return ids
}

func user_notice_rows(db *sql.DB, name string, offset int) []user_notice_row {
	rows := tool.Get_user_notice_rows(db, name, offset)
	defer rows.Close()

	result := []user_notice_row{}
	for rows.Next() {
		var id, data, date, readme string
		if rows.Scan(&id, &data, &date, &readme) != nil {
			continue
		}
		result = append(result, user_notice_row{id: id, data: data, date: date, readme: readme})
	}
	return result
}

func page3_ids(db *sql.DB, name string, offset int) []string {
	return user_notice_ids(db, name, offset)
}

func ids_of_rows(rows []user_notice_row) []string {
	ids := []string{}
	for _, row := range rows {
		ids = append(ids, row.id)
	}
	return ids
}

func ids_desc(start, stop int) []string {
	out := []string{}
	for i := start; i >= stop; i-- {
		out = append(out, fmt.Sprint(i))
	}
	return out
}

func ids_asc(start, stop int) []string {
	out := []string{}
	for i := start; i <= stop; i++ {
		out = append(out, fmt.Sprint(i))
	}
	return out
}

func ids_except(start, stop int, remove string) []string {
	out := []string{}
	for i := start; i <= stop; i++ {
		id := fmt.Sprint(i)
		if id == remove {
			continue
		}
		out = append(out, id)
	}
	return out
}

func slice_equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func sorted_slice_equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	sa := append([]string{}, a...)
	sb := append([]string{}, b...)
	sort.Strings(sa)
	sort.Strings(sb)
	return slice_equal(sa, sb)
}
