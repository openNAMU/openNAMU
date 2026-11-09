package test

import (
	"strconv"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

func Test_Api_bbs_w_comment_page_preserves_offset(t *testing.T) {
	db := apiSlowDB(t)
	defer tool.DB_close(db)

	const board_id = "43"
	bbs_comment_seed_board(db, t, board_id, "1", "topic", "admin_ip")
	for comment_number := 1; comment_number <= 52; comment_number++ {
		bbs_comment_seed_parent(db, t, board_id, "1", strconv.Itoa(comment_number), "admin_ip")
	}

	config := tool.Config{IP: "admin_ip", Session: &emptySession{}}
	page := route.Api_bbs_w_comment_page(config, board_id+"-1", true, "around", 2)
	comments, ok := page["data"].([]map[string]string)
	if page["response"] != "ok" || !ok || len(comments) != 2 {
		t.Fatalf("page 2 = response %v, data %T(%v); want ok and 2 comments", page["response"], page["data"], page["data"])
	}
	if comments[0]["code"] != "51" || comments[1]["code"] != "52" {
		t.Fatalf("page 2 comment codes = %q, %q; want 51, 52", comments[0]["code"], comments[1]["code"])
	}
}

func Test_Api_bbs_w_comment_traversal_preserves_nested_branch_order(t *testing.T) {
	db := apiSlowDB(t)
	defer tool.DB_close(db)

	const board_id = "44"
	bbs_comment_seed_board(db, t, board_id, "1", "topic", "admin_ip")
	for _, comment := range []struct {
		parent_id   string
		parent_code string
		code        string
	}{
		{parent_id: board_id, parent_code: "1", code: "1"},
		{parent_id: board_id, parent_code: "1", code: "2"},
		{parent_id: board_id + "-1", parent_code: "1", code: "1"},
		{parent_id: board_id + "-1", parent_code: "1", code: "2"},
		{parent_id: board_id + "-1-1", parent_code: "1", code: "1"},
		{parent_id: board_id + "-1-1", parent_code: "2", code: "1"},
		{parent_id: board_id + "-1", parent_code: "2", code: "1"},
		{parent_id: board_id + "-1", parent_code: "2", code: "2"},
		{parent_id: board_id + "-1-2", parent_code: "2", code: "1"},
	} {
		bbs_comment_seed_parent(db, t, comment.parent_id, comment.parent_code, comment.code, "admin_ip")
	}
	if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values ('blind', '1', ?, 'O')`, board_id+"-1-1-1"); err != nil {
		t.Fatalf("seed blind nested comment: %v", err)
	}

	config := tool.Config{IP: "unprivileged_ip", Session: &emptySession{}}
	result := route.Api_bbs_w_comment_all(config, board_id+"-1", true, "around")
	comments, ok := result["data"].([]map[string]string)
	if result["response"] != "ok" || !ok || len(comments) != 9 {
		t.Fatalf("nested comments = response %v, data %T(%v); want ok and 9 comments", result["response"], result["data"], result["data"])
	}
	want := []struct {
		id   string
		code string
	}{
		{id: board_id + "-1", code: "1"},
		{id: board_id + "-1-1", code: "1"},
		{id: board_id + "-1-1-1", code: "1"},
		{id: board_id + "-1-1", code: "2"},
		{id: board_id + "-1-1-2", code: "1"},
		{id: board_id + "-1", code: "2"},
		{id: board_id + "-1-2", code: "1"},
		{id: board_id + "-1-2", code: "2"},
		{id: board_id + "-1-2-2", code: "1"},
	}
	for index, expected := range want {
		if comments[index]["id"] != expected.id || comments[index]["code"] != expected.code {
			t.Errorf("comment %d = id %q code %q, want id %q code %q", index, comments[index]["id"], comments[index]["code"], expected.id, expected.code)
		}
	}
	if comments[2]["comment"] != "" {
		t.Errorf("blind nested comment content = %q, want empty for non-manager", comments[2]["comment"])
	}
}

func Test_Api_bbs_w_comment_traversal_is_depth_first(t *testing.T) {
	db := apiSlowDB(t)
	defer tool.DB_close(db)

	const board_id = "42"
	bbs_comment_seed_board(db, t, board_id, "1", "topic", "admin_ip")
	bbs_comment_seed_parent(db, t, board_id, "1", "1", "admin_ip")
	bbs_comment_seed_parent(db, t, board_id+"-1", "1", "1", "admin_ip")
	bbs_comment_seed_parent(db, t, board_id, "1", "2", "admin_ip")

	config := tool.Config{IP: "admin_ip", Session: &emptySession{}}
	one := route.Api_bbs_w_comment_one(config, true, "", board_id+"-1-1")
	one_comments, ok := one["data"].([]map[string]string)
	if one["response"] != "ok" || !ok || len(one_comments) != 1 || one_comments[0]["code"] != "1" {
		t.Fatalf("single-comment API = response %v, data %T(%v); want one comment with code 1", one["response"], one["data"], one["data"])
	}

	result := route.Api_bbs_w_comment_all(config, board_id+"-1", true, "around")
	if result["response"] != "ok" {
		t.Fatalf("Api_bbs_w_comment_all response = %v, want ok", result["response"])
	}
	comments, ok := result["data"].([]map[string]string)
	if !ok || len(comments) != 3 {
		t.Fatalf("Api_bbs_w_comment_all data = %T(%v), want 3 comments", result["data"], result["data"])
	}
	want := []struct {
		id   string
		code string
	}{
		{id: board_id + "-1", code: "1"},
		{id: board_id + "-1-1", code: "1"},
		{id: board_id + "-1", code: "2"},
	}
	for index, expected := range want {
		if comments[index]["id"] != expected.id || comments[index]["code"] != expected.code {
			t.Errorf("comment %d = id %q code %q, want id %q code %q", index, comments[index]["id"], comments[index]["code"], expected.id, expected.code)
		}
	}

	page := route.Api_bbs_w_comment_page(config, board_id+"-1", true, "around", 1)
	page_comments, ok := page["data"].([]map[string]string)
	if page["response"] != "ok" || !ok || len(page_comments) != len(want) {
		t.Fatalf("Api_bbs_w_comment_page = response %v, data %T(%v); want ok and 3 comments", page["response"], page["data"], page["data"])
	}
	for index, expected := range want {
		if page_comments[index]["id"] != expected.id || page_comments[index]["code"] != expected.code {
			t.Errorf("page comment %d = id %q code %q, want id %q code %q", index, page_comments[index]["id"], page_comments[index]["code"], expected.id, expected.code)
		}
	}

	if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values ('pinned', '1', ?, '1')`, board_id+"-1"); err != nil {
		t.Fatalf("seed pinned comment: %v", err)
	}
	pinned_comments := route.Api_bbs_w_comment_pinned_data(config, board_id, "1")
	if len(pinned_comments) != 1 || pinned_comments[0]["code"] != "1" {
		t.Fatalf("pinned comments = %v, want comment code 1", pinned_comments)
	}

	if _, err := db.Exec(`insert into bbs_set (set_name, set_code, set_id, set_data) values ('bbs_view_acl', '', ?, 'owner')`, board_id); err != nil {
		t.Fatalf("seed board view ACL: %v", err)
	}
	denied_config := tool.Config{IP: "unprivileged_ip", Session: &emptySession{}}
	denied_one := route.Api_bbs_w_comment_one(denied_config, false, "around", board_id+"-1")
	if denied_one["response"] != "require auth" {
		t.Errorf("unauthorized comment response = %v, want require auth", denied_one["response"])
	}
	denied_all := route.Api_bbs_w_comment_all(denied_config, board_id+"-1", false, "around")
	denied_comments, ok := denied_all["data"].([]map[string]string)
	if denied_all["response"] != "ok" || !ok || len(denied_comments) != 0 {
		t.Errorf("unauthorized all-comments response = %v, data %T(%v); want legacy ok with an empty slice", denied_all["response"], denied_all["data"], denied_all["data"])
	}
	denied_page := route.Api_bbs_w_comment_page(denied_config, board_id+"-1", false, "around", 1)
	denied_page_comments, ok := denied_page["data"].([]map[string]string)
	if denied_page["response"] != "ok" || !ok || len(denied_page_comments) != 0 {
		t.Errorf("unauthorized comment page response = %v, data %T(%v); want ok with an empty slice", denied_page["response"], denied_page["data"], denied_page["data"])
	}
}
