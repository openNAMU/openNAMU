package test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"

	"opennamu/route"
	"opennamu/route/tool"
)

func Test_Api_bbs_global_excellent_ranking_exclusion_and_cleanup(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso", "mysql"} {
		t.Run(backend, func(t *testing.T) {
			db := new_record_db(t, backend)
			defer tool.DB_close(db)
			config := tool.Config{IP: "admin_ip", Session: &emptySession{}}
			for _, threshold_case := range []struct {
				board_value  string
				global_value string
				want         int
			}{
				{board_value: "2", global_value: "5", want: 2},
				{board_value: "invalid", global_value: "5", want: 5},
				{board_value: "", global_value: "-1", want: 5},
			} {
				if got := route.Bbs_excellent_min_from_values(threshold_case.board_value, threshold_case.global_value); got != threshold_case.want {
					t.Fatalf("excellent threshold fallback (%q, %q) = %d, want %d", threshold_case.board_value, threshold_case.global_value, got, threshold_case.want)
				}
			}

			Bbs_global_excellent_seed_threshold(t, db, "5")
			Bbs_global_excellent_seed_board(t, db, "301", "normal", "2")
			Bbs_global_excellent_seed_board(t, db, "302", "normal", "invalid")
			Bbs_global_excellent_seed_board(t, db, "303", "normal", "")
			Bbs_global_excellent_seed_board(t, db, "304", "normal", "0")

			Bbs_global_excellent_seed_post(t, db, "301", "1", "custom threshold pinned", "admin_ip", 4, 1, true, false)
			Bbs_global_excellent_seed_post(t, db, "301", "2", "custom threshold boundary", "admin_ip", 2, 0, false, false)
			Bbs_global_excellent_seed_post(t, db, "302", "1", "global fallback", "admin_ip", 8, 2, false, false)
			Bbs_global_excellent_seed_post(t, db, "302", "2", "global fallback boundary", "admin_ip", 5, 0, false, false)
			Bbs_global_excellent_seed_post(t, db, "303", "1", "global default", "admin_ip", 7, 1, false, false)
			Bbs_global_excellent_seed_post(t, db, "304", "1", "zero threshold", "admin_ip", 2, 1, false, false)
			Bbs_global_excellent_seed_post(t, db, "-1", "1", "internal thread must not appear", "admin_ip", 100, 0, false, false)
			Bbs_global_excellent_seed_post(t, db, "-2", "1", "report must not appear", "admin_ip", 100, 0, false, false)
			Bbs_global_excellent_seed_post(t, db, "0", "1", "wiki comment must not appear", "admin_ip", 100, 0, false, false)

			first_page := route.Api_bbs(config, "", "1", "excellent")
			posts := Bbs_global_excellent_posts(t, first_page)
			want_keys := []string{"302-1", "303-1", "301-1", "304-1"}
			if len(posts) != len(want_keys) || first_page["has_next"] != false {
				t.Fatalf("global excellent result = %#v, has_next=%v; want %d eligible posts and no next page", posts, first_page["has_next"], len(want_keys))
			}
			for index, want_key := range want_keys {
				got_key := posts[index]["set_id"] + "-" + posts[index]["set_code"]
				if got_key != want_key {
					t.Fatalf("global excellent post %d = %q, want %q", index, got_key, want_key)
				}
			}
			if posts[2]["pinned"] != "1" {
				t.Fatalf("pinned global post state = %q, want 1", posts[2]["pinned"])
			}
			for _, post := range posts {
				if strings.Contains(post["title"], "must not appear") {
					t.Fatalf("special-board post leaked into global recommendations: %#v", post)
				}
			}

			board_page := route.Api_bbs(config, "301", "1", "excellent")
			board_posts := Bbs_global_excellent_posts(t, board_page)
			board_post_count := 0
			for _, post := range board_posts {
				if post["set_id"] == "301" && post["set_code"] == "1" {
					board_post_count++
				}
			}
			if board_post_count != 2 {
				t.Fatalf("existing pinned-plus-ordinary per-board behavior changed: %#v", board_posts)
			}

			if response := route.Api_bbs_excellent_exclusion_set(config, "301", "1", true); response["response"] != "ok" || response["data"] != "excluded" {
				t.Fatalf("exclude post response = %#v, want excluded", response)
			}
			if response := route.Api_bbs_excellent_exclusion_set(config, "301", "1", true); response["response"] != "ok" {
				t.Fatalf("repeat explicit exclude response = %#v, want ok", response)
			}
			if count := Bbs_global_excellent_marker_count(t, db, "301", "1"); count != 1 {
				t.Fatalf("repeated exclude left %d marker rows, want exactly one", count)
			}
			if Bbs_global_excellent_has_post(t, route.Api_bbs(config, "", "1", "excellent"), "301", "1") {
				t.Fatal("excluded post remained in the global excellent list")
			}
			if Bbs_global_excellent_has_post(t, route.Api_bbs(config, "301", "1", "excellent"), "301", "1") {
				t.Fatal("excluded pinned/ordinary post remained in the per-board excellent list")
			}
			if !Bbs_global_excellent_has_post(t, route.Api_bbs(config, "301", "1", "tabom"), "301", "1") {
				t.Fatal("excellent exclusion changed the existing tabom-sorted list")
			}

			if response := route.Api_bbs_excellent_exclusion_set(config, "301", "1", false); response["response"] != "ok" || response["data"] != "restored" {
				t.Fatalf("restore post response = %#v, want restored", response)
			}
			if response := route.Api_bbs_excellent_exclusion_set(config, "301", "1", false); response["response"] != "ok" {
				t.Fatalf("repeat explicit restore response = %#v, want ok", response)
			}
			if count := Bbs_global_excellent_marker_count(t, db, "301", "1"); count != 0 {
				t.Fatalf("restored post left %d marker rows, want none", count)
			}
			if count := Bbs_global_excellent_post_count(t, route.Api_bbs(config, "", "1", "excellent"), "301", "1"); count != 1 {
				t.Fatalf("restored global post appeared %d times, want once", count)
			}

			unauthorized := route.Api_bbs_excellent_exclusion_set(tool.Config{IP: "ordinary_ip"}, "301", "1", true)
			if unauthorized["response"] != "require auth" || Bbs_global_excellent_marker_count(t, db, "301", "1") != 0 {
				t.Fatalf("unauthorized exclusion response/state = %#v / %d marker rows", unauthorized, Bbs_global_excellent_marker_count(t, db, "301", "1"))
			}
			for _, target := range [][2]string{{"invalid", "1"}, {"301", "01"}, {"301", "999"}, {"999", "1"}} {
				response := route.Api_bbs_excellent_exclusion_set(config, target[0], target[1], true)
				if response["response"] != "not exist" {
					t.Errorf("tampered/missing identifier %q/%q response = %#v, want not exist", target[0], target[1], response)
				}
			}

			if response := route.Api_bbs_excellent_exclusion_set(config, "301", "1", true); response["response"] != "ok" {
				t.Fatalf("exclude before post deletion = %#v", response)
			}
			if response := route.Api_bbs_w_delete(config, "301", "1"); response["response"] != "ok" {
				t.Fatalf("post delete response = %#v", response)
			}
			if count := Bbs_global_excellent_marker_count(t, db, "301", "1"); count != 0 {
				t.Fatalf("post deletion retained %d exclusion markers", count)
			}

			if response := route.Api_bbs_excellent_exclusion_set(config, "304", "1", true); response["response"] != "ok" {
				t.Fatalf("exclude before board deletion = %#v", response)
			}
			if response := route.Api_bbs_delete(config, "304"); response["response"] != "ok" {
				t.Fatalf("board delete response = %#v", response)
			}
			if count := Bbs_global_excellent_marker_count(t, db, "304", "1"); count != 0 {
				t.Fatalf("board deletion retained %d exclusion markers", count)
			}
		})
	}
}

