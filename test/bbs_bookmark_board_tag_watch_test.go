package test

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"

	"github.com/gin-gonic/gin"
)

func Test_bbs_bookmark_board_tag_watch_end_to_end(t *testing.T) {
	backends := []string{"sqlite", "turso"}
	if os.Getenv("MYSQL_AUDIT_DSN") != "" {
		backends = append(backends, "mysql")
	}

	for _, backend := range backends {
		t.Run(backend, func(t *testing.T) {
			db := new_record_db(t, backend)
			defer tool.DB_close(db)
			Seed_bbs_subscription_board(db, t, "88", "owner")

			outsider := tool.Config{IP: "subscription_outsider", Session: &emptySession{}}
			denied := route.Api_bbs_board_watch_post(outsider, "88")
			if denied["response"] != "require auth" {
				t.Fatalf("private board watch response = %#v, want require auth", denied)
			}
			if count := Bbs_subscription_count(db, t, "bbs_board_watchlist", outsider.IP, "88"); count != 0 {
				t.Fatalf("private-board watch created %d rows, want 0", count)
			}

			if _, err := db.Exec(`update bbs_set set set_data = 'normal' where set_id = '88' and set_name = 'bbs_view_acl'`); err != nil {
				t.Fatalf("make board public: %v", err)
			}
			watchers := []string{"watch_all", "board_only", "tag_only", "tag_both", "admin_ip"}
			for _, watcher := range watchers {
				config := tool.Config{IP: watcher, Session: &emptySession{}}
				if watcher == "watch_all" || watcher == "board_only" || watcher == "tag_both" || watcher == "admin_ip" {
					if result := route.Api_bbs_board_watch_post(config, "88"); result["response"] != "ok" {
						t.Fatalf("board watch for %q = %#v, want ok", watcher, result)
					}
				}
			}
			for watcher, tags := range map[string][]string{
				"watch_all": {"release", "security"},
				"tag_only":  {"release"},
				"tag_both":  {"release", "security"},
				"admin_ip":  {"release", "security"},
			} {
				config := tool.Config{IP: watcher, Session: &emptySession{}}
				for _, tag := range tags {
					if result := route.Api_bbs_tag_watch_post(config, "88", tag); result["response"] != "ok" {
						t.Fatalf("tag watch for %q / %q = %#v, want ok", watcher, tag, result)
					}
				}
			}
			if result := route.Api_bbs_tag_watch_post(outsider, "88", "bad,tag"); result["response"] != "error" {
				t.Fatalf("invalid tag watch response = %#v, want error", result)
			}

			author := tool.Config{IP: "admin_ip", Session: &emptySession{}}
			created := route.Api_bbs_w_edit_post(author, "88", "", "", "Tagged post", "post body", "", "release,security", "")
			if created["response"] != "ok" || created["data"] != "1" {
				t.Fatalf("create tagged post response = %#v, want ok/1", created)
			}
			for _, expected := range []struct {
				user  string
				count int
			}{
				{"watch_all", 1},
				{"board_only", 1},
				{"tag_only", 1},
				{"tag_both", 1},
				{"admin_ip", 0},
			} {
				if count := Bbs_user_notice_count(db, t, expected.user); count != expected.count {
					t.Errorf("notices for %q = %d, want %d", expected.user, count, expected.count)
				}
			}

			watcher_config := tool.Config{IP: "watch_all", Session: &emptySession{}}
			if result := route.Api_bbs_watch_post(watcher_config, "88", "1"); result["response"] != "ok" {
				t.Fatalf("comment subscription toggle = %#v, want ok", result)
			}
			if result := route.Api_bbs_bookmark_post(watcher_config, "88", "1"); result["response"] != "ok" {
				t.Fatalf("bookmark toggle = %#v, want ok", result)
			}
			if count := Bbs_subscription_count(db, t, "bbs_bookmark", watcher_config.IP, "88-1"); count != 1 {
				t.Fatalf("bookmark row count = %d, want 1", count)
			}
			if count := Bbs_subscription_count(db, t, "bbs_watchlist", watcher_config.IP, "88-1"); count != 1 {
				t.Fatalf("comment subscription row count = %d, want 1", count)
			}
			legacy_watch_list := route.Api_user_watch_list(watcher_config, watcher_config.IP, "1", "bbs_watchlist")
			legacy_json, err := json.Marshal(legacy_watch_list["data"])
			if err != nil {
				t.Fatalf("marshal legacy BBS watch list: %v", err)
			}
			if strings.Contains(string(legacy_json), `"Tag"`) {
				t.Fatalf("existing BBS watch-list JSON gained a Tag field: %s", legacy_json)
			}

			for _, watch_type := range []string{"bbs_bookmark", "bbs_board_watchlist", "bbs_tag_watchlist"} {
				result := route.Api_user_watch_list(watcher_config, watcher_config.IP, "1", watch_type)
				if result["response"] != "ok" {
					t.Fatalf("watch list %q response = %#v, want ok", watch_type, result)
				}
				items, ok := result["data"].([]route.User_watch_item)
				if !ok || len(items) == 0 {
					t.Fatalf("watch list %q data = %#v, want non-empty []User_watch_item", watch_type, result["data"])
				}
			}

			user_tool_page := route.View_user_safe(watcher_config, "")
			for _, path := range []string{`href="/bbs_bookmark_list"`, `href="/bbs_board_watch_list"`, `href="/bbs_tag_watch_list"`} {
				if !strings.Contains(user_tool_page, path) {
					t.Errorf("user tool page lacks management link %s", path)
				}
			}

			bookmark_list := route.View_user_watch_list(watcher_config, "1", "bbs_bookmark")
			if !strings.Contains(bookmark_list, `/bbs_bookmark/88/1`) || !strings.Contains(bookmark_list, `/bbs/w/88/1`) {
				t.Fatalf("bookmark management HTML lacks view/remove links: %s", bookmark_list)
			}
			board_list := route.View_user_watch_list(watcher_config, "1", "bbs_board_watchlist")
			if !strings.Contains(board_list, `/bbs_board_watch/88`) || !strings.Contains(board_list, `/bbs/in/88`) {
				t.Fatalf("board-watch management HTML lacks view/remove links: %s", board_list)
			}
			tag_list := route.View_user_watch_list(watcher_config, "1", "bbs_tag_watchlist")
			if !strings.Contains(tag_list, `/bbs_tag_watch/cmVsZWFzZQ`) || !strings.Contains(tag_list, `#release`) {
				t.Fatalf("global tag-watch management HTML lacks tag/remove links: %s", tag_list)
			}

			bookmark_form := route.View_bbs_bookmark(watcher_config, "88", "1", nil)
			if !strings.Contains(bookmark_form, `action="/bbs_bookmark/88/1"`) {
				t.Fatalf("bookmark confirmation lacks POST control: %s", bookmark_form)
			}
			gin.SetMode(gin.TestMode)
			context, _ := gin.CreateTestContext(httptest.NewRecorder())
			context.Request = httptest.NewRequest("GET", "/bbs/w/88/1", nil)
			context.Params = gin.Params{{Key: "comment_select", Value: ""}, {Key: "page", Value: "1"}}
			post_html := route.View_bbs_in_w(context, watcher_config, "88", "1")
			if !strings.Contains(post_html, `/bbs_bookmark/88/1`) || !strings.Contains(post_html, `/bbs_tag_watch/cmVsZWFzZQ`) {
				t.Fatalf("BBS post HTML lacks bookmark/tag watch controls: %s", post_html)
			}
			for _, label := range []string{"게시글 북마크", "게시글 댓글 주시", "태그 주시"} {
				if !strings.Contains(post_html, label) {
					t.Errorf("BBS post controls lack distinct label %q", label)
				}
			}
			board_html := route.View_bbs_in(watcher_config, "88", "1", "")
			if !strings.Contains(board_html, `bbs_board_watch/88`) {
				t.Fatalf("BBS board HTML lacks board-watch control: %s", board_html)
			}

			if result := route.Api_bbs_bookmark_post(watcher_config, "88", "1"); result["response"] != "ok" {
				t.Fatalf("bookmark removal = %#v, want ok", result)
			}
			if count := Bbs_subscription_count(db, t, "bbs_bookmark", watcher_config.IP, "88-1"); count != 0 {
				t.Fatalf("removed bookmark row count = %d, want 0", count)
			}
			if count := Bbs_subscription_count(db, t, "bbs_watchlist", watcher_config.IP, "88-1"); count != 1 {
				t.Fatalf("bookmark removal changed comment subscription: %d rows, want 1", count)
			}
			if result := route.Api_bbs_bookmark_post(watcher_config, "88", "1"); result["response"] != "ok" {
				t.Fatalf("restore bookmark = %#v, want ok", result)
			}

			if _, err := db.Exec(`update bbs_set set set_data = 'owner' where set_id = '88' and set_name = 'bbs_view_acl'`); err != nil {
				t.Fatalf("make board private: %v", err)
			}
			private_post := route.Api_bbs_w_edit_post(author, "88", "", "", "Private tagged post", "private body", "", "release,security", "")
			if private_post["response"] != "ok" || private_post["data"] != "2" {
				t.Fatalf("create private post response = %#v, want ok/2", private_post)
			}
			for _, user := range []string{"watch_all", "board_only", "tag_only", "tag_both"} {
				if count := Bbs_user_notice_count(db, t, user); count != 1 {
					t.Errorf("private post notification leaked to %q: total notices %d, want unchanged 1", user, count)
				}
			}
			for _, watch_type := range []string{"bbs_bookmark", "bbs_board_watchlist"} {
				result := route.Api_user_watch_list(watcher_config, watcher_config.IP, "1", watch_type)
				items, ok := result["data"].([]route.User_watch_item)
				if result["response"] != "ok" || !ok || len(items) == 0 {
					t.Errorf("inaccessible %s list = %#v, want removable placeholder", watch_type, result)
					continue
				}
				if items[0].Available || items[0].Display != "" {
					t.Errorf("inaccessible %s item exposed data: %#v", watch_type, items[0])
				}
			}
			hidden_bookmark_list := route.View_user_watch_list(watcher_config, "1", "bbs_bookmark")
			if strings.Contains(hidden_bookmark_list, "Tagged post") || !strings.Contains(hidden_bookmark_list, "접근할 수 없는 게시글") || !strings.Contains(hidden_bookmark_list, "/bbs_bookmark/88/1") {
				t.Errorf("hidden bookmark list should hide its title and keep removal control: %s", hidden_bookmark_list)
			}
			hidden_board_list := route.View_user_watch_list(watcher_config, "1", "bbs_board_watchlist")
			if strings.Contains(hidden_board_list, "Subscription board") || strings.Contains(hidden_board_list, `href="/bbs/in/88"`) || !strings.Contains(hidden_board_list, "/bbs_board_watch/88") {
				t.Errorf("inaccessible board list should hide its name and keep removal control: %s", hidden_board_list)
			}
			stale_bookmark_control := route.View_bbs_bookmark(watcher_config, "88", "1", nil)
			if !strings.Contains(stale_bookmark_control, "접근할 수 없는 게시글") || !strings.Contains(stale_bookmark_control, "게시글 북마크 해제") {
				t.Errorf("stale bookmark control should explain and allow removal: %s", stale_bookmark_control)
			}
			if result := route.Api_bbs_bookmark_post(watcher_config, "88", "1"); result["response"] != "ok" {
				t.Fatalf("remove inaccessible bookmark = %#v, want ok", result)
			}
			if count := Bbs_subscription_count(db, t, "bbs_bookmark", watcher_config.IP, "88-1"); count != 0 {
				t.Fatalf("inaccessible bookmark removal left %d rows", count)
			}
			if count := Bbs_subscription_count(db, t, "bbs_watchlist", watcher_config.IP, "88-1"); count != 1 {
				t.Fatalf("inaccessible bookmark removal changed comment subscription: %d rows", count)
			}
			if result := route.Api_bbs_board_watch_post(watcher_config, "88"); result["response"] != "ok" {
				t.Fatalf("remove inaccessible board watch = %#v, want ok", result)
			}
			if count := Bbs_subscription_count(db, t, "bbs_board_watchlist", watcher_config.IP, "88"); count != 0 {
				t.Fatalf("inaccessible board-watch removal left %d rows", count)
			}

			post_deleted := route.Api_bbs_w_delete(author, "88", "1")
			if post_deleted["response"] != "ok" {
				t.Fatalf("delete bookmarked post = %#v, want ok", post_deleted)
			}
			for _, watch_type := range []string{"bbs_bookmark", "bbs_watchlist"} {
				if count := Bbs_subscription_count(db, t, watch_type, watcher_config.IP, "88-1"); count != 0 {
					t.Errorf("post deletion left %s row: %d", watch_type, count)
				}
			}

			deleted := route.Api_bbs_delete(author, "88")
			if deleted["response"] != "ok" {
				t.Fatalf("delete subscription test board = %#v, want ok", deleted)
			}
			for _, entry := range []struct {
				name string
				data string
			}{
				{"bbs_bookmark", "88-1"},
				{"bbs_watchlist", "88-1"},
				{"bbs_board_watchlist", "88"},
			} {
				if count := Bbs_subscription_count(db, t, entry.name, watcher_config.IP, entry.data); count != 0 {
					t.Errorf("board deletion left %s row %q: %d", entry.name, entry.data, count)
				}
			}
			if count := Bbs_subscription_count(db, t, "bbs_tag_watchlist", watcher_config.IP, "release"); count != 1 {
				t.Errorf("deleting one board removed global tag watch: row count = %d, want 1", count)
			}

			Seed_bbs_subscription_board(db, t, "-1", "normal")
			thread_board_watcher := tool.Config{IP: "thread_board_watcher", Session: &emptySession{}}
			thread_tag_watcher := tool.Config{IP: "thread_tag_watcher", Session: &emptySession{}}
			thread_multi_watcher := tool.Config{IP: "thread_multi_watcher", Session: &emptySession{}}
			if result := route.Api_bbs_board_watch_post(thread_board_watcher, "-1"); result["response"] != "ok" {
				t.Fatalf("thread board watch = %#v, want ok", result)
			}
			if result := route.Api_bbs_tag_watch_post(thread_tag_watcher, "-1", "Discussion page"); result["response"] != "ok" {
				t.Fatalf("thread tag watch = %#v, want ok", result)
			}
			if result := route.Api_bbs_board_watch_post(thread_multi_watcher, "-1"); result["response"] != "ok" {
				t.Fatalf("thread board watch for combined watcher = %#v, want ok", result)
			}
			if result := route.Api_bbs_tag_watch_post(thread_multi_watcher, "-1", "Discussion page"); result["response"] != "ok" {
				t.Fatalf("thread tag watch for combined watcher = %#v, want ok", result)
			}
			if _, err := db.Exec(`insert into user_set (id, name, data) values (?, 'thread_watchlist', '1')`, thread_multi_watcher.IP); err != nil {
				t.Fatalf("seed legacy thread-comment watcher for deduplication: %v", err)
			}
			thread_post := route.Api_thread_bbs_post(author, "0", "Discussion page", "thread starter", "", "New discussion")
			if thread_post["response"] != "ok" || thread_post["topic_num"] != "1" {
				t.Fatalf("create thread topic response = %#v, want ok/1", thread_post)
			}
			if count := Bbs_user_notice_count(db, t, thread_board_watcher.IP); count != 1 {
				t.Errorf("thread board watcher notices = %d, want 1", count)
			}
			if count := Bbs_user_notice_count(db, t, thread_tag_watcher.IP); count != 1 {
				t.Errorf("thread tag watcher notices = %d, want 1", count)
			}
			if count := Bbs_user_notice_count(db, t, thread_multi_watcher.IP); count != 1 {
				t.Errorf("thread board+tag watcher notices = %d, want 1", count)
			}
		})
	}
}

