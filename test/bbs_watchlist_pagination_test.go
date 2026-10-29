package test

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

func Test_bbs_watchlist_pagination(t *testing.T) {
	backends := []string{"sqlite", "turso"}
	if os.Getenv("MYSQL_AUDIT_DSN") != "" {
		backends = append(backends, "mysql")
	}

	for _, backend := range backends {
		t.Run(backend, func(t *testing.T) {
			t.Run("totals_page_links", func(t *testing.T) {
				for _, total := range []int{49, 50, 51, 100, 101} {
					t.Run(fmt.Sprintf("total=%d", total), func(t *testing.T) {
						db := new_record_db(t, backend)
						defer tool.DB_close(db)

						user := "pag_user"
						seed_user_auth(db, t, user)
						seed_bbs_watch_codes(db, t, user, "2", bbs_watch_codes_zero(total))

						page1 := bbs_watch_page_items(db, t, user, 1)
						page2 := bbs_watch_page_items(db, t, user, 2)
						page3 := bbs_watch_page_items(db, t, user, 3)

						want1 := bbs_watch_min(total, 50)
						want2 := bbs_watch_min(bbs_watch_sub(total, 50), 50)
						want3 := bbs_watch_max(total - 100, 0)

						if len(page1) != want1 {
							t.Fatalf("total=%d page1 len = %d, want %d", total, len(page1), want1)
						}
						if len(page2) != want2 {
							t.Fatalf("total=%d page2 len = %d, want %d", total, len(page2), want2)
						}
						if len(page3) != want3 {
							t.Fatalf("total=%d page3 len = %d, want %d", total, len(page3), want3)
						}

						all := collect_bbs_watch_pages(db, t, user)
						want_all := bbs_watch_codes_with_set("2", bbs_watch_codes_zero(total))
						if !slice_equal(all, want_all) {
							t.Fatalf("total=%d all pages = %v, want %v", total, all, want_all)
						}

						html1 := view_bbs_watch_page(db, t, user, 1)
						assert_page_links(t, html1, 1, false, total > 50)

						html2 := view_bbs_watch_page(db, t, user, 2)
						assert_page_links(t, html2, 2, true, total > 100)

						html3 := view_bbs_watch_page(db, t, user, 3)
						assert_page_links(t, html3, 3, true, false)
					})
				}
			})

			t.Run("exactly_50_no_false_next", func(t *testing.T) {
				db := new_record_db(t, backend)
				defer tool.DB_close(db)

				user := "exactly_50_user"
				seed_user_auth(db, t, user)
				seed_bbs_watch_codes(db, t, user, "2", bbs_watch_codes_zero(50))

				html := view_bbs_watch_page(db, t, user, 1)
				assert_page_links(t, html, 1, false, false)

				if strings.Contains(html, `href="/bbs_watch_list_page/2">`) {
					t.Fatalf("exactly 50 entries still show a next link to page 2")
				}
			})

			t.Run("stable_order_and_insertion_order", func(t *testing.T) {
				codes := bbs_watch_codes_zero(12)

				t.Run("repeated_calls", func(t *testing.T) {
					db := new_record_db(t, backend)
					defer tool.DB_close(db)

					user := "stable_user"
					seed_user_auth(db, t, user)
					seed_bbs_watch_codes(db, t, user, "2", codes)

					first := bbs_watch_page_items(db, t, user, 1)
					second := bbs_watch_page_items(db, t, user, 1)
					if !slice_equal(first, second) {
						t.Fatalf("repeated calls differ: first=%v second=%v", first, second)
					}
				})

				t.Run("varying_insertion_order", func(t *testing.T) {
					user := "insertion_user"
				ascending := append([]int{}, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12)
				descending := make([]int, len(ascending))
				for i := range ascending {
					descending[i] = ascending[len(ascending)-1-i]
				}

					first_db := new_record_db(t, backend)
					defer tool.DB_close(first_db)
					seed_user_auth(first_db, t, user)
					seed_bbs_watch_codes_order(first_db, t, user, "2", ascending)
					ascending_ids := bbs_watch_page_items(first_db, t, user, 1)

					second_db := new_record_db(t, backend)
					defer tool.DB_close(second_db)
					seed_user_auth(second_db, t, user)
					seed_bbs_watch_codes_order(second_db, t, user, "2", descending)
					descending_ids := bbs_watch_page_items(second_db, t, user, 1)

					if !slice_equal(ascending_ids, descending_ids) {
						t.Fatalf("insertion order changed API order: ascending=%v descending=%v", ascending_ids, descending_ids)
					}
				if !slice_equal(ascending_ids, bbs_watch_codes_with_set("2", codes)) {
					t.Fatalf("stable expected order not produced: got %v want %v", ascending_ids, bbs_watch_codes_with_set("2", codes))
				}
				})
			})

			t.Run("legacy_thread_only", func(t *testing.T) {
				db := new_record_db(t, backend)
				defer tool.DB_close(db)

				user := "legacy_only_user"
				seed_user_auth(db, t, user)
				thread_seed_bbs(db, t, "42", "legacy thread 42")
				seed_user_watch(db, t, "thread_watchlist", user, "42")

				items := bbs_watch_api_items(db, t, user, 1)
				if len(items) != 1 {
					t.Fatalf("legacy only items = %d, want 1", len(items))
				}
				if items[0].Set_id != "-1" || items[0].Set_code != "42" {
					t.Fatalf("legacy item set = %v/%v, want -1/42", items[0].Set_id, items[0].Set_code)
				}

				html := view_bbs_watch_page(db, t, user, 1)
				if !strings.Contains(html, `href="/bbs/w/-1/42">`) {
					t.Fatalf("legacy thread 42 view link missing")
				}
				if !strings.Contains(html, `href="/bbs_watch/-1/42">`) {
					t.Fatalf("legacy thread 42 unwatch link missing")
				}
			})

			t.Run("modern_thread_only", func(t *testing.T) {
				db := new_record_db(t, backend)
				defer tool.DB_close(db)

				user := "modern_only_user"
				seed_user_auth(db, t, user)
				thread_seed_bbs(db, t, "42", "modern thread 42")
				seed_user_watch(db, t, "bbs_watchlist", user, "-1-42")

				items := bbs_watch_api_items(db, t, user, 1)
				if len(items) != 1 {
					t.Fatalf("modern thread only items = %d, want 1", len(items))
				}
				if items[0].Set_id != "-1" || items[0].Set_code != "42" {
					t.Fatalf("modern thread item set = %v/%v, want -1/42", items[0].Set_id, items[0].Set_code)
				}

				html := view_bbs_watch_page(db, t, user, 1)
				if !strings.Contains(html, `href="/bbs/w/-1/42">`) {
					t.Fatalf("modern thread 42 view link missing")
				}
				if !strings.Contains(html, `href="/bbs_watch/-1/42">`) {
					t.Fatalf("modern thread 42 unwatch link missing")
				}
			})

			t.Run("duplicate_legacy_and_bbs_shown_once", func(t *testing.T) {
				db := new_record_db(t, backend)
				defer tool.DB_close(db)

				user := "dup_user"
				seed_user_auth(db, t, user)
				thread_seed_bbs(db, t, "42", "duplicate thread 42")
				seed_user_watch(db, t, "thread_watchlist", user, "42")
				seed_user_watch(db, t, "bbs_watchlist", user, "-1-42")

				items := bbs_watch_api_items(db, t, user, 1)
				if len(items) != 1 {
					t.Fatalf("duplicate items = %d, want 1", len(items))
				}
				if items[0].Set_id != "-1" || items[0].Set_code != "42" {
					t.Fatalf("duplicate item set = %v/%v, want -1/42", items[0].Set_id, items[0].Set_code)
				}

				html := view_bbs_watch_page(db, t, user, 1)
				if strings.Count(html, `href="/bbs/w/-1/42">`) != 1 {
					t.Errorf("duplicate thread 42 view link count = %d, want 1", strings.Count(html, `href="/bbs/w/-1/42">`))
				}
				if strings.Count(html, `href="/bbs_watch/-1/42">`) != 1 {
					t.Errorf("duplicate thread 42 unwatch link count = %d, want 1", strings.Count(html, `href="/bbs_watch/-1/42">`))
				}
			})

			t.Run("dedupe_51_with_boundary_duplicates", func(t *testing.T) {
				db := new_record_db(t, backend)
				defer tool.DB_close(db)

				user := "dedupe_user"
				seed_user_auth(db, t, user)

				for _, code := range bbs_watch_codes_zero(51) {
					seed_user_watch(db, t, "bbs_watchlist", user, "-1-"+code)
					seed_user_watch(db, t, "thread_watchlist", user, code)
				}

				page1 := bbs_watch_page_items(db, t, user, 1)
				page2 := bbs_watch_page_items(db, t, user, 2)

				if len(page1) != 50 {
					t.Fatalf("dedupe 51 page1 len = %d, want 50", len(page1))
				}
				if len(page2) != 1 {
					t.Fatalf("dedupe 51 page2 len = %d, want 1", len(page2))
				}

				all := collect_bbs_watch_pages(db, t, user)
				if !slice_equal(all, bbs_watch_codes_zero(51)) {
					t.Fatalf("dedupe 51 all pages = %v, want %v", all, bbs_watch_codes_zero(51))
				}

				html1 := view_bbs_watch_page(db, t, user, 1)
				assert_page_links(t, html1, 1, false, true)
				html2 := view_bbs_watch_page(db, t, user, 2)
				assert_page_links(t, html2, 2, true, false)
			})

			t.Run("normal_bbs_not_confused_with_thread", func(t *testing.T) {
				db := new_record_db(t, backend)
				defer tool.DB_close(db)

				user := "normal_user"
				seed_user_auth(db, t, user)

				seed_bbs_post(db, t, "-1", "42", "thread title 42")
				seed_bbs_post(db, t, "2", "42", "normal bbs title 2-42")

				seed_user_watch(db, t, "thread_watchlist", user, "42")
				seed_user_watch(db, t, "bbs_watchlist", user, "2-42")

				items := bbs_watch_api_items(db, t, user, 1)
				if len(items) != 2 {
					t.Fatalf("normal/thread mixed items = %d, want 2", len(items))
				}

				seen := map[string]bool{}
				for _, item := range items {
					if item.Set_id == "-1" && item.Set_code == "42" {
						if seen["thread_42"] {
							t.Fatalf("thread 42 appears more than once")
						}
						seen["thread_42"] = true
					}
					if item.Set_id == "2" && item.Set_code == "42" {
						if seen["normal_2_42"] {
							t.Fatalf("normal BBS 2-42 appears more than once")
						}
						seen["normal_2_42"] = true
					}
					if item.Set_id == "2" && item.Set_code == "42" && item.Data != "2-42" {
						t.Fatalf("normal BBS 2-42 Data = %q, want 2-42", item.Data)
					}
					if item.Set_id == "-1" && item.Set_code == "42" && item.Data != "42" {
						t.Fatalf("thread 42 Data = %q, want 42", item.Data)
					}
				}
				if !seen["thread_42"] || !seen["normal_2_42"] {
					t.Fatalf("expected both thread 42 and normal BBS 2-42 to be visible")
				}

				html := view_bbs_watch_page(db, t, user, 1)
				if !strings.Contains(html, `href="/bbs/w/-1/42">`) {
					t.Fatalf("thread 42 view link missing")
				}
				if !strings.Contains(html, `href="/bbs_watch/-1/42">`) {
					t.Fatalf("thread 42 unwatch link missing")
				}
				if !strings.Contains(html, `href="/bbs/w/2/42">`) {
					t.Fatalf("normal BBS 2-42 view link missing")
				}
				if !strings.Contains(html, `href="/bbs_watch/2/42">`) {
					t.Fatalf("normal BBS 2-42 unwatch link missing")
				}
			})

			t.Run("unwatch_removes_thread_representations_only", func(t *testing.T) {
				db := new_record_db(t, backend)
				defer tool.DB_close(db)

				user := "unwatch_user"
				seed_user_auth(db, t, user)

				seed_bbs_post(db, t, "-1", "42", "thread title 42")
				seed_bbs_post(db, t, "2", "42", "normal title 2-42")
				seed_bbs_post(db, t, "3", "42", "normal title 3-42")

				seed_user_watch(db, t, "bbs_watchlist", user, "-1-42")
				seed_user_watch(db, t, "thread_watchlist", user, "42")
				seed_user_watch(db, t, "bbs_watchlist", user, "2-42")
				seed_user_watch(db, t, "bbs_watchlist", user, "3-42")

				out := route.Api_bbs_watch_post(tool.Config{IP: user, Session: &emptySession{}}, "-1", "42")
				if out["response"] != "ok" {
					t.Fatalf("unwatch response = %v, want ok", out)
				}

				if count_user_watch_rows(db, "bbs_watchlist", user, "-1-42") != 0 {
					t.Fatalf("modern thread watch row remains after unwatch")
				}
				if count_user_watch_rows(db, "thread_watchlist", user, "42") != 0 {
					t.Fatalf("legacy thread watch row remains after unwatch")
				}
				if count_user_watch_rows(db, "bbs_watchlist", user, "2-42") == 0 {
					t.Fatalf("unrelated normal BBS 2-42 row was deleted by thread unwatch")
				}
				if count_user_watch_rows(db, "bbs_watchlist", user, "3-42") == 0 {
					t.Fatalf("unrelated normal BBS 3-42 row was deleted by thread unwatch")
				}

				items := bbs_watch_api_items(db, t, user, 1)
				for _, item := range items {
					if item.Set_id == "-1" && item.Set_code == "42" {
						t.Fatalf("thread 42 still appears after unwatch")
					}
				}
				if len(items) != 2 {
					t.Fatalf("items after unwatch = %d, want 2 unrelated normal entries", len(items))
				}
			})

			t.Run("account_separation", func(t *testing.T) {
				db := new_record_db(t, backend)
				defer tool.DB_close(db)

				user_a := "user_a"
				user_b := "user_b"
				other_user := "user_c"

				seed_user_auth(db, t, user_a)
				seed_user_auth(db, t, user_b)
				seed_user_auth(db, t, other_user)

				seed_bbs_post(db, t, "5", "1", "doc a1")
				seed_bbs_post(db, t, "5", "2", "doc a2")
				seed_bbs_post(db, t, "5", "3", "doc b1")

				seed_user_watch(db, t, "bbs_watchlist", user_a, "5-1")
				seed_user_watch(db, t, "bbs_watchlist", user_a, "5-2")
				seed_user_watch(db, t, "bbs_watchlist", user_b, "5-3")

				items_a := bbs_watch_api_items(db, t, user_a, 1)
				items_b := bbs_watch_api_items(db, t, user_b, 1)

				if len(items_a) != 2 {
					t.Fatalf("user_a items = %d, want 2", len(items_a))
				}
				if len(items_b) != 1 {
					t.Fatalf("user_b items = %d, want 1", len(items_b))
				}

				for _, item := range items_a {
					if item.Set_id == "5" && item.Set_code == "3" {
						t.Fatalf("user_a saw user_b watch entry 5-3")
					}
				}
				for _, item := range items_b {
					if item.Set_id == "5" && (item.Set_code == "1" || item.Set_code == "2") {
						t.Fatalf("user_b saw user_a watch entry")
					}
				}
			})

			t.Run("permission_check_unchanged", func(t *testing.T) {
				db := new_record_db(t, backend)
				defer tool.DB_close(db)

				owner := "watch_owner"
				viewer := "watch_viewer"
				seed_user_auth(db, t, owner)

				seed_bbs_post(db, t, "2", "99", "permission bbs 2-99")
				seed_user_watch(db, t, "bbs_watchlist", owner, "2-99")

				res := route.Api_user_watch_list(tool.Config{IP: viewer, Session: &emptySession{}}, owner, "1", "bbs_watchlist")
				if res["response"] != "require auth" {
					t.Fatalf("viewer response = %v, want require auth", res["response"])
				}
				data, ok := res["data"].([]route.User_watch_item)
				if !ok || len(data) != 0 {
					t.Fatalf("viewer data type or length = %v, want []route.User_watch_item{}", res["data"])
				}

				html := route.View_user_watch_list(tool.Config{IP: viewer, Session: &emptySession{}}, "1", "bbs_watchlist")
				if strings.Contains(html, `/bbs/w/2/99`) {
					t.Errorf("unauthenticated view rendered owner watch data")
				}
			})
		})
	}
}