func Test_Api_bbs_global_excellent_visibility_and_pagination(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso", "mysql"} {
		t.Run(backend, func(t *testing.T) {
			db := new_record_db(t, backend)
			defer tool.DB_close(db)
			Bbs_global_excellent_seed_threshold(t, db, "5")
			Bbs_global_excellent_seed_permission_group(t, db, "excellent_member_group", "view", "bbs_use")
			if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl', 'excellent_member', 'excellent_member_group')`); err != nil {
				t.Fatalf("seed member account: %v", err)
			}
			member := tool.Config{IP: "excellent_member"}

			Bbs_global_excellent_seed_board(t, db, "401", "normal", "")
			Bbs_global_excellent_seed_board(t, db, "402", "owner", "")
			Bbs_global_excellent_seed_board(t, db, "403", "normal", "")
			if _, err := db.Exec(`insert into bbs_set (set_name, set_code, set_id, set_data) values ('bbs_only_my_data_view_acl', '', '403', 'owner')`); err != nil {
				t.Fatalf("seed private-post ACL: %v", err)
			}

			for post_code := 1; post_code <= 50; post_code++ {
				Bbs_global_excellent_seed_post(t, db, "401", fmt.Sprint(post_code), fmt.Sprintf("visible post %02d", post_code), "other_author", 6, 0, false, false)
			}
			Bbs_global_excellent_seed_post(t, db, "402", "1", "secret board title must not leak", "other_author", 100, 0, false, false)
			Bbs_global_excellent_seed_post(t, db, "403", "1", "private post title must not leak", "other_author", 100, 0, false, false)
			Bbs_global_excellent_seed_post(t, db, "403", "2", "member owned private post", "excellent_member", 6, 0, false, false)
			Bbs_global_excellent_seed_post(t, db, "401", "51", "blind post title must not leak", "other_author", 100, 0, false, true)

			first_result := route.Api_bbs(member, "", "1", "excellent")
			first_page := Bbs_global_excellent_posts(t, first_result)
			if len(first_page) != 50 || first_result["has_next"] != true {
				t.Fatalf("member page 1 = %d posts, has_next=%v; want exactly 50 visible posts and a next page", len(first_page), first_result["has_next"])
			}
			for _, post := range first_page {
				if strings.Contains(post["title"], "must not leak") || post["title"] == "member owned private post" {
					t.Fatalf("inaccessible/second-page title leaked on page 1: %#v", post)
				}
			}
			second_result := route.Api_bbs(member, "", "2", "excellent")
			second_page := Bbs_global_excellent_posts(t, second_result)
			if len(second_page) != 1 || second_result["has_next"] != false || second_page[0]["title"] != "member owned private post" {
				t.Fatalf("member page 2 = %#v, has_next=%v; want only their visible private post", second_page, second_result["has_next"])
			}

			Bbs_global_excellent_seed_board(t, db, "410", "normal", "")
			for post_code := 1; post_code <= 55; post_code++ {
				Bbs_global_excellent_seed_post(t, db, "410", fmt.Sprint(post_code), fmt.Sprintf("page boundary %02d", post_code), "admin_ip", 6, 0, false, false)
			}
			page_one_result := route.Api_bbs(member, "", "1", "excellent")
			page_two_result := route.Api_bbs(member, "", "2", "excellent")
			page_three_result := route.Api_bbs(member, "", "3", "excellent")
			page_one := Bbs_global_excellent_posts(t, page_one_result)
			page_one_again := Bbs_global_excellent_posts(t, route.Api_bbs(member, "", "1", "excellent"))
			page_two := Bbs_global_excellent_posts(t, page_two_result)
			page_three := Bbs_global_excellent_posts(t, page_three_result)
			if len(page_one) != 50 || len(page_two) != 50 || len(page_three) != 6 || page_one_result["has_next"] != true || page_two_result["has_next"] != true || page_three_result["has_next"] != false {
				t.Fatalf("pagination lengths/statuses = %d/%v, %d/%v, %d/%v; want 50/true, 50/true, 6/false", len(page_one), page_one_result["has_next"], len(page_two), page_two_result["has_next"], len(page_three), page_three_result["has_next"])
			}
			if fmt.Sprint(Bbs_global_excellent_posts_tie_keys(page_one)) != fmt.Sprint(Bbs_global_excellent_posts_tie_keys(page_one_again)) {
				t.Fatalf("equal-score order is unstable: %v vs %v", Bbs_global_excellent_posts_tie_keys(page_one), Bbs_global_excellent_posts_tie_keys(page_one_again))
			}
			seen := map[string]bool{}
			for _, post := range append(append(page_one, page_two...), page_three...) {
				key := post["set_id"] + "-" + post["set_code"]
				if seen[key] {
					t.Fatalf("post %q appeared more than once across pages", key)
				}
				seen[key] = true
			}
			if len(seen) != 106 {
				t.Fatalf("combined page result count = %d, want 106 eligible posts (50 public + 1 member-owned + 55 boundary rows)", len(seen))
			}
		})
	}
}

