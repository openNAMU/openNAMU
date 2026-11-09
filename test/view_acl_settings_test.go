package test

import (
	"strings"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

func Test_view_acl_document_settings_and_owner_controls(t *testing.T) {
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
			exec(`insert into data (title, data, type) values ('AclFormDoc', 'document content', '')`)
			for set_name, set_data := range map[string]string{
				"document_top":               "<b>TOP_SETTING_MARKER</b>",
				"document_top_markup":        "plain",
				"document_editor_top":        "<i>EDITOR_TOP_SETTING_MARKER</i>",
				"document_editor_top_markup": "plain",
			} {
				exec(`insert into data_set (doc_name, doc_rev, set_name, set_data) values ('AclFormDoc', '', ?, ?)`, set_name, set_data)
			}

			manager_config := tool.Config{IP: "acl_manager_ip", Session: &emptySession{}}
			manager_form := route.View_acl(manager_config, "AclFormDoc", false, nil)
			for _, expected := range []string{
				`name="document_top" disabled`,
				`name="document_editor_top" disabled`,
				`&lt;b&gt;TOP_SETTING_MARKER&lt;/b&gt;`,
				`&lt;i&gt;EDITOR_TOP_SETTING_MARKER&lt;/i&gt;`,
				`value="plain" selected`,
			} {
				if !strings.Contains(manager_form, expected) {
					t.Errorf("non-owner ACL form does not contain %q", expected)
				}
			}

			multiple_form := route.View_acl(manager_config, "", true, nil)
			if !strings.Contains(multiple_form, `name="title_name"`) {
				t.Fatal("multiple-document ACL form lost its title list input")
			}
			if strings.Contains(multiple_form, "TOP_SETTING_MARKER") || strings.Contains(multiple_form, "EDITOR_TOP_SETTING_MARKER") {
				t.Fatal("multiple-document ACL form unexpectedly contains a single-document setting")
			}

			owner_config := tool.Config{IP: "admin_ip", Session: &emptySession{}}
			owner_form := route.View_acl(owner_config, "AclFormDoc", false, nil)
			if strings.Contains(owner_form, `name="document_top" disabled`) || strings.Contains(owner_form, `name="document_editor_top" disabled`) {
				t.Fatal("owner ACL form unexpectedly disables owner-only settings")
			}
		})
	}
}