func Test_bbs_tag_watch_global_scope_exactness_and_visibility(t *testing.T) {
	backends := []string{"sqlite", "turso"}
	if os.Getenv("MYSQL_AUDIT_DSN") != "" {
		backends = append(backends, "mysql")
	}

	for _, backend := range backends {
		t.Run(backend, func(t *testing.T) {
			db := new_record_db(t, backend)
			defer tool.DB_close(db)
			Seed_bbs_subscription_board(db, t, "88", "normal")
			Seed_bbs_subscription_board(db, t, "89", "normal")
			Seed_bbs_subscription_board(db, t, "90", "owner")

			global_watcher := tool.Config{IP: "global_tag_watcher", Session: &emptySession{}}
			multi_watcher := tool.Config{IP: "board_and_tag_watcher", Session: &emptySession{}}
			case_watcher := tool.Config{IP: "exact_case_watcher", Session: &emptySession{}}
			unicode_watcher := tool.Config{IP: "unicode_tag_watcher", Session: &emptySession{}}
			special_watcher := tool.Config{IP: "special_tag_watcher", Session: &emptySession{}}
			legacy_watcher := tool.Config{IP: "legacy_tag_watcher", Session: &emptySession{}}
			author := tool.Config{IP: "admin_ip", Session: &emptySession{}}

			for _, test_case := range []struct {
				config tool.Config
				board  string
				tag    string
			}{
				{global_watcher, "88", "release"},
				{multi_watcher, "88", "release"},
				{multi_watcher, "88", "설명"},
				{case_watcher, "88", "Release"},
				{case_watcher, "88", "release"},
				{unicode_watcher, "88", "설명"},
				{special_watcher, "88", "C++ & 'release'"},
				{author, "88", "release"},
			} {
				if result := route.Api_bbs_tag_watch_post(test_case.config, test_case.board, test_case.tag); result["response"] != "ok" {
					t.Fatalf("tag subscribe %q on board %s = %#v", test_case.tag, test_case.board, result)
				}
			}
			if result := route.Api_bbs_board_watch_post(multi_watcher, "89"); result["response"] != "ok" {
				t.Fatalf("board watch on second board = %#v", result)
			}

			status := route.Api_bbs_tag_watch_view(global_watcher, "89", "release")
			status_data, ok := status["data"].(map[string]string)
			if status["response"] != "ok" || !ok || status_data["watched"] != "true" {
				t.Fatalf("global tag state from another board = %#v, want watched", status)
			}
			user_tags := route.Api_user_watch_list(global_watcher, global_watcher.IP, "1", "bbs_tag_watchlist")
			items, ok := user_tags["data"].([]route.User_watch_item)
			if !ok || len(items) != 1 || items[0].Tag != "release" || items[0].Set_id != "" || items[0].Display != "#release" {
				t.Fatalf("global tag list item = %#v, want board-independent release", user_tags["data"])
			}

			if _, err := db.Exec(`update bbs_set set set_data = 'Cross & board' where set_id = '89' and set_name = 'bbs_name'`); err != nil {
				t.Fatalf("set escaped notification board name: %v", err)
			}
			created := route.Api_bbs_w_edit_post(author, "89", "", "", "Global & tag post", "post body", "", "release,설명,C++ & 'release'", "")
			if created["response"] != "ok" || created["data"] != "1" {
				t.Fatalf("create tagged post on second board = %#v", created)
			}
			for _, expected := range []struct {
				user  string
				count int
			}{
				{global_watcher.IP, 1},
				{multi_watcher.IP, 1},
				{case_watcher.IP, 1},
				{unicode_watcher.IP, 1},
				{special_watcher.IP, 1},
				{author.IP, 0},
			} {
				if count := Bbs_user_notice_count(db, t, expected.user); count != expected.count {
					t.Errorf("global tag notices for %q = %d, want %d", expected.user, count, expected.count)
				}
			}
			if result := route.Api_bbs_tag_watch_post(case_watcher, "", "Release"); result["response"] != "ok" {
				t.Fatalf("remove exact-case Release watch = %#v", result)
			}
			for _, test_case := range []struct {
				config tool.Config
				tag    string
				state  string
			}{{case_watcher, "Release", "false"}, {case_watcher, "release", "true"}} {
				state_result := route.Api_bbs_tag_watch_view(test_case.config, "", test_case.tag)
				state_data, _ := state_result["data"].(map[string]string)
				if state_result["response"] != "ok" || state_data["watched"] != test_case.state {
					t.Errorf("exact-case state for %q = %#v, want %s", test_case.tag, state_result, test_case.state)
				}
			}

			var notice_data string
			if err := db.QueryRow(`select data from user_notice where name = ?`, global_watcher.IP).Scan(&notice_data); err != nil {
				t.Fatalf("read global tag notification: %v", err)
			}
			if !strings.Contains(notice_data, "Cross &amp; board") || !strings.Contains(notice_data, "Global &amp; tag post") || strings.Contains(notice_data, "Cross & board") {
				t.Errorf("notification HTML was not safely escaped: %q", notice_data)
			}

			commented := route.Api_bbs_w_comment_post(author, "89", "1", "", "A comment does not create a new-post notification.")
			if commented["response"] != "ok" {
				t.Fatalf("create comment under watched post = %#v", commented)
			}
			edited := route.Api_bbs_w_edit_post(author, "89", "1", "", "Global tag post edited", "edited body", "", "release,설명", "")
			if edited["response"] != "ok" {
				t.Fatalf("edit tagged post = %#v", edited)
			}
			private_created := route.Api_bbs_w_edit_post(author, "90", "", "", "Private tagged post", "private body", "", "release", "")
			if private_created["response"] != "ok" {
				t.Fatalf("create private tagged post = %#v", private_created)
			}
			for _, user := range []string{global_watcher.IP, multi_watcher.IP, case_watcher.IP, unicode_watcher.IP} {
				if count := Bbs_user_notice_count(db, t, user); count != 1 {
					t.Errorf("private post, comment, or edit changed tag notices for %q: %d, want 1", user, count)
				}
			}

			if _, err := db.Exec(`insert into user_set (id, name, data) values (?, 'bbs_tag_watchlist', ?)`, legacy_watcher.IP, route.Bbs_tag_watch_key("88", "legacy")); err != nil {
				t.Fatalf("seed legacy-form tag watch: %v", err)
			}
			if _, err := db.Exec(`insert into user_set (id, name, data) values (?, 'bbs_tag_watchlist', 'legacy')`, legacy_watcher.IP); err != nil {
				t.Fatalf("seed canonical equivalent tag watch: %v", err)
			}
			for _, unrelated := range []struct {
				name string
				data string
			}{{"bbs_tag_watchlist", "other-tag"}, {"bbs_other_watchlist", "legacy"}} {
				if _, err := db.Exec(`insert into user_set (id, name, data) values (?, ?, ?)`, legacy_watcher.IP, unrelated.name, unrelated.data); err != nil {
					t.Fatalf("seed unrelated watcher row %q/%q: %v", unrelated.name, unrelated.data, err)
				}
			}
			legacy_status := route.Api_bbs_tag_watch_view(legacy_watcher, "", "legacy")
			legacy_data, ok := legacy_status["data"].(map[string]string)
			if legacy_status["response"] != "ok" || !ok || legacy_data["watched"] != "true" {
				t.Fatalf("legacy tag status = %#v, want watched", legacy_status)
			}
			legacy_list := route.View_user_watch_list(legacy_watcher, "1", "bbs_tag_watchlist")
			legacy_items_data := route.Api_user_watch_list(legacy_watcher, legacy_watcher.IP, "1", "bbs_tag_watchlist")
			legacy_items, ok := legacy_items_data["data"].([]route.User_watch_item)
			if !strings.Contains(legacy_list, "#legacy") || strings.Contains(legacy_list, `/bbs/in/88`) || !strings.Contains(legacy_list, "/bbs_tag_watch/") || legacy_items_data["response"] != "ok" || !ok || len(legacy_items) != 2 {
				t.Fatalf("legacy and canonical tags should deduplicate while unrelated tag remains: html=%s data=%#v", legacy_list, legacy_items_data)
			}
			legacy_post := route.Api_bbs_w_edit_post(author, "88", "", "", "Legacy scoped post", "legacy body", "", "legacy", "")
			if legacy_post["response"] != "ok" {
				t.Fatalf("create legacy-tag post on watched board = %#v", legacy_post)
			}
			if count := Bbs_user_notice_count(db, t, legacy_watcher.IP); count != 1 {
				t.Errorf("legacy-form tag notice count on first board = %d, want 1", count)
			}
			if deleted := route.Api_bbs_delete(author, "88"); deleted["response"] != "ok" {
				t.Fatalf("delete legacy source board = %#v", deleted)
			}
			if count := Bbs_subscription_count(db, t, "bbs_tag_watchlist", legacy_watcher.IP, route.Bbs_tag_watch_key("88", "legacy")); count != 1 {
				t.Errorf("deleting source board removed global legacy tag watch: %d", count)
			}
			if count := Bbs_subscription_count(db, t, "bbs_tag_watchlist", legacy_watcher.IP, "legacy"); count != 1 {
				t.Errorf("deleting source board removed canonical equivalent tag watch: %d", count)
			}
			legacy_elsewhere := route.Api_bbs_w_edit_post(author, "89", "", "", "Other legacy-tag post", "other body", "", "legacy", "")
			if legacy_elsewhere["response"] != "ok" {
				t.Fatalf("create legacy tag on another board = %#v", legacy_elsewhere)
			}
			if count := Bbs_user_notice_count(db, t, legacy_watcher.IP); count != 2 {
				t.Errorf("legacy-form and canonical tag watches should produce one notice per post: notices = %d, want 2", count)
			}
			if result := route.Api_bbs_tag_watch_post(legacy_watcher, "", "legacy"); result["response"] != "ok" {
				t.Fatalf("remove legacy tag watch globally = %#v", result)
			}
			if count := Bbs_subscription_count(db, t, "bbs_tag_watchlist", legacy_watcher.IP, route.Bbs_tag_watch_key("88", "legacy")); count != 0 {
				t.Errorf("tag toggle left legacy board-scoped record: %d", count)
			}
			if count := Bbs_subscription_count(db, t, "bbs_tag_watchlist", legacy_watcher.IP, "legacy"); count != 0 {
				t.Errorf("legacy removal left canonical record: %d", count)
			}
			for _, unrelated := range []struct {
				name string
				data string
			}{{"bbs_tag_watchlist", "other-tag"}, {"bbs_other_watchlist", "legacy"}} {
				if count := Bbs_subscription_count(db, t, unrelated.name, legacy_watcher.IP, unrelated.data); count != 1 {
					t.Errorf("removing legacy tag changed unrelated %s/%s row count to %d", unrelated.name, unrelated.data, count)
				}
			}
		})
	}
}

