package test

import (
	"os"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

func Test_Api_bbs_search_fallback_preserves_board_acl_filter(t *testing.T) {
	backends := []string{"sqlite", "turso"}
	if os.Getenv("MYSQL_AUDIT_DSN") != "" {
		backends = append(backends, "mysql")
	}

	for _, backend := range backends {
		t.Run(backend, func(t *testing.T) {
			db := new_record_db(t, backend)
			for _, row := range []struct {
				name string
				acl  string
			}{
				{name: "search_reader", acl: "view"},
				{name: "search_reader", acl: "bbs_main_view"},
			} {
				if _, err := db.Exec(`insert into alist (name, acl) values (?, ?)`, row.name, row.acl); err != nil {
					t.Fatalf("seed auth permission: %v", err)
				}
			}
			if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl', 'search_user', 'search_reader')`); err != nil {
				t.Fatalf("seed search user: %v", err)
			}
			for _, board := range []struct {
				id  string
				acl string
			}{
				{id: "2", acl: "search_reader"},
				{id: "3", acl: "owner"},
			} {
				if _, err := db.Exec(`insert into bbs_set (set_name, set_code, set_id, set_data) values ('bbs_name', '', ?, ?)`, board.id, "Search board "+board.id); err != nil {
					t.Fatalf("seed board %s: %v", board.id, err)
				}
				if _, err := db.Exec(`insert into bbs_set (set_name, set_code, set_id, set_data) values ('bbs_view_acl', '', ?, ?)`, board.id, board.acl); err != nil {
					t.Fatalf("seed board ACL %s: %v", board.id, err)
				}
				for _, post_data := range []struct {
					name string
					data string
				}{
					{name: "title", data: "needle topic " + board.id},
					{name: "date", data: "2026-01-01 00:00:00"},
					{name: "user_id", data: "post_author"},
				} {
					if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, '1', ?, ?)`, post_data.name, board.id, post_data.data); err != nil {
						t.Fatalf("seed post on board %s: %v", board.id, err)
					}
				}
			}

			config := tool.Config{IP: "search_user", Session: &emptySession{}}
			if !tool.Check_permission(db, "bbs_main_view", config.IP) {
				t.Fatal("fixture error: search user lacks the BBS main-view permission")
			}
			if !tool.Check_acl(db, "2", "", "bbs_view", config.IP) {
				t.Fatal("fixture error: search user cannot view the visible board")
			}
			if tool.Check_acl(db, "3", "", "bbs_view", config.IP) {
				t.Fatal("fixture error: search user can view the private board")
			}
			if _, index_ready := tool.Search_bbs_index_search("needle", "", 0, 50); index_ready {
				t.Skip("BBS search index is ready; this regression targets the SQL fallback")
			}

			result := route.Api_bbs_search(config, "needle", "", "1")
			if result["response"] != "ok" {
				t.Fatalf("BBS search response = %#v, want ok", result)
			}
			data_list, ok := result["data"].([]map[string]string)
			if !ok {
				t.Fatalf("BBS search data type = %T, want []map[string]string", result["data"])
			}
			if len(data_list) != 1 || data_list[0]["set_id"] != "2" || data_list[0]["title"] != "needle topic 2" {
				t.Fatalf("BBS search results = %#v, want only the visible board's post", data_list)
			}
		})
	}
}
