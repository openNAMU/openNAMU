package test

import (
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

func Test_Api_bbs_sort_modes(t *testing.T) {
	db := apiSlowDB(t)
	defer tool.DB_close(db)

	const board_id = "41"
	posts := []struct {
		code          string
		views         string
		comments      string
		tabom         string
		tabom_down    string
		last_activity string
	}{
		{code: "1", views: "5", comments: "1", tabom: "7", tabom_down: "1", last_activity: "2026-01-01 00:00:00"},
		{code: "2", views: "20", comments: "3", tabom: "12", tabom_down: "1", last_activity: "2026-01-03 00:00:00"},
		{code: "3", views: "10", comments: "2", tabom: "9", tabom_down: "1", last_activity: "2026-01-02 00:00:00"},
	}
	for _, post := range posts {
		bbs_comment_seed_board(db, t, board_id, post.code, "post "+post.code, "admin_ip")
		for _, value := range []struct {
			name string
			data string
		}{
			{name: "view_count", data: post.views},
			{name: "comment_count", data: post.comments},
			{name: "tabom_count", data: post.tabom},
			{name: "tabom_down_count", data: post.tabom_down},
			{name: "last_activity", data: post.last_activity},
			{name: "tag", data: "phase5"},
		} {
			if _, err := db.Exec(
				`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`,
				value.name,
				post.code,
				board_id,
				value.data,
			); err != nil {
				t.Fatalf("seed %s for post %s: %v", value.name, post.code, err)
			}
		}
	}

	config := tool.Config{IP: "admin_ip", Session: &emptySession{}}
	for _, test_case := range []struct {
		sort_type string
		want_code string
	}{
		{sort_type: "date", want_code: "3"},
		{sort_type: "activity", want_code: "2"},
		{sort_type: "view", want_code: "2"},
		{sort_type: "comment", want_code: "2"},
		{sort_type: "tabom", want_code: "2"},
		{sort_type: "excellent", want_code: "2"},
	} {
		t.Run(test_case.sort_type, func(t *testing.T) {
			result := route.Api_bbs(config, board_id, "1", test_case.sort_type)
			if result["response"] != "ok" {
				t.Fatalf("Api_bbs response = %v, want ok", result["response"])
			}
			data, ok := result["data"].([]map[string]string)
			if !ok || len(data) != len(posts) {
				t.Fatalf("Api_bbs data = %T(%v), want %d posts", result["data"], result["data"], len(posts))
			}
			if data[0]["set_code"] != test_case.want_code {
				t.Errorf("first %s-sorted post = %q, want %q", test_case.sort_type, data[0]["set_code"], test_case.want_code)
			}
			if data[0]["user_id"] != "admin_ip" {
				t.Errorf("post user_id = %q, want admin_ip", data[0]["user_id"])
			}
			if _, exists := data[0]["user_id_render"]; !exists {
				t.Errorf("post user_id_render missing from API response")
			}
			if data[0]["tags"] != "phase5" {
				t.Errorf("post tags = %q, want phase5", data[0]["tags"])
			}
			if result["has_next"] != false {
				t.Errorf("has_next = %v, want false", result["has_next"])
			}
		})
	}

	filtered := route.Api_bbs_filter(config, board_id, "comment/2/1")
	if filtered["response"] != "ok" {
		t.Fatalf("Api_bbs_filter response = %v, want ok", filtered["response"])
	}
	filtered_posts, ok := filtered["data"].([]map[string]string)
	if !ok || len(filtered_posts) != 2 || filtered_posts[0]["set_code"] != "3" || filtered_posts[1]["set_code"] != "2" {
		t.Fatalf("comment-filtered posts = %T(%v), want posts 3 and 2", filtered["data"], filtered["data"])
	}
}