func Test_bbs_subscription_notifications_recheck_visibility_and_delivery_errors(t *testing.T) {
	backends := []string{"sqlite", "turso"}
	if os.Getenv("MYSQL_AUDIT_DSN") != "" {
		backends = append(backends, "mysql")
	}

	for _, backend := range backends {
		t.Run(backend, func(t *testing.T) {
			db := new_record_db(t, backend)
			defer tool.DB_close(db)
			Seed_bbs_subscription_board(db, t, "88", "normal")
			author := tool.Config{IP: "admin_ip", Session: &emptySession{}}
			watcher := tool.Config{IP: "visibility_watcher", Session: &emptySession{}}

			created := route.Api_bbs_w_edit_post(author, "88", "", "", "Visibility recheck", "body", "", "recheck", "")
			if created["response"] != "ok" || created["data"] != "1" {
				t.Fatalf("create visibility test post: %#v", created)
			}
			if result := route.Api_bbs_board_watch_post(watcher, "88"); result["response"] != "ok" {
				t.Fatalf("subscribe to board after post creation: %#v", result)
			}
			if result := route.Api_bbs_tag_watch_post(watcher, "", "recheck"); result["response"] != "ok" {
				t.Fatalf("subscribe to tag after post creation: %#v", result)
			}

			if _, err := db.Exec(`update bbs_set set set_data = 'owner' where set_id = '88' and set_name = 'bbs_view_acl'`); err != nil {
				t.Fatalf("revoke board visibility: %v", err)
			}
			route.Bbs_post_watch_notify(db, author, "88", "1", "Subscription board", "Visibility recheck", author.IP, []string{"recheck"})
			if count := Bbs_user_notice_count(db, t, watcher.IP); count != 0 {
				t.Fatalf("revoked board permission still received %d notifications", count)
			}

			if _, err := db.Exec(`update bbs_set set set_data = 'normal' where set_id = '88' and set_name = 'bbs_view_acl'`); err != nil {
				t.Fatalf("restore board visibility: %v", err)
			}
			if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values ('blind', '1', '88', 'O')`); err != nil {
				t.Fatalf("mark post blind: %v", err)
			}
			route.Bbs_post_watch_notify(db, author, "88", "1", "Subscription board", "Visibility recheck", author.IP, []string{"recheck"})
			if count := Bbs_user_notice_count(db, t, watcher.IP); count != 0 {
				t.Fatalf("blind post delivered %d notifications to an unauthorized watcher", count)
			}
			if _, err := db.Exec(`delete from bbs_data where set_name = 'blind' and set_code = '1' and set_id = '88'`); err != nil {
				t.Fatalf("unblind post: %v", err)
			}
			route.Bbs_post_watch_notify(db, author, "88", "1", "Subscription board", "Visibility recheck", author.IP, []string{"recheck"})
			if count := Bbs_user_notice_count(db, t, watcher.IP); count != 1 {
				t.Fatalf("visible post notification count = %d, want 1", count)
			}

			Seed_bbs_subscription_board(db, t, "-1", "owner")
			thread_seed_bbs(db, t, "25", "Legacy visibility topic")
			thread_watcher := "legacy_visibility_watcher"
			if _, err := db.Exec(`insert into user_set (id, name, data) values (?, 'thread_watchlist', '25')`, thread_watcher); err != nil {
				t.Fatalf("seed legacy thread watcher: %v", err)
			}
			route.Bbs_watch_notify(db, author, "-1", "25", "1", "thread", "Legacy visibility topic", "bbs_author", "")
			if count := Bbs_user_notice_count(db, t, thread_watcher); count != 0 {
				t.Fatalf("legacy thread watcher received %d notifications after access revocation", count)
			}
			if _, err := db.Exec(`update bbs_set set set_data = 'normal' where set_id = '-1' and set_name = 'bbs_view_acl'`); err != nil {
				t.Fatalf("restore legacy topic visibility: %v", err)
			}
			route.Bbs_watch_notify(db, author, "-1", "25", "1", "thread", "Legacy visibility topic", "bbs_author", "")
			if count := Bbs_user_notice_count(db, t, thread_watcher); count != 1 {
				t.Fatalf("visible legacy thread watcher notices = %d, want 1", count)
			}

			if backend == "sqlite" {
				for _, user_id := range []string{"delivery_fail", "delivery_success"} {
					if _, err := db.Exec(`insert into user_set (id, name, data) values (?, 'bbs_tag_watchlist', 'delivery')`, user_id); err != nil {
						t.Fatalf("seed delivery watcher %q: %v", user_id, err)
					}
				}
				if _, err := db.Exec(`create trigger fail_subscription_notice before insert on user_notice when new.name = 'delivery_fail' begin select raise(abort, 'forced delivery failure'); end`); err != nil {
					t.Fatalf("create delivery-failure trigger: %v", err)
				}
				post := route.Api_bbs_w_edit_post(author, "88", "", "", "Delivery failure isolation", "body", "", "delivery", "")
				if post["response"] != "ok" || post["data"] != "2" {
					t.Fatalf("delivery failure changed committed post result: %#v", post)
				}
				if count := Bbs_user_notice_count(db, t, "delivery_fail"); count != 0 {
					t.Errorf("failed recipient unexpectedly got %d notices", count)
				}
				if count := Bbs_user_notice_count(db, t, "delivery_success"); count != 1 {
					t.Errorf("healthy recipient notices = %d, want 1 despite another recipient's failure", count)
				}
			}
		})
	}
}

func Test_bbs_notification_database_call_rethrows_programming_panics(t *testing.T) {
	defer func() {
		recovered := recover()
		if recovered != "unexpected notification bug" {
			t.Fatalf("recovered panic = %#v, want original unexpected panic", recovered)
		}
	}()
	route.Bbs_notification_db_call(func() {
		panic("unexpected notification bug")
	})
}

func Test_bbs_notification_db_call_rethrows_error_valued_programming_panics(t *testing.T) {
	expected := errors.New("unexpected notification programming error")
	defer func() {
		if recovered := recover(); recovered != expected {
			t.Fatalf("recovered panic = %#v, want original error-valued panic", recovered)
		}
	}()
	route.Bbs_notification_db_call(func() {
		panic(expected)
	})
}

func Test_bbs_tag_watch_concurrent_cross_board_toggle(t *testing.T) {
	backends := []string{"sqlite", "turso"}
	if os.Getenv("MYSQL_AUDIT_DSN") != "" {
		backends = append(backends, "mysql")
	}

	for _, backend := range backends {
		t.Run(backend, func(t *testing.T) {
			db := new_record_db(t, backend)
			defer tool.DB_close(db)
			Seed_bbs_subscription_board(db, t, "88", "normal")
			Seed_bbs_subscription_board(db, t, "89", "normal")
			watcher := tool.Config{IP: "concurrent_tag_watcher", Session: &emptySession{}}

			toggle_concurrently := func(count int) {
				t.Helper()
				var wait_group sync.WaitGroup
				failures := make(chan string, count)
				for index := 0; index < count; index++ {
					wait_group.Add(1)
					go func(index int) {
						defer wait_group.Done()
						board_id := "88"
						if index%2 == 1 {
							board_id = "89"
						}
						if result := route.Api_bbs_tag_watch_post(watcher, board_id, "cross-board"); result["response"] != "ok" {
							failures <- fmt.Sprintf("board %s: %#v", board_id, result)
							return
						}
					}(index)
				}
				wait_group.Wait()
				close(failures)
				for failure := range failures {
					t.Errorf("concurrent toggle failed: %s", failure)
				}
			}

			toggle_concurrently(12)
			if count := Bbs_subscription_count(db, t, "bbs_tag_watchlist", watcher.IP, "cross-board"); count != 0 {
				t.Fatalf("even cross-board toggles left %d rows, want 0", count)
			}
			toggle_concurrently(13)
			if count := Bbs_subscription_count(db, t, "bbs_tag_watchlist", watcher.IP, "cross-board"); count != 1 {
				t.Fatalf("odd cross-board toggles left %d rows, want exactly 1", count)
			}

			independent_cases := []struct {
				user  string
				board string
				tag   string
			}{
				{"concurrent_tag_user_one", "88", "alpha"},
				{"concurrent_tag_user_one", "89", "beta"},
				{"concurrent_tag_user_two", "88", "alpha"},
				{"concurrent_tag_user_three", "89", "gamma"},
			}
			var independent_group sync.WaitGroup
			start_independent := make(chan struct{})
			independent_failures := make(chan string, len(independent_cases)*5)
			for _, test_case := range independent_cases {
				for attempt := 0; attempt < 5; attempt++ {
					independent_group.Add(1)
					go func(test_case struct {
						user  string
						board string
						tag   string
					}) {
						defer independent_group.Done()
						<-start_independent
						config := tool.Config{IP: test_case.user, Session: &emptySession{}}
						if result := route.Api_bbs_tag_watch_post(config, test_case.board, test_case.tag); result["response"] != "ok" {
							independent_failures <- fmt.Sprintf("user=%s board=%s tag=%q: %#v", test_case.user, test_case.board, test_case.tag, result)
						}
					}(test_case)
				}
			}
			close(start_independent)
			independent_group.Wait()
			close(independent_failures)
			for failure := range independent_failures {
				t.Errorf("independent concurrent tag toggle failed: %s", failure)
			}
			for _, test_case := range independent_cases {
				if count := Bbs_subscription_count(db, t, "bbs_tag_watchlist", test_case.user, test_case.tag); count != 1 {
					t.Errorf("independent odd toggles for user=%s tag=%q left %d rows, want one", test_case.user, test_case.tag, count)
				}
			}
		})
	}
}

func Test_bbs_subscription_management_pagination_and_stale_removal(t *testing.T) {
	db := new_record_db(t, "sqlite")
	defer tool.DB_close(db)
	watcher := tool.Config{IP: "pagination_watcher", Session: &emptySession{}}

	for index := 0; index < 53; index++ {
		for _, watch_type := range []string{"bbs_bookmark", "bbs_board_watchlist", "bbs_tag_watchlist"} {
			value := fmt.Sprintf("%d-1", 2000+index)
			if watch_type == "bbs_board_watchlist" {
				value = fmt.Sprintf("%d", 2000+index)
			} else if watch_type == "bbs_tag_watchlist" {
				value = fmt.Sprintf("page_%02d", index)
			}
			if _, err := db.Exec(`insert into user_set (id, name, data) values (?, ?, ?)`, watcher.IP, watch_type, value); err != nil {
				t.Fatalf("seed %s page row: %v", watch_type, err)
			}
		}
	}

	for _, watch_type := range []string{"bbs_bookmark", "bbs_board_watchlist", "bbs_tag_watchlist"} {
		first_page := route.Api_user_watch_list(watcher, watcher.IP, "1", watch_type)
		first_items, ok := first_page["data"].([]route.User_watch_item)
		if first_page["response"] != "ok" || !ok || len(first_items) != 50 || first_page["has_next"] != true {
			t.Errorf("first page of %s = %#v, want 50 items and next page", watch_type, first_page)
		}
		second_page := route.Api_user_watch_list(watcher, watcher.IP, "2", watch_type)
		second_items, ok := second_page["data"].([]route.User_watch_item)
		if second_page["response"] != "ok" || !ok || len(second_items) != 3 || second_page["has_next"] != false {
			t.Errorf("second page of %s = %#v, want 3 items and no next page", watch_type, second_page)
		}
	}
	tag_page := route.View_user_watch_list(watcher, "1", "bbs_tag_watchlist")
	if !strings.Contains(tag_page, "/bbs_tag_watch_list_page/2") {
		t.Errorf("first tag page lacks next link: %s", tag_page)
	}
	tag_last_page := route.View_user_watch_list(watcher, "2", "bbs_tag_watchlist")
	if strings.Contains(tag_last_page, "/bbs_tag_watch_list_page/3") {
		t.Errorf("last tag page unexpectedly links to a third page: %s", tag_last_page)
	}

	missing_board_state := route.Api_bbs_board_watch_view(watcher, "2001")
	missing_board_data, ok := missing_board_state["data"].(map[string]string)
	if missing_board_state["response"] != "ok" || !ok || missing_board_data["available"] != "false" {
		t.Errorf("stale board watch state = %#v, want removable unavailable state", missing_board_state)
	}
	if result := route.Api_bbs_board_watch_post(watcher, "2001"); result["response"] != "ok" {
		t.Errorf("remove stale board watch = %#v, want ok", result)
	}
	if result := route.Api_bbs_bookmark_post(watcher, "2001", "1"); result["response"] != "ok" {
		t.Errorf("remove stale bookmark = %#v, want ok", result)
	}
	if result := route.Api_bbs_bookmark_post(watcher, "invalid", "1"); result["response"] != "not exist" {
		t.Errorf("invalid bookmark identifier response = %#v, want not exist", result)
	}
	if result := route.Api_bbs_board_watch_post(watcher, "invalid"); result["response"] != "not exist" {
		t.Errorf("invalid board identifier response = %#v, want not exist", result)
	}
}

func Seed_bbs_subscription_board(db *sql.DB, t *testing.T, board_id string, view_acl string) {
	t.Helper()
	var anchor_count int
	if err := db.QueryRow(`select count(*) from other where name = 'ver'`).Scan(&anchor_count); err != nil {
		t.Fatalf("check global transaction anchor: %v", err)
	}
	if anchor_count == 0 {
		if _, err := db.Exec(`insert into other (name, data, coverage) values ('ver', 'subscription_test', '')`); err != nil {
			t.Fatalf("seed global transaction anchor: %v", err)
		}
	}
	for _, value := range []struct {
		name string
		data string
	}{
		{"bbs_name", "Subscription board"},
		{"bbs_view_acl", view_acl},
		{"bbs_edit_acl", "normal"},
		{"bbs_comment_acl", "normal"},
		{"bbs_type", "comment"},
	} {
		if _, err := db.Exec(`insert into bbs_set (set_name, set_code, set_id, set_data) values (?, '', ?, ?)`, value.name, board_id, value.data); err != nil {
			t.Fatalf("seed subscription board %s: %v", value.name, err)
		}
	}
}

func Bbs_subscription_count(db *sql.DB, t *testing.T, watch_type string, user string, value string) int {
	t.Helper()
	var count int
	if err := db.QueryRow(`select count(*) from user_set where name = ? and id = ? and data = ?`, watch_type, user, value).Scan(&count); err != nil {
		t.Fatalf("count %s rows for %q: %v", watch_type, user, err)
	}
	return count
}

func Bbs_user_notice_count(db *sql.DB, t *testing.T, user string) int {
	t.Helper()
	var count int
	if err := db.QueryRow(`select count(*) from user_notice where name = ?`, user).Scan(&count); err != nil {
		t.Fatalf("count notices for %q: %v", user, err)
	}
	return count
}

func Test_bbs_tag_watch_local_lock_without_report_board_settings(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso"} {
		t.Run(backend, func(t *testing.T) {
			db := new_record_db(t, backend)
			defer tool.DB_close(db)
			var version_rows int
			if err := db.QueryRow(`select count(*) from other where name = 'ver'`).Scan(&version_rows); err != nil {
				t.Fatalf("check required application version anchor: %v", err)
			}
			if version_rows == 0 {
				if _, err := db.Exec(`insert into other (name, data, coverage) values ('ver', 'subscription_lock_test', '')`); err != nil {
					t.Fatalf("seed required application version anchor: %v", err)
				}
			}
			var report_board_rows int
			if err := db.QueryRow(`select count(*) from bbs_set where set_id = '-2'`).Scan(&report_board_rows); err != nil {
				t.Fatalf("check optional report-board settings: %v", err)
			}
			if report_board_rows != 0 {
				t.Fatalf("test DB unexpectedly has %d report-board setting rows", report_board_rows)
			}
			watcher := tool.Config{IP: "anchor_watcher", Session: &emptySession{}}
			if result := route.Api_bbs_tag_watch_post(watcher, "", "anchor"); result["response"] != "ok" {
				t.Fatalf("global tag toggle with absent report-board settings = %#v", result)
			}
			if count := Bbs_subscription_count(db, t, "bbs_tag_watchlist", watcher.IP, "anchor"); count != 1 {
				t.Fatalf("global tag row count = %d, want 1", count)
			}
		})
	}
}

func Test_bbs_tag_watch_local_lock_uses_user_anchor(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso"} {
		t.Run(backend, func(t *testing.T) {
			db := new_record_db(t, backend)
			defer tool.DB_close(db)
			user_id := "scoped_tag_watcher"
			if _, err := db.Exec(`insert into user_set (id, name, data) values (?, 'acl', 'user')`, user_id); err != nil {
				t.Fatalf("seed stable user lock row: %v", err)
			}
			if _, err := db.Exec(`create trigger reject_global_tag_watch_anchor before update on other when old.name = 'ver' begin select raise(abort, 'global tag-watch lock used'); end`); err != nil {
				t.Fatalf("install global-anchor guard: %v", err)
			}

			watcher := tool.Config{IP: user_id, Session: &emptySession{}}
			if result := route.Api_bbs_tag_watch_post(watcher, "", "scoped"); result["response"] != "ok" {
				t.Fatalf("tag toggle with per-user anchor = %#v", result)
			}
			if count := Bbs_subscription_count(db, t, "bbs_tag_watchlist", user_id, "scoped"); count != 1 {
				t.Fatalf("global tag row count = %d, want 1", count)
			}
		})
	}
}

func Test_bbs_tag_watch_advisory_lock_scope(t *testing.T) {
	lock_name := func(database string, user string, tag string) string {
		return route.Bbs_identifier_lock_name(database, "tag-watch:"+user+":"+tag)
	}
	first := lock_name("subscription_db", "user_one", "Release")
	if len(first) > 64 || first != lock_name("subscription_db", "user_one", "Release") {
		t.Fatalf("same tag-watch scope produced unstable or oversized lock name %q", first)
	}
	for _, scope := range [][3]string{
		{"subscription_db", "user_two", "Release"},
		{"subscription_db", "user_one", "release"},
		{"other_db", "user_one", "Release"},
	} {
		if other := lock_name(scope[0], scope[1], scope[2]); other == first {
			t.Errorf("independent tag-watch scope %v shares advisory lock name %q", scope, first)
		}
	}
}

func Test_bbs_subscription_labels_localized(t *testing.T) {
	db := new_record_db(t, "sqlite")
	defer tool.DB_close(db)
	if _, err := db.Exec(`insert into other (name, data, coverage) values ('language', 'ko-KR', '')`); err != nil {
		t.Fatalf("set Korean test language: %v", err)
	}

	labels := []struct {
		key     string
		korean  string
		english string
	}{
		{"bbs_post_bookmark", "게시글 북마크", "Bookmark post"},
		{"bbs_post_bookmark_remove", "게시글 북마크 해제", "Remove post bookmark"},
		{"bbs_comment_watch", "게시글 댓글 주시", "Watch post comments"},
		{"bbs_comment_watch_remove", "게시글 댓글 주시 해제", "Stop watching post comments"},
		{"bbs_board_watch", "게시판 주시", "Watch board"},
		{"bbs_board_watch_remove", "게시판 주시 해제", "Stop watching board"},
		{"bbs_tag_watch", "태그 주시", "Watch tag"},
		{"bbs_tag_watch_remove", "태그 주시 해제", "Stop watching tag"},
		{"bbs_post_bookmark_list", "북마크한 게시글", "Bookmarked BBS posts"},
		{"bbs_board_watch_list", "주시 중인 게시판", "Watched BBS boards"},
		{"bbs_tag_watch_list", "주시 중인 태그", "Watched tags"},
		{"bbs_subscription_unavailable", "접근할 수 없는 게시글 또는 게시판", "Unavailable item"},
	}
	for _, label := range labels {
		if got := tool.Get_language(db, label.key, true); got != label.korean {
			t.Errorf("Korean %s = %q, want %q", label.key, got, label.korean)
		}
	}
	if _, err := db.Exec(`update other set data = 'en-US' where name = 'language'`); err != nil {
		t.Fatalf("set English test language: %v", err)
	}
	for _, label := range labels {
		if got := tool.Get_language(db, label.key, true); got != label.english {
			t.Errorf("English %s = %q, want %q", label.key, got, label.english)
		}
	}
}

func Test_bbs_subscription_tag_value_validation(t *testing.T) {
	for _, test_case := range []struct {
		tag   string
		valid bool
	}{
		{"release", true},
		{"tag with spaces", true},
		{"  trimmed  ", true},
		{"", false},
		{"one,two", false},
		{"one\ntwo", false},
	} {
		if got := route.Bbs_tag_watch_value_valid(test_case.tag); got != test_case.valid {
			t.Errorf("Bbs_tag_watch_value_valid(%q) = %t, want %t", test_case.tag, got, test_case.valid)
		}
	}

}