func seed_bbs_post(db *sql.DB, t *testing.T, set_id, set_code, title string) {
	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed bbs %q: %v", q, err)
		}
	}
	exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`, "title", set_code, set_id, title)
	exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`, "date", set_code, set_id, "2026-01-01 00:00:00")
	exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`, "user_id", set_code, set_id, "bbs_author")
}

func bbs_watch_codes_zero(total int) []string {
	codes := make([]string, 0, total)
	for i := 1; i <= total; i++ {
		codes = append(codes, fmt.Sprintf("%03d", i))
	}
	return codes
}

func bbs_watch_codes_with_set(set_id string, codes []string) []string {
	out := make([]string, 0, len(codes))
	for _, code := range codes {
		out = append(out, set_id+"-"+code)
	}
	return out
}

func seed_bbs_watch_codes(db *sql.DB, t *testing.T, user, set_id string, codes []string) {
	for _, code := range codes {
		data := set_id + "-" + code
		if _, err := db.Exec(`insert into user_set (name, id, data) values (?, ?, ?)`, "bbs_watchlist", user, data); err != nil {
			t.Fatalf("seed bbs_watchlist %q: %v", data, err)
		}
	}
}

func seed_bbs_watch_codes_order(db *sql.DB, t *testing.T, user, set_id string, codes []int) {
	for _, code := range codes {
		data := set_id + "-" + fmt.Sprintf("%03d", code)
		if _, err := db.Exec(`insert into user_set (name, id, data) values (?, ?, ?)`, "bbs_watchlist", user, data); err != nil {
			t.Fatalf("seed bbs_watchlist %q: %v", data, err)
		}
	}
}

func bbs_watch_page_items(db *sql.DB, t *testing.T, user string, page int) []string {
	items := bbs_watch_api_items(db, t, user, page)
	ids := []string{}
	for _, item := range items {
		ids = append(ids, item.Data)
	}
	return ids
}

func bbs_watch_api_items(db *sql.DB, t *testing.T, user string, page int) []route.User_watch_item {
	res := route.Api_user_watch_list(tool.Config{IP: user, Session: &emptySession{}}, user, fmt.Sprint(page), "bbs_watchlist")

	response, ok := res["response"].(string)
	if !ok || response != "ok" {
		t.Fatalf("bbs watch api response = %v, want ok", res["response"])
	}

	items, ok := res["data"].([]route.User_watch_item)
	if !ok {
		t.Fatalf("bbs watch api data is not []route.User_watch_item: %v", res["data"])
	}
	return items
}

func collect_bbs_watch_pages(db *sql.DB, t *testing.T, user string) []string {
	all := []string{}
	page := 1
	for {
		items := bbs_watch_page_items(db, t, user, page)
		all = append(all, items...)
		if len(items) < 50 {
			break
		}
		page++
	}
	return all
}

func view_bbs_watch_page(db *sql.DB, t *testing.T, user string, page int) string {
	html := route.View_user_watch_list(tool.Config{IP: user, Session: &emptySession{}}, fmt.Sprint(page), "bbs_watchlist")
	if html == "" {
		t.Fatalf("view_bbs_watch_page returned empty html")
	}
	return html
}

func assert_page_links(t *testing.T, html string, page int, expect_prev bool, expect_next bool) {
	if page > 1 {
		prev_href := `href="/bbs_watch_list_page/` + strconv.Itoa(page-1) + `">`
		if expect_prev {
			if !strings.Contains(html, prev_href) {
				t.Fatalf("page %d is missing previous link %s", page, prev_href)
			}
		} else {
			if strings.Contains(html, prev_href) {
				t.Fatalf("page %d has an unexpected previous link %s", page, prev_href)
			}
		}
	}

	next_href := `href="/bbs_watch_list_page/` + strconv.Itoa(page+1) + `">`
	if expect_next {
		if !strings.Contains(html, next_href) {
			t.Fatalf("page %d is missing next link %s", page, next_href)
		}
	} else {
		if strings.Contains(html, next_href) {
			t.Fatalf("page %d has an unexpected next link %s", page, next_href)
		}
	}
}

func bbs_watch_min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func bbs_watch_max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func bbs_watch_sub(a, b int) int {
	if a > b {
		return a - b
	}
	return 0
}