func Test_bbs_global_excellent_concurrent_explicit_exclusion_is_idempotent(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso", "mysql"} {
		t.Run(backend, func(t *testing.T) {
			db := new_record_db(t, backend)
			defer tool.DB_close(db)
			Bbs_global_excellent_seed_board(t, db, "501", "normal", "")
			Bbs_global_excellent_seed_post(t, db, "501", "1", "concurrent moderation target", "admin_ip", 6, 0, false, false)

			const workers = 8
			start := make(chan struct{})
			results := make(chan map[string]any, workers)
			var wait_group sync.WaitGroup
			for worker := 0; worker < workers; worker++ {
				wait_group.Add(1)
				go func() {
					defer wait_group.Done()
					<-start
					results <- route.Api_bbs_excellent_exclusion_set(tool.Config{IP: "admin_ip"}, "501", "1", true)
				}()
			}
			close(start)
			wait_group.Wait()
			close(results)
			for result := range results {
				if result["response"] != "ok" {
					t.Fatalf("concurrent exclusion response = %#v, want ok", result)
				}
			}
			if count := Bbs_global_excellent_marker_count(t, db, "501", "1"); count != 1 {
				t.Fatalf("concurrent explicit exclusions left %d marker rows, want exactly one", count)
			}
		})
	}
}

