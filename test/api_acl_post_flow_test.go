package test

import (
	"net/url"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

func Test_api_acl_post_preserves_owner_document_settings(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso"} {
		t.Run(backend, func(t *testing.T) {
			db := audit_setup(t, backend)
			audit_seed_auth(db, t)

			exec := func(query string, args ...any) {
				t.Helper()
				if _, err := db.Exec(query, args...); err != nil {
					t.Fatalf("seed %q: %v", query, err)
				}
			}
			exec(`insert into alist (name, acl) values ('acl_manager_group', 'document_acl_manage')`)
			exec(`insert into user_set (name, id, data) values ('acl', 'acl_manager_ip', 'acl_manager_group')`)
			exec(`insert into data (title, data, type) values ('AclPostTarget', 'document content', '')`)

			protected_settings := map[string]string{
				"document_markup":            "namumark_beta",
				"document_top":               "preserved top",
				"document_top_markup":        "plain",
				"document_editor_top":        "preserved editor top",
				"document_editor_top_markup": "plain",
			}
			for set_name, set_data := range protected_settings {
				exec(`insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, '', ?, ?)`, "AclPostTarget", set_name, set_data)
			}

			config := tool.Config{IP: "acl_manager_ip", Session: &emptySession{}}
			if !tool.Check_permission(db, "document_acl_manage", config.IP) {
				t.Fatal("fixture error: ACL manager permission was not granted")
			}
			if tool.Check_permission(db, "owner", config.IP) {
				t.Fatal("fixture error: ACL manager unexpectedly has owner permission")
			}

			result := route.Api_acl_post(config, "AclPostTarget", false, url.Values{
				"document_markup":            {"namumark_beta"},
				"document_top":               {"untrusted replacement"},
				"document_top_markup":        {"html"},
				"document_editor_top":        {"untrusted editor replacement"},
				"document_editor_top_markup": {"html"},
				"why":                        {"ACL update"},
			})
			if response, _ := result["response"].(string); response != "ok" {
				t.Fatalf("ACL save response = %#v, want ok", result)
			}

			for set_name, want := range protected_settings {
				got := tool.Get_document_setting_value(db, "AclPostTarget", set_name, "")
				if got != want {
					t.Errorf("protected %s = %q, want preserved value %q", set_name, got, want)
				}
			}
		})
	}
}
