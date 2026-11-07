package test

import (
	"reflect"
	"strings"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

func Test_Api_bbs_search_adapters_match_shared_operation(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso"} {
		t.Run(backend, func(t *testing.T) {
			db := audit_setup(t, backend)
			audit_seed_auth(db, t)

			if _, err := db.Exec("insert into bbs_set (set_name, set_code, set_id, set_data) values ('bbs_name', '', '2', 'Search board')"); err != nil {
				t.Fatal(err)
			}
			for _, row := range []struct {
				name string
				data string
			}{
				{name: "title", data: "needle topic"},
				{name: "data", data: "needle body"},
				{name: "date", data: "2026-11-06 00:00:00"},
				{name: "user_id", data: "192.0.2.44"},
			} {
				if _, err := db.Exec("insert into bbs_data (set_name, set_id, set_code, set_data) values (?, '2', '1', ?)", row.name, row.data); err != nil {
					t.Fatal(err)
				}
			}

			config := tool.Config{IP: "admin_ip"}
			for _, test_case := range []struct {
				name        string
				adapter     func() map[string]any
				search_type string
			}{
				{
					name:        "title",
					adapter:     func() map[string]any { return route.Api_bbs_search(config, "needle", "2", "1") },
					search_type: "title",
				},
				{
					name:        "data",
					adapter:     func() map[string]any { return route.Api_bbs_search_data(config, "needle", "2", "1") },
					search_type: "data",
				},
			} {
				t.Run(test_case.name+" API", func(t *testing.T) {
					adapter_result := test_case.adapter()
					shared_result := route.Api_bbs_search_internal(config, "needle", "2", "1", test_case.search_type)
					if !reflect.DeepEqual(adapter_result, shared_result) {
						t.Fatalf("adapter result=%v, shared result=%v", adapter_result, shared_result)
					}
					if adapter_result["response"] != "ok" {
						t.Fatalf("search response=%v, want ok", adapter_result)
					}
				})
			}

			view_config := tool.Config{IP: "admin_ip", Session: &emptySession{}}
			for _, test_case := range []struct {
				name        string
				view        func() string
				want_result string
			}{
				{
					name:        "title search",
					view:        func() string { return route.View_bbs_search(view_config, "2", "needle", "1") },
					want_result: "topic",
				},
				{
					name:        "body search",
					view:        func() string { return route.View_bbs_search_data(view_config, "2", "needle", "1") },
					want_result: "body",
				},
			} {
				t.Run(test_case.name+" View", func(t *testing.T) {
					output := test_case.view()
					if !strings.Contains(output, "Search board") || !strings.Contains(strings.ToLower(output), test_case.want_result) {
						t.Fatalf("search View output does not contain board and expected result %q: %q", test_case.want_result, output)
					}
				})
			}
		})
	}
}