func Test_bbs_global_excellent_registered_http_routes_and_management_form(t *testing.T) {
	db := new_record_db(t, "sqlite")
	defer tool.DB_close(db)
	Bbs_global_excellent_seed_board(t, db, "811", "normal", "")
	Bbs_global_excellent_seed_post(t, db, "811", "1", "HTTP recommended <post>", "admin_ip", 8, 0, false, false)
	Bbs_global_excellent_seed_permission_group(t, db, "excellent_moderator_group", "view", "bbs")
	Bbs_global_excellent_seed_permission_group(t, db, "excellent_member_group", "view", "bbs_use")
	for _, user := range []struct {
		id    string
		group string
	}{
		{id: "excellent_moderator", group: "excellent_moderator_group"},
		{id: "excellent_member", group: "excellent_member_group"},
	} {
		if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl', ?, ?)`, user.id, user.group); err != nil {
			t.Fatalf("seed HTTP user %q: %v", user.id, err)
		}
	}

	previous_mode := gin.Mode()
	gin.SetMode(gin.TestMode)
	t.Cleanup(func() { gin.SetMode(previous_mode) })
	t.Setenv("NAMU_SESSION_KEY", "global_bbs_excellent_http_test_session_key")
	session_middleware := tool.Session_middleware()
	moderator_cookie := Bbs_subscription_http_session_cookie(t, "excellent_moderator", session_middleware)
	member_cookie := Bbs_subscription_http_session_cookie(t, "excellent_member", session_middleware)

	http_router := gin.New()
	http_router.Use(session_middleware)
	http_router.Use(route.Site_view_middleware())
	route.Register_bbs_excellent_routes(http_router)
	route.Register_bbs_excellent_api_routes(http_router)
	route.Register_bbs_post_tool_routes(http_router)
	server := httptest.NewServer(http_router)
	t.Cleanup(server.Close)

	status, body := Bbs_subscription_http_request(t, server.URL, http.MethodGet, "/bbs/excellent", moderator_cookie, nil)
	if status != http.StatusOK || !strings.Contains(string(body), "HTTP recommended &lt;post&gt;") || !strings.Contains(string(body), tool.Get_language(db, "global_excellent_posts", true)) {
		t.Fatalf("global recommendation page = %d %s; want escaped post and localized heading", status, body)
	}
	status, body = Bbs_subscription_http_request(t, server.URL, http.MethodGet, "/bbs/excellent/page/1", moderator_cookie, nil)
	if status != http.StatusOK || !strings.Contains(string(body), "HTTP recommended &lt;post&gt;") {
		t.Fatalf("paginated recommendation route = %d %s", status, body)
	}
	status, body = Bbs_subscription_http_request(t, server.URL, http.MethodGet, "/api/v2/bbs/excellent/1", moderator_cookie, nil)
	var global_api_response map[string]any
	if status != http.StatusOK || json.Unmarshal(body, &global_api_response) != nil || global_api_response["response"] != "ok" {
		t.Fatalf("global recommendation API = %d %s, want response ok", status, body)
	}

	status, body = Bbs_subscription_http_request(t, server.URL, http.MethodGet, "/bbs/tool/811/1", moderator_cookie, nil)
	if status != http.StatusOK || !strings.Contains(string(body), tool.Get_language(db, "excellent_post_exclude", true)) || !strings.Contains(string(body), tool.Get_language(db, "excellent_post_included_status", true)) {
		t.Fatalf("moderation form before exclusion = %d %s", status, body)
	}
	status, body = Bbs_subscription_http_request(t, server.URL, http.MethodGet, "/bbs/tool/811/1", member_cookie, nil)
	if status != http.StatusOK || strings.Contains(string(body), `name="action" value="excellent_exclusion"`) {
		t.Fatalf("non-moderator post tool = %d %s, want no exclusion control", status, body)
	}
	member_form := url.Values{"action": {"excellent_exclusion"}, "excluded": {"1"}}
	status, _ = Bbs_subscription_http_request(t, server.URL, http.MethodPost, "/bbs/tool/811/1", member_cookie, member_form)
	if status != http.StatusOK || Bbs_global_excellent_marker_count(t, db, "811", "1") != 0 {
		t.Fatalf("non-moderator tool POST changed state: status=%d markers=%d", status, Bbs_global_excellent_marker_count(t, db, "811", "1"))
	}
	form := url.Values{"action": {"excellent_exclusion"}, "excluded": {"1"}}
	status, body, location := Bbs_subscription_http_request_with_location(t, server.URL, http.MethodPost, "/bbs/tool/811/1", moderator_cookie, form)
	if status != http.StatusFound || !strings.Contains(location, "/bbs/tool/811/1") || Bbs_global_excellent_marker_count(t, db, "811", "1") != 1 {
		t.Fatalf("moderation POST = %d %s location=%q, markers=%d", status, body, location, Bbs_global_excellent_marker_count(t, db, "811", "1"))
	}
	status, body = Bbs_subscription_http_request(t, server.URL, http.MethodGet, "/bbs/tool/811/1", moderator_cookie, nil)
	if status != http.StatusOK || !strings.Contains(string(body), tool.Get_language(db, "excellent_post_restore", true)) || !strings.Contains(string(body), tool.Get_language(db, "excellent_post_excluded_status", true)) {
		t.Fatalf("moderation form after exclusion = %d %s", status, body)
	}
	status, body = Bbs_subscription_http_request(t, server.URL, http.MethodGet, "/bbs/excellent", moderator_cookie, nil)
	if status != http.StatusOK || strings.Contains(string(body), "HTTP recommended &lt;post&gt;") || !strings.Contains(string(body), tool.Get_language(db, "no_excellent_posts", true)) {
		t.Fatalf("excluded post remained on global page = %d %s", status, body)
	}

	form.Set("excluded", "0")
	status, body, location = Bbs_subscription_http_request_with_location(t, server.URL, http.MethodPost, "/bbs/tool/811/1", moderator_cookie, form)
	if status != http.StatusFound || !strings.Contains(location, "/bbs/tool/811/1") || Bbs_global_excellent_marker_count(t, db, "811", "1") != 0 {
		t.Fatalf("restore moderation POST = %d %s location=%q, markers=%d", status, body, location, Bbs_global_excellent_marker_count(t, db, "811", "1"))
	}

	status, body = Bbs_subscription_http_request(t, server.URL, http.MethodPost, "/api/v2/bbs/excellent_exclusion/811/1/exclude", member_cookie, nil)
	Bbs_subscription_http_expect_response(t, status, body, "require auth")
	if count := Bbs_global_excellent_marker_count(t, db, "811", "1"); count != 0 {
		t.Fatalf("unauthorized HTTP API request created %d exclusion markers", count)
	}
	status, body = Bbs_subscription_http_request(t, server.URL, http.MethodPost, "/api/v2/bbs/excellent_exclusion/811/1/exclude", moderator_cookie, nil)
	Bbs_subscription_http_expect_response(t, status, body, "ok")
	if count := Bbs_global_excellent_marker_count(t, db, "811", "1"); count != 1 {
		t.Fatalf("authorized HTTP API request created %d markers, want one", count)
	}
	status, body = Bbs_subscription_http_request(t, server.URL, http.MethodPost, "/api/v2/bbs/excellent_exclusion/811/1/restore", moderator_cookie, nil)
	Bbs_subscription_http_expect_response(t, status, body, "ok")
	if count := Bbs_global_excellent_marker_count(t, db, "811", "1"); count != 0 {
		t.Fatalf("HTTP API restore left %d markers", count)
	}
	status, body = Bbs_subscription_http_request(t, server.URL, http.MethodPost, "/api/v2/bbs/excellent_exclusion/811/1/invalid", moderator_cookie, nil)
	Bbs_subscription_http_expect_response(t, status, body, "error")
	if count := Bbs_global_excellent_marker_count(t, db, "811", "1"); count != 0 {
		t.Fatalf("invalid HTTP API state changed %d markers", count)
	}
	status, _ = Bbs_subscription_http_request(t, server.URL, http.MethodGet, "/api/v2/bbs/excellent_exclusion/811/1/exclude", moderator_cookie, nil)
	if status == http.StatusOK {
		t.Fatalf("mutation API accepted GET with status 200")
	}
}

