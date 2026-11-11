package test

import (
	"database/sql"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

func bbs_best_comment_vote(db *sql.DB, t *testing.T, set_id string, set_code string, up string, down string) {
	t.Helper()
	for _, value := range []struct {
		name string
		data string
	}{
		{name: "tabom_count", data: up},
		{name: "tabom_down_count", data: down},
	} {
		if _, err := db.Exec(
			`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`,
			value.name,
			set_code,
			set_id,
			value.data,
		); err != nil {
			t.Fatalf("seed best comment vote %s: %v", value.name, err)
		}
	}
}

func bbs_best_comment_set_value(db *sql.DB, t *testing.T, set_id string, set_name string, value string) {
	t.Helper()
	if _, err := db.Exec(
		`insert into bbs_set (set_name, set_code, set_id, set_data) values (?, '', ?, ?)`,
		set_name,
		set_id,
		value,
	); err != nil {
		t.Fatalf("seed best comment setting %q: %v", set_name, err)
	}
}

func Test_Api_bbs_w_comment_best_selection_settings_and_privacy(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso"} {
		t.Run(backend, func(t *testing.T) {
			db := new_record_db(t, backend)
			defer tool.DB_close(db)

			const board_id = "71"
			bbs_comment_seed_board(db, t, board_id, "1", "BEST topic", "admin_ip")
			for _, comment_code := range []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"} {
				bbs_comment_seed_parent(db, t, board_id, "1", comment_code, "commenter")
			}
			for _, vote := range []struct {
				code string
				up   string
				down string
			}{
				{code: "1", up: "12", down: "2"},
				{code: "2", up: "15", down: "5"},
				{code: "3", up: "7", down: "1"},
				{code: "4", up: "8", down: "2"},
				{code: "5", up: "100", down: "0"},
				{code: "6", up: "99", down: "0"},
				{code: "7", up: "100", down: "0"},
				{code: "8", up: "5", down: "0"},
				{code: "9", up: "999999999999999999999999", down: "malformed"},
				{code: "10", up: "100", down: "0"},
				{code: "11", up: "100", down: "0"},
			} {
				bbs_best_comment_vote(db, t, board_id+"-1", vote.code, vote.up, vote.down)
			}
			for _, code := range []string{"3", "4"} {
				if _, err := db.Exec(
					`update bbs_data set set_data = '2026-01-03 00:00:00' where set_name = 'comment_date' and set_id = ? and set_code = ?`,
					board_id+"-1",
					code,
				); err != nil {
					t.Fatalf("set tied comment date: %v", err)
				}
			}
			for code, date := range map[string]string{"1": "2026-01-02 00:00:00", "2": "2026-01-01 00:00:00"} {
				if _, err := db.Exec(
					`update bbs_data set set_data = ? where set_name = 'comment_date' and set_id = ? and set_code = ?`,
					date,
					board_id+"-1",
					code,
				); err != nil {
					t.Fatalf("set ranked comment date: %v", err)
				}
			}
			if _, err := db.Exec(
				`update bbs_data set set_data = '' where set_name = 'comment' and set_id = ? and set_code = '10'`,
				board_id+"-1",
			); err != nil {
				t.Fatalf("tombstone comment 10: %v", err)
			}
			if _, err := db.Exec(
				`insert into bbs_data (set_name, set_code, set_id, set_data) values ('pinned', '5', ?, '2026-01-01')`,
				board_id+"-1",
			); err != nil {
				t.Fatalf("pin comment 5: %v", err)
			}
			if _, err := db.Exec(
				`insert into bbs_data (set_name, set_code, set_id, set_data) values ('blind', '6', ?, 'O')`,
				board_id+"-1",
			); err != nil {
				t.Fatalf("blind comment 6: %v", err)
			}
			if _, err := db.Exec(
				`update bbs_data set set_data = '   ' where set_name = 'comment' and set_id = ? and set_code = '7'`,
				board_id+"-1",
			); err != nil {
				t.Fatalf("blank comment 7: %v", err)
			}
			if _, err := db.Exec(
				`update bbs_data set set_data = '　' where set_name = 'comment' and set_id = ? and set_code = '11'`,
				board_id+"-1",
			); err != nil {
				t.Fatalf("unicode-whitespace-only comment 11: %v", err)
			}
			// A duplicate legacy body row must not create a second preview.
			if _, err := db.Exec(
				`insert into bbs_data (set_name, set_code, set_id, set_data) values ('comment', '1', ?, 'duplicate legacy body')`,
				board_id+"-1",
			); err != nil {
				t.Fatalf("duplicate comment body: %v", err)
			}

			viewer := tool.Config{IP: "viewer_ip", Session: &emptySession{}}
			admin := tool.Config{IP: "admin_ip", Session: &emptySession{}}
			result := route.Api_bbs_w_comment_best(viewer, board_id, "1")
			comments, ok := result["data"].([]map[string]string)
			if result["response"] != "ok" || !ok || len(comments) != 3 {
				t.Fatalf("default BEST response/data = %v / %T(%v), want ok and 3 comments", result["response"], result["data"], result["data"])
			}
			want_codes := []string{"2", "1", "3"}
			want_scores := []string{"10", "10", "6"}
			for index, want_code := range want_codes {
				if comments[index]["code"] != want_code {
					t.Errorf("BEST comment %d code = %q, want %q", index, comments[index]["code"], want_code)
				}
				if comments[index]["best_score"] != want_scores[index] {
					t.Errorf("BEST comment %q score = %q, want %q", want_code, comments[index]["best_score"], want_scores[index])
				}
			}
			for _, comment := range comments {
				if comment["code"] == "5" || comment["code"] == "6" || comment["code"] == "7" || comment["code"] == "9" || comment["code"] == "10" || comment["code"] == "11" {
					t.Errorf("ineligible comment %q appeared in BEST: %#v", comment["code"], comment)
				}
			}

			if route.Bbs_best_comment_enabled(db, board_id) != true || route.Bbs_best_comment_min_score(db, board_id) != 5 {
				t.Fatalf("missing board settings do not default to enabled/minimum 5")
			}
			admin_result := route.Api_bbs_w_comment_best(admin, board_id, "1")
			admin_comments, _ := admin_result["data"].([]map[string]string)
			if len(admin_comments) == 0 || admin_comments[0]["code"] != "6" {
				t.Errorf("comment manager did not receive the visible blinded comment: %#v", admin_comments)
			}
			read_setting := route.Api_bbs_w_set(viewer, board_id, "bbs_best_comment_min_score")
			if read_setting["response"] != "ok" {
				t.Fatalf("public BEST setting read response = %v, want ok", read_setting["response"])
			}
			invalid_settings_page := route.View_bbs_set(admin, board_id, url.Values{
				"bbs_best_comment_enabled":   {"0"},
				"bbs_best_comment_min_score": {"not-a-score"},
			})
			if !strings.Contains(invalid_settings_page, tool.Get_language(db, "error", true)) || !route.Bbs_best_comment_enabled(db, board_id) || route.Bbs_best_comment_min_score(db, board_id) != 5 {
				t.Errorf("invalid settings POST changed values or failed to show an error")
			}

			if response := route.Api_bbs_w_set_put(admin, board_id, "bbs_best_comment_min_score", "10", ""); response["response"] != "ok" {
				t.Fatalf("set BEST threshold response = %v, want ok", response["response"])
			}
			result = route.Api_bbs_w_comment_best(viewer, board_id, "1")
			comments, _ = result["data"].([]map[string]string)
			if len(comments) != 2 || comments[0]["code"] != "2" || comments[1]["code"] != "1" {
				t.Errorf("threshold 10 BEST = %#v, want comments 2 and 1 in creation-time order", comments)
			}

			if response := route.Api_bbs_w_set_put(admin, board_id, "bbs_best_comment_enabled", "0", ""); response["response"] != "ok" {
				t.Fatalf("disable BEST response = %v, want ok", response["response"])
			}
			result = route.Api_bbs_w_comment_best(viewer, board_id, "1")
			comments, _ = result["data"].([]map[string]string)
			if result["response"] != "ok" || len(comments) != 0 {
				t.Errorf("disabled BEST response/data = %v / %#v, want ok and empty", result["response"], comments)
			}
			if route.Api_bbs_w_set_put(admin, board_id, "bbs_best_comment_min_score", "-1", "")["response"] != "error" ||
				route.Api_bbs_w_set_put(admin, board_id, "bbs_best_comment_min_score", "not-a-score", "")["response"] != "error" ||
				route.Api_bbs_w_set_put(admin, board_id, "bbs_best_comment_min_score", "1000000001", "")["response"] != "error" ||
				route.Api_bbs_w_set_put(admin, board_id, "bbs_best_comment_enabled", "enabled", "")["response"] != "error" {
				t.Errorf("invalid BEST settings were not rejected")
			}
			if response := route.Api_bbs_w_set_put(viewer, board_id, "bbs_best_comment_enabled", "1", ""); response["response"] != "require auth" {
				t.Errorf("unauthorized BEST setting write response = %v, want require auth", response["response"])
			}

			bbs_best_comment_set_value(db, t, board_id, "bbs_only_my_data_view_acl", "owner")
			private_result := route.Api_bbs_w_comment_best(viewer, board_id, "1")
			if private_result["response"] != "require auth" {
				t.Errorf("private post BEST response = %v, want require auth", private_result["response"])
			}
			if _, err := db.Exec(
				`update bbs_set set set_data = 'normal' where set_name = 'bbs_only_my_data_view_acl' and set_id = ?`,
				board_id,
			); err != nil {
				t.Fatalf("restore public post visibility: %v", err)
			}
			bbs_best_comment_set_value(db, t, board_id, "bbs_view_acl", "owner")
			board_acl_result := route.Api_bbs_w_comment_best(viewer, board_id, "1")
			if board_acl_result["response"] != "require auth" {
				t.Errorf("board-ACL-protected BEST response = %v, want require auth", board_acl_result["response"])
			}
		})
	}
}

