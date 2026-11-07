package test

import (
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

func Test_Api_edit_and_revert_share_content_validation(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso"} {
		t.Run(backend, func(t *testing.T) {
			db := audit_setup(t, backend)
			audit_seed_auth(db, t)

			if _, err := db.Exec("insert into other (name, data, coverage) values ('document_content_max_length', '1', '')"); err != nil {
				t.Fatal(err)
			}
			if _, err := db.Exec("insert into data (title, data, type) values ('Example', '[[https://example.com/doc|external example]]', '')"); err != nil {
				t.Fatal(err)
			}
			if _, err := db.Exec("insert into history (id, title, data, date, ip, send, leng, hide, type) values ('1', 'Example', 'too long', '', '', '', '0', '', '')"); err != nil {
				t.Fatal(err)
			}

			config := tool.Config{IP: "admin_ip"}
			responses := []struct {
				name     string
				response map[string]any
			}{
				{
					name:     "edit",
					response: route.Api_edit_post(config, "Example", "too long", "", "", ""),
				},
				{
					name:     "revert",
					response: route.Api_edit_revert_post(config, "Example", "1", "", ""),
				},
			}

			for _, test_case := range responses {
				t.Run(test_case.name, func(t *testing.T) {
					if test_case.response["response"] != "error" || test_case.response["data"] != "overflow max length" {
						t.Fatalf("response=%v, want overflow max length", test_case.response)
					}
				})
			}

			if _, err := db.Exec("insert into other (name, data, coverage) values ('copyright_checkbox_text', 'required', '')"); err != nil {
				t.Fatal(err)
			}
			for _, test_case := range []struct {
				name     string
				response map[string]any
			}{
				{name: "edit checkbox precedence", response: route.Api_edit_post(config, "Example", "too long", "", "", "")},
				{name: "revert checkbox precedence", response: route.Api_edit_revert_post(config, "Example", "1", "", "")},
			} {
				t.Run(test_case.name, func(t *testing.T) {
					if test_case.response["response"] != "error" || test_case.response["data"] != "checkbox check require" {
						t.Fatalf("response=%v, want checkbox validation before max length", test_case.response)
					}
				})
			}

			var data string
			if err := db.QueryRow("select data from data where title = 'Example'").Scan(&data); err != nil {
				t.Fatal(err)
			}
			if data != "[[https://example.com/doc|external example]]" {
				t.Fatalf("rejected edit/revert changed document content: %q", data)
			}
			var history_count int
			if err := db.QueryRow("select count(*) from history where title = 'Example'").Scan(&history_count); err != nil {
				t.Fatal(err)
			}
			if history_count != 1 {
				t.Fatalf("rejected edit/revert wrote history rows: %d, want original row only", history_count)
			}
		})
	}
}