func Bbs_global_excellent_seed_threshold(t *testing.T, db *sql.DB, value string) {
	t.Helper()
	if _, err := db.Exec(`delete from other where name = 'bbs_excellent_min' and coverage = ''`); err != nil {
		t.Fatalf("clear global excellent threshold: %v", err)
	}
	if _, err := db.Exec(`insert into other (name, data, coverage) values ('bbs_excellent_min', ?, '')`, value); err != nil {
		t.Fatalf("seed global excellent threshold: %v", err)
	}
}

func Bbs_global_excellent_seed_board(t *testing.T, db *sql.DB, board_id string, view_acl string, threshold string) {
	t.Helper()
	Seed_bbs_subscription_board(db, t, board_id, view_acl)
	if threshold != "" {
		if _, err := db.Exec(`insert into bbs_set (set_name, set_code, set_id, set_data) values ('bbs_excellent_min', '', ?, ?)`, board_id, threshold); err != nil {
			t.Fatalf("seed board %s excellent threshold %q: %v", board_id, threshold, err)
		}
	}
}

func Bbs_global_excellent_seed_permission_group(t *testing.T, db *sql.DB, group string, permissions ...string) {
	t.Helper()
	for _, permission := range permissions {
		if _, err := db.Exec(`insert into alist (name, acl) values (?, ?)`, group, permission); err != nil {
			t.Fatalf("seed permission %q for group %q: %v", permission, group, err)
		}
	}
}

