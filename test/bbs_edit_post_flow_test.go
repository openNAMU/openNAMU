package test

import (
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

func Test_bbs_edit_post_workflow_contract(t *testing.T) {
	db := apiSlowDB(t)
	defer tool.DB_close(db)

	const board_id = "72"
	const post_code = "1"
	bbs_comment_seed_board(db, t, board_id, post_code, "original title", "admin_ip")
	if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values ('data', ?, ?, ?)`, post_code, board_id, "original post"); err != nil {
		t.Fatalf("seed post body: %v", err)
	}
	if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values ('comment', '1', ?, ?)`, board_id+"-"+post_code, "original comment"); err != nil {
		t.Fatalf("seed comment: %v", err)
	}
	if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values ('comment_user_id', '1', ?, ?)`, board_id+"-"+post_code, "admin_ip"); err != nil {
		t.Fatalf("seed comment author: %v", err)
	}

	config := tool.Config{IP: "admin_ip", Session: &emptySession{}}
	post_edit := route.Api_bbs_w_edit_post(config, board_id, post_code, "", "updated title", "updated post", "", "", "")
	if post_edit["response"] != "ok" || post_edit["data"] != post_code {
		t.Fatalf("existing post edit response = %#v, want ok with post code %q", post_edit, post_code)
	}
	if title, _ := tool.Get_bbs_data_value(db, board_id, post_code, "title"); title != "updated title" {
		t.Fatalf("edited post title = %q, want updated title", title)
	}
	if body, _ := tool.Get_bbs_data_value(db, board_id, post_code, "data"); body != "updated post" {
		t.Fatalf("edited post body = %q, want updated post", body)
	}

	comment_edit := route.Api_bbs_w_edit_post(config, board_id, post_code, "1", "", "updated comment", "", "", "")
	if comment_edit["response"] != "ok" || comment_edit["data"] != post_code {
		t.Fatalf("comment edit response = %#v, want ok with parent post code %q", comment_edit, post_code)
	}
	if comment, _ := tool.Get_bbs_data_value(db, board_id+"-"+post_code, "1", "comment"); comment != "updated comment" {
		t.Fatalf("edited comment = %q, want updated comment", comment)
	}

	new_post := route.Api_bbs_w_edit_post(config, board_id, "", "", "new title", "new body", "", "", "")
	if new_post["response"] != "ok" || new_post["data"] != "2" {
		t.Fatalf("new post response = %#v, want ok with post code 2", new_post)
	}
	if body, _ := tool.Get_bbs_data_value(db, board_id, "2", "data"); body != "new body" {
		t.Fatalf("new post body = %q, want new body", body)
	}
}