func Test_bbs_best_comment_page_link_and_rendering(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso"} {
		t.Run(backend, func(t *testing.T) {
			db := new_record_db(t, backend)
			defer tool.DB_close(db)

			const board_id = "72"
			bbs_comment_seed_board(db, t, board_id, "1", "Page navigation", "admin_ip")
			for comment_number := 1; comment_number <= 48; comment_number++ {
				bbs_comment_seed_parent(db, t, board_id, "1", strconv.Itoa(comment_number), "commenter")
			}
			if _, err := db.Exec(
				`insert into bbs_data (set_name, set_code, set_id, set_data) values ('comment_date', '49', ?, '2026-01-01 00:00:00')`,
				board_id+"-1",
			); err != nil {
				t.Fatalf("seed metadata-only comment node 49: %v", err)
			}
			bbs_comment_seed_parent(db, t, board_id+"-1", "49", "1", "reply_author")
			bbs_comment_seed_parent(db, t, board_id+"-1", "49", "2", "reply_author")
			bbs_comment_seed_parent(db, t, board_id, "1", "50", "commenter")
			bbs_best_comment_vote(db, t, board_id+"-1-49", "2", "30", "0")
			bbs_best_comment_vote(db, t, board_id+"-1", "50", "20", "0")
			bbs_best_comment_vote(db, t, board_id+"-1", "1", "40", "0")
			if _, err := db.Exec(
				`insert into bbs_data (set_name, set_code, set_id, set_data) values ('pinned', '1', ?, '2026-01-01')`,
				board_id+"-1",
			); err != nil {
				t.Fatalf("pin comment 1: %v", err)
			}

			viewer := tool.Config{IP: "viewer_ip", Session: &emptySession{}}
			page_one := route.Api_bbs_w_comment_page(viewer, board_id+"-1", true, "around", 1)
			page_one_comments, ok := page_one["data"].([]map[string]string)
			if page_one["response"] != "ok" || !ok || len(page_one_comments) != 50 || strings.TrimPrefix(page_one_comments[49]["id"]+"-"+page_one_comments[49]["code"], board_id+"-1-") != "49-1" {
				t.Fatalf("page one traversal end = response %v, count %d, want comment 49-1 at position 50", page_one["response"], len(page_one_comments))
			}
			page_two := route.Api_bbs_w_comment_page(viewer, board_id+"-1", true, "around", 2)
			page_two_comments, ok := page_two["data"].([]map[string]string)
			if page_two["response"] != "ok" || !ok || len(page_two_comments) != 2 || strings.TrimPrefix(page_two_comments[0]["id"]+"-"+page_two_comments[0]["code"], board_id+"-1-") != "49-2" || page_two_comments[1]["code"] != "50" {
				t.Fatalf("page two traversal = response %v, count %d, want comments 49-2 and 50", page_two["response"], len(page_two_comments))
			}

			best_result := route.Api_bbs_w_comment_best(viewer, board_id, "1")
			best_comments, ok := best_result["data"].([]map[string]string)
			if best_result["response"] != "ok" || !ok || len(best_comments) == 0 {
				t.Fatalf("BEST API response/data = %v / %T(%v), want ok and at least one preview", best_result["response"], best_result["data"], best_result["data"])
			}
			if strings.TrimPrefix(best_comments[0]["id"]+"-"+best_comments[0]["code"], board_id+"-1-") != "49-2" || best_comments[0]["best_page"] != "2" {
				t.Errorf("BEST target page = id %q code %q page %q, want 49-2 on page 2", best_comments[0]["id"], best_comments[0]["code"], best_comments[0]["best_page"])
			}
			if best_comments[1]["code"] != "50" || best_comments[1]["best_page"] != "2" {
				t.Errorf("BEST top-level page = code %q page %q, want comment 50 on page 2 after prior nested subtree", best_comments[1]["code"], best_comments[1]["best_page"])
			}

			comment_html := route.View_bbs_in_w_comment(db, viewer, board_id, "1", "", 1)
			want_link := `/bbs/w/72/1/page/2/comment/49-2#49-2`
			if !strings.Contains(comment_html, want_link) {
				t.Errorf("BEST preview does not link to original paginated comment %q", want_link)
			}
			best_heading := tool.Get_language(db, "bbs_best_comment", true)
			if !strings.Contains(comment_html, best_heading) || !strings.Contains(comment_html, `id="best-49-2"`) || !strings.Contains(comment_html, `id="pinned-1"`) {
				t.Errorf("page-one render lacks BEST/pinned preview markers")
			}
			if strings.Index(comment_html, `id="pinned-1"`) > strings.Index(comment_html, best_heading) {
				t.Errorf("BEST preview was not rendered after the pinned comment")
			}
			page_two_html := route.View_bbs_in_w_comment(db, viewer, board_id, "1", "", 2)
			if strings.Contains(page_two_html, best_heading) || strings.Contains(page_two_html, `id="best-49-2"`) {
				t.Errorf("page-two render unexpectedly duplicated first-page BEST previews")
			}

			admin := tool.Config{IP: "admin_ip", Session: &emptySession{}}
			settings_html := route.View_bbs_set(admin, board_id, nil)
			if !strings.Contains(settings_html, `name="bbs_best_comment_enabled"`) || !strings.Contains(settings_html, `name="bbs_best_comment_min_score"`) {
				t.Errorf("BBS settings form lacks BEST enable/threshold controls")
			}
			post_result := route.View_bbs_set(admin, board_id, url.Values{
				"bbs_view_acl":               {"normal"},
				"bbs_only_my_data_view_acl":  {"normal"},
				"bbs_acl":                    {"normal"},
				"bbs_edit_acl":               {"normal"},
				"bbs_comment_acl":            {"normal"},
				"bbs_markup":                 {""},
				"bbs_name":                   {"comment board"},
				"bbs_prefix":                 {""},
				"bbs_placeholder":            {""},
				"bbs_comment_placeholder":    {""},
				"bbs_excellent_min":          {""},
				"bbs_best_comment_enabled":   {"0"},
				"bbs_best_comment_min_score": {"8"},
				"bbs_slow_mode_seconds":      {""},
			})
			if !strings.Contains(post_result, "/bbs/set/72") {
				t.Errorf("valid BEST settings POST did not retain existing redirect behavior")
			}
			if route.Bbs_best_comment_enabled(db, board_id) || route.Bbs_best_comment_min_score(db, board_id) != 8 {
				t.Errorf("settings POST did not persist BEST enable/threshold values")
			}
			disabled_html := route.View_bbs_in_w_comment(db, viewer, board_id, "1", "", 1)
			if strings.Contains(disabled_html, tool.Get_language(db, "bbs_best_comment", true)) {
				t.Errorf("disabled BEST setting still renders a preview heading")
			}
		})
	}
}

func Test_bbs_best_comment_page_locator_ignores_unreachable_malformed_paths(t *testing.T) {
	db := new_record_db(t, "sqlite")
	defer tool.DB_close(db)

	const board_id = "73"
	bbs_comment_seed_board(db, t, board_id, "1", "Malformed path", "admin_ip")
	bbs_comment_seed_parent(db, t, board_id, "1", "1", "commenter")
	bbs_comment_seed_parent(db, t, board_id+"-1-99", "1", "1", "orphan")
	bbs_best_comment_vote(db, t, board_id+"-1-99", "1", "20", "0")

	viewer := tool.Config{IP: "viewer_ip", Session: &emptySession{}}
	result := route.Api_bbs_w_comment_best(viewer, board_id, "1")
	comments, ok := result["data"].([]map[string]string)
	if result["response"] != "ok" || !ok || len(comments) != 0 {
		t.Fatalf("unreachable malformed best comment response/data = %v / %#v, want ok and empty", result["response"], comments)
	}
}