func Bbs_global_excellent_seed_post(t *testing.T, db *sql.DB, board_id string, post_code string, title string, author string, up int, down int, pinned bool, blind bool) {
	t.Helper()
	values := []struct {
		name string
		data string
	}{
		{name: "title", data: title},
		{name: "date", data: "2026-02-01 12:00:00"},
		{name: "user_id", data: author},
		{name: "tabom_count", data: fmt.Sprint(up)},
		{name: "tabom_down_count", data: fmt.Sprint(down)},
	}
	if blind {
		values = append(values, struct {
			name string
			data string
		}{name: "blind", data: "O"})
	}
	for _, value := range values {
		if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`, value.name, post_code, board_id, value.data); err != nil {
			t.Fatalf("seed %s for %s/%s: %v", value.name, board_id, post_code, err)
		}
	}
	if pinned {
		if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values ('pinned', ?, ?, '2026-02-01 12:00:00')`, post_code, board_id); err != nil {
			t.Fatalf("seed pinned state for %s/%s: %v", board_id, post_code, err)
		}
	}
}

func Bbs_global_excellent_posts(t *testing.T, response map[string]any) []map[string]string {
	t.Helper()
	if response["response"] != "ok" {
		t.Fatalf("global excellent response = %#v, want ok", response)
	}
	posts, ok := response["data"].([]map[string]string)
	if !ok {
		t.Fatalf("global excellent data type = %T, want []map[string]string", response["data"])
	}
	return posts
}

func Bbs_global_excellent_posts_tie_keys(posts []map[string]string) []string {
	keys := make([]string, 0, len(posts))
	for _, post := range posts {
		keys = append(keys, post["set_id"]+"-"+post["set_code"])
	}
	return keys
}

func Bbs_global_excellent_has_post(t *testing.T, response map[string]any, set_id string, set_code string) bool {
	t.Helper()
	return Bbs_global_excellent_post_count(t, response, set_id, set_code) > 0
}

func Bbs_global_excellent_post_count(t *testing.T, response map[string]any, set_id string, set_code string) int {
	t.Helper()
	count := 0
	for _, post := range Bbs_global_excellent_posts(t, response) {
		if post["set_id"] == set_id && post["set_code"] == set_code {
			count++
		}
	}
	return count
}

func Bbs_global_excellent_marker_count(t *testing.T, db *sql.DB, set_id string, set_code string) int {
	t.Helper()
	var count int
	if err := db.QueryRow(`select count(*) from bbs_data where set_name = 'excellent_excluded' and set_id = ? and set_code = ?`, set_id, set_code).Scan(&count); err != nil {
		t.Fatalf("count exclusion markers for %s/%s: %v", set_id, set_code, err)
	}
	return count
}
