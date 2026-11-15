package test

import (
	"database/sql"
	"net/url"
	"strings"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
	"opennamu/route/tool/markup"
)

const markup_xss_payload = `<script>window.__retired_markup_xss__ = true;</script><b>SAFE_MARKUP_BODY</b>`

func Markup_xss_seed_document(db *sql.DB, t *testing.T, title string, data string, document_markup string) {
	t.Helper()
	if _, err := db.Exec(`insert into data (title, data, type) values (?, ?, '')`, title, data); err != nil {
		t.Fatalf("seed document %q: %v", title, err)
	}
	if document_markup != "" {
		if _, err := db.Exec(`insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, '', 'document_markup', ?)`, title, document_markup); err != nil {
			t.Fatalf("seed document markup for %q: %v", title, err)
		}
	}
}

func Markup_xss_response(t *testing.T, result map[string]any, want string) {
	t.Helper()
	response, _ := result["response"].(string)
	if response != want {
		t.Fatalf("API response = %#v, want %q", result, want)
	}
}

func Markup_xss_assert_escaped(t *testing.T, context string, rendered string) {
	t.Helper()
	if strings.Contains(rendered, `<script>window.__retired_markup_xss__`) {
		t.Fatalf("%s emitted executable retired markup: %s", context, rendered)
	}
	if !strings.Contains(rendered, `&lt;script&gt;window.__retired_markup_xss__`) {
		t.Fatalf("%s did not render script source as escaped text: %s", context, rendered)
	}
}

func Test_markup_retirement_legacy_rendering_and_surfaces(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso"} {
		t.Run(backend, func(t *testing.T) {
			db := audit_setup(t, backend)
			audit_seed_auth(db, t)
			seed := func(query string, args ...any) {
				t.Helper()
				if _, err := db.Exec(query, args...); err != nil {
					t.Fatalf("seed %q: %v", query, err)
				}
			}
			seed(`insert into other (name, data, coverage) values ('markup', 'html', '')`)
			Markup_xss_seed_document(db, t, "OrdinaryLegacyHtml", markup_xss_payload, "html")
			Markup_xss_seed_document(db, t, "user:alice", markup_xss_payload, "custom")
			Markup_xss_seed_document(db, t, "InheritedLegacyMarkup", markup_xss_payload, "")
			seed(`insert into data_set (doc_name, doc_rev, set_name, set_data) values ('OrdinaryLegacyHtml', '', 'document_top', ?)`, markup_xss_payload)
			seed(`insert into data_set (doc_name, doc_rev, set_name, set_data) values ('OrdinaryLegacyHtml', '', 'document_top_markup', 'html')`)
			seed(`insert into data_set (doc_name, doc_rev, set_name, set_data) values ('OrdinaryLegacyHtml', '', 'document_editor_top', ?)`, markup_xss_payload)
			seed(`insert into data_set (doc_name, doc_rev, set_name, set_data) values ('OrdinaryLegacyHtml', '', 'document_editor_top_markup', 'custom')`)

			if got := tool.Get_document_markup(db, "OrdinaryLegacyHtml", "document"); got != "raw" {
				t.Fatalf("legacy html effective markup = %q, want raw", got)
			}
			if got := tool.Get_document_markup(db, "user:alice", "document"); got != "raw" {
				t.Fatalf("legacy custom effective markup = %q, want raw", got)
			}
			if got := tool.Get_document_markup(db, "InheritedLegacyMarkup", "document"); got != "raw" {
				t.Fatalf("inherited global html markup = %q, want raw", got)
			}

			for _, doc_name := range []string{"OrdinaryLegacyHtml", "user:alice", "InheritedLegacyMarkup"} {
				rendered := markup.Get_render_with_requester(db, doc_name, markup_xss_payload, "normal", nil, "admin_ip")["data"]
				Markup_xss_assert_escaped(t, "stored document "+doc_name, rendered)
			}
			for _, legacy_markup := range []string{"html", "custom"} {
				rendered := markup.Get_render_direct(db, "DirectLegacyMarkup", markup_xss_payload, legacy_markup, "test", "view")["data"]
				Markup_xss_assert_escaped(t, "direct renderer "+legacy_markup, rendered)
			}
			for _, safe_markup := range []string{"raw", "plain"} {
				rendered := markup.Get_render_direct(db, "SafeMarkup", markup_xss_payload, safe_markup, "test", "view")["data"]
				Markup_xss_assert_escaped(t, "safe renderer "+safe_markup, rendered)
			}
			for _, legacy_markup := range []string{"html", "custom"} {
				if got := route.Setting_markup_normalize(legacy_markup); got != "raw" {
					t.Errorf("legacy setting markup %q normalized to %q, want raw", legacy_markup, got)
				}
				Markup_xss_assert_escaped(t, "top markup helper "+legacy_markup, route.Setting_render_markup(db, markup_xss_payload, legacy_markup))
			}
			Markup_xss_assert_escaped(t, "editor top markup", route.Document_editor_top_render(db, "OrdinaryLegacyHtml"))
			// Test the document-top setting separately after removing editor-top content.
			seed(`delete from data_set where doc_name = 'OrdinaryLegacyHtml' and set_name = 'document_editor_top'`)
			Markup_xss_assert_escaped(t, "document top markup", route.Document_editor_top_render(db, "OrdinaryLegacyHtml"))

			markup_options := markup.List_markup()
			setting_markup_options := route.Setting_markup_options()
			for _, retired := range []string{"html", "custom"} {
				if tool.Arr_in_str(markup_options, retired) || tool.Arr_in_str(setting_markup_options, retired) {
					t.Errorf("retired markup %q remains selectable: markup=%v setting=%v", retired, markup_options, setting_markup_options)
				}
			}
			for _, supported := range []string{"raw", "plain", "namumark", "namumark_beta", "markdown", "macromark"} {
				if !tool.Arr_in_str(markup_options, supported) || !tool.Arr_in_str(setting_markup_options, supported) {
					t.Errorf("supported markup %q missing from options: markup=%v setting=%v", supported, markup_options, setting_markup_options)
				}
			}
			init_markup_options, ok := tool.Get_init_set_list("markup")["markup"]["list"].([]string)
			if !ok || tool.Arr_in_str(init_markup_options, "html") || tool.Arr_in_str(init_markup_options, "custom") || !tool.Arr_in_str(init_markup_options, "raw") {
				t.Errorf("global markup options contain retired values or omit raw: %#v", init_markup_options)
			}
			acl_ui := tool.Get_markup_select_ui(db, tool.Config{IP: "admin_ip", Session: &emptySession{}}, "OrdinaryLegacyHtml", "raw", "", "")
			if !strings.Contains(acl_ui, `<option value="raw" selected>raw</option>`) || strings.Contains(acl_ui, `<option value="html"`) || strings.Contains(acl_ui, `<option value="custom"`) {
				t.Errorf("ACL markup selector did not offer/select the safe raw replacement: %s", acl_ui)
			}

			for _, legacy_markup := range []string{"html", "custom"} {
				if _, err := db.Exec(`update other set data = ? where name = 'markup'`, legacy_markup); err != nil {
					t.Fatalf("set legacy global markup %q: %v", legacy_markup, err)
				}
				if got := tool.Get_document_markup(db, "InheritedLegacyMarkup", "document"); got != "raw" {
					t.Errorf("legacy global %q resolves to %q, want raw", legacy_markup, got)
				}
				rendered := markup.Get_render_with_requester(db, "InheritedLegacyMarkup", markup_xss_payload, "normal", nil, "admin_ip")["data"]
				Markup_xss_assert_escaped(t, "inherited global "+legacy_markup, rendered)
			}

			admin_config := tool.Config{IP: "admin_ip", Session: &emptySession{}}
			for set_id, legacy_markup := range map[string]string{"LegacyHtmlBBS": "html", "LegacyCustomBBS": "custom"} {
				seed(`insert into bbs_set (set_name, set_code, set_id, set_data) values ('bbs_markup', '', ?, ?)`, set_id, legacy_markup)
				Markup_xss_assert_escaped(t, "legacy BBS markup "+legacy_markup, route.Get_bbs_render(db, set_id, markup_xss_payload, "view", admin_config))
			}
			include_payload := tool.Base64_encode(`{"version":1,"name":"OrdinaryLegacyHtml"}`)
			Markup_xss_assert_escaped(t, "standalone include", route.View_include(admin_config, include_payload))
			tool.Do_add_history(db, "OrdinaryLegacyHtml", markup_xss_payload, "2024-01-01 00:00:00", "admin_ip", "", "0", "add", "")
			history_revision := ""
			if err := db.QueryRow(`select id from history where title = 'OrdinaryLegacyHtml'`).Scan(&history_revision); err != nil {
				t.Fatalf("read seeded history revision: %v", err)
			}
			Markup_xss_assert_escaped(t, "historical revision", route.View_render(admin_config, "OrdinaryLegacyHtml", history_revision))
			Markup_xss_assert_escaped(t, "edit preview", route.View_edit_preview(admin_config, "OrdinaryLegacyHtml", markup_xss_payload, "normal", ""))

			for doc_name, expected_markup := range map[string]string{"OrdinaryLegacyHtml": "html", "user:alice": "custom"} {
				stored_markup := ""
				if err := db.QueryRow(`select set_data from data_set where doc_name = ? and set_name = 'document_markup'`, doc_name).Scan(&stored_markup); err != nil || stored_markup != expected_markup {
					t.Errorf("legacy markup setting for %q was migrated: got %q, want stored %q (%v)", doc_name, stored_markup, expected_markup, err)
				}
			}
			stored_global_markup := ""
			if err := db.QueryRow(`select data from other where name = 'markup'`).Scan(&stored_global_markup); err != nil || stored_global_markup != "custom" {
				t.Errorf("legacy global markup setting was migrated: got %q, want stored custom (%v)", stored_global_markup, err)
			}

			for _, doc_name := range []string{"OrdinaryLegacyHtml", "user:alice", "InheritedLegacyMarkup"} {
				var stored_body string
				if err := db.QueryRow(`select data from data where title = ?`, doc_name).Scan(&stored_body); err != nil {
					t.Fatalf("read preserved body for %q: %v", doc_name, err)
				}
				if stored_body != markup_xss_payload {
					t.Errorf("stored body for %q changed: %q", doc_name, stored_body)
				}
			}
			seed(`insert into bbs_set (set_name, set_code, set_id, set_data) values ('bbs_markup', '', 'LegacyNormalBBS', 'normal')`)
			normal_bbs := route.Get_bbs_render(db, "LegacyNormalBBS", `'''NORMAL_BBS_FORMAT'''`, "view", admin_config)
			if !strings.Contains(normal_bbs, "<b>NORMAL_BBS_FORMAT</b>") {
				t.Errorf("BBS normal markup sentinel stopped using the default Namumark renderer: %s", normal_bbs)
			}
		})
	}
}

func Test_markup_retirement_rejects_new_writes_without_partial_updates(t *testing.T) {
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
			exec(`insert into data (title, data, type) values ('MarkupWriteTarget', 'BODY_UNCHANGED', '')`)
			exec(`insert into data_set (doc_name, doc_rev, set_name, set_data) values ('MarkupWriteTarget', '', 'document_markup', 'namumark')`)
			exec(`insert into data_set (doc_name, doc_rev, set_name, set_data) values ('MarkupWriteTarget', '', 'document_top_markup', 'plain')`)
			exec(`insert into data_set (doc_name, doc_rev, set_name, set_data) values ('MarkupWriteTarget', '', 'document_editor_top_markup', 'plain')`)
			exec(`insert into acl (title, data, type) values ('MarkupWriteTarget', 'prior reason', 'why')`)

			admin_config := tool.Config{IP: "admin_ip", Session: &emptySession{}}
			for _, retired := range []string{"html", "custom"} {
				for set_name, expected_markup := range map[string]string{
					"document_markup":            "namumark",
					"document_top_markup":        "plain",
					"document_editor_top_markup": "plain",
				} {
					Markup_xss_response(t, route.Api_w_set_put(admin_config, "MarkupWriteTarget", set_name, retired), "error")
					got := tool.Get_document_setting_value(db, "MarkupWriteTarget", set_name, "")
					if got != expected_markup {
						t.Errorf("rejected Api_w_set_put changed %s to %q, want %q", set_name, got, expected_markup)
					}
				}
				result := route.Api_acl_post(admin_config, "MarkupWriteTarget", false, url.Values{
					"document_markup": {retired},
					"why":             {"should not persist"},
				})
				Markup_xss_response(t, result, "error")
				if got := tool.Get_document_setting_value(db, "MarkupWriteTarget", "document_markup", ""); got != "namumark" {
					t.Errorf("rejected ACL markup %q changed setting to %q", retired, got)
				}
				if got := route.Acl_value(db, "MarkupWriteTarget", "why"); got != "prior reason" {
					t.Errorf("rejected ACL markup %q partially changed why to %q", retired, got)
				}
				top_result := route.Api_acl_post(admin_config, "MarkupWriteTarget", false, url.Values{
					"document_markup":     {"plain"},
					"document_top_markup": {retired},
					"why":                 {"should not partially update"},
				})
				Markup_xss_response(t, top_result, "error")
				if got := tool.Get_document_setting_value(db, "MarkupWriteTarget", "document_top_markup", ""); got != "plain" {
					t.Errorf("rejected top markup %q changed setting to %q", retired, got)
				}
				if got := tool.Get_document_setting_value(db, "MarkupWriteTarget", "document_markup", ""); got != "namumark" {
					t.Errorf("rejected top markup %q partially changed document markup to %q", retired, got)
				}
				if got := route.Acl_value(db, "MarkupWriteTarget", "why"); got != "prior reason" {
					t.Errorf("rejected top markup %q partially changed why to %q", retired, got)
				}
			}

			Markup_xss_response(t, route.Api_w_set_put(admin_config, "MarkupWriteTarget", "document_markup", "raw"), "ok")
			if got := tool.Get_document_setting_value(db, "MarkupWriteTarget", "document_markup", ""); got != "plain" {
				t.Errorf("raw document markup stored as %q, want canonical plain", got)
			}
			Markup_xss_response(t, route.Api_w_set_put(admin_config, "MarkupWriteTarget", "document_top_markup", "plain"), "ok")
			Markup_xss_response(t, route.Api_acl_post(admin_config, "MarkupWriteTarget", false, url.Values{
				"document_markup": {"plain"},
				"why":             {"safe update"},
			}), "ok")

			seed_setting_permission := func(permission string, group string, identity string) {
				exec(`insert into alist (name, acl) values (?, ?)`, group, permission)
				exec(`insert into user_set (name, id, data) values ('acl', ?, ?)`, identity, group)
			}
			seed_setting_permission("setting_phrase", "phrase_admin_group", "phrase_admin_ip")
			exec(`insert into other (name, data, coverage) values ('body', 'OLD_BODY_SETTING', '')`)
			exec(`insert into other (name, data, coverage) values ('body_markup', 'plain', '')`)
			exec(`insert into other (name, data, coverage) values ('bottom_body', 'OLD_BOTTOM_BODY_SETTING', '')`)
			exec(`insert into other (name, data, coverage) values ('bottom_body_markup', 'plain', '')`)
			for _, retired := range []string{"html", "custom"} {
				phrase_result := route.Api_setting_phrase_post(tool.Config{IP: "phrase_admin_ip", Session: &emptySession{}}, map[string]string{
					"body":               "NEW_BODY_SETTING",
					"body_markup":        retired,
					"bottom_body":        "NEW_BOTTOM_BODY",
					"bottom_body_markup": "plain",
				})
				Markup_xss_response(t, phrase_result, "error")
				var body_setting string
				if err := db.QueryRow(`select data from other where name = 'body' and coverage = ''`).Scan(&body_setting); err != nil || body_setting != "OLD_BODY_SETTING" {
					t.Errorf("rejected phrase markup %q partially changed body setting: %q, %v", retired, body_setting, err)
				}
				var bottom_body_setting string
				if err := db.QueryRow(`select data from other where name = 'bottom_body' and coverage = ''`).Scan(&bottom_body_setting); err != nil || bottom_body_setting != "OLD_BOTTOM_BODY_SETTING" {
					t.Errorf("rejected phrase markup %q partially changed bottom body setting: %q, %v", retired, bottom_body_setting, err)
				}
			}

			seed_setting_permission("bbs_setting", "bbs_manager_group", "bbs_manager_ip")
			exec(`insert into bbs_set (set_name, set_code, set_id, set_data) values ('bbs_name', '', '301', 'Markup test board')`)
			exec(`insert into bbs_set (set_name, set_code, set_id, set_data) values ('bbs_markup', '', '301', 'plain')`)
			for _, retired := range []string{"html", "custom"} {
				result := route.Api_bbs_w_set_put(tool.Config{IP: "bbs_manager_ip", Session: &emptySession{}}, "301", "bbs_markup", retired, "")
				Markup_xss_response(t, result, "error")
				var current_markup string
				if err := db.QueryRow(`select set_data from bbs_set where set_name = 'bbs_markup' and set_id = '301'`).Scan(&current_markup); err != nil || current_markup != "plain" {
					t.Errorf("rejected BBS markup %q changed stored value to %q: %v", retired, current_markup, err)
				}
			}
			Markup_xss_response(t, route.Api_bbs_w_set_put(tool.Config{IP: "bbs_manager_ip", Session: &emptySession{}}, "301", "bbs_markup", "raw", ""), "ok")
			var bbs_markup string
			if err := db.QueryRow(`select set_data from bbs_set where set_name = 'bbs_markup' and set_id = '301'`).Scan(&bbs_markup); err != nil || bbs_markup != "plain" {
				t.Errorf("raw BBS markup stored as %q, want canonical plain: %v", bbs_markup, err)
			}
			Markup_xss_response(t, route.Api_bbs_w_set_put(tool.Config{IP: "bbs_manager_ip", Session: &emptySession{}}, "301", "bbs_markup", "normal", ""), "ok")
			if err := db.QueryRow(`select set_data from bbs_set where set_name = 'bbs_markup' and set_id = '301'`).Scan(&bbs_markup); err != nil || bbs_markup != "normal" {
				t.Errorf("BBS normal sentinel stored as %q, want normal: %v", bbs_markup, err)
			}

			exec(`insert into other (name, data, coverage) values ('markup', 'markdown', '')`)
			for _, retired := range []string{"html", "custom"} {
				if err := route.Emergency_set_other(db, "markup", retired); err == nil {
					t.Errorf("emergency global markup setter accepted %q", retired)
				}
				var global_markup string
				if err := db.QueryRow(`select data from other where name = 'markup'`).Scan(&global_markup); err != nil || global_markup != "markdown" {
					t.Errorf("rejected emergency markup %q changed global setting to %q: %v", retired, global_markup, err)
				}
			}
			if err := route.Emergency_set_other(db, "markup", "raw"); err != nil {
				t.Fatalf("emergency global raw markup rejected: %v", err)
			}
			var global_markup string
			if err := db.QueryRow(`select data from other where name = 'markup'`).Scan(&global_markup); err != nil || global_markup != "plain" {
				t.Errorf("raw global markup stored as %q, want canonical plain: %v", global_markup, err)
			}
		})
	}
}

func Test_markup_link_xss_and_external_target_parsing(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso"} {
		t.Run(backend, func(t *testing.T) {
			db := audit_setup(t, backend)
			audit_seed_auth(db, t)
			for _, test_case := range []struct {
				name   string
				markup string
				data   string
			}{
				{name: "macromark javascript", markup: "macromark", data: `[ex(javascript:alert(1),unsafe)]`},
				{name: "macromark entity scheme", markup: "macromark", data: `[ex(javascript&#58;alert(1),unsafe)]`},
				{name: "macromark vbscript", markup: "macromark", data: `[ex(vbscript:alert(1),unsafe)]`},
				{name: "macromark data scheme", markup: "macromark", data: `[ex(data:text/html%3Bbase64,unsafe)]`},
				{name: "macromark file scheme", markup: "macromark", data: `[ex(file:///etc/passwd,unsafe)]`},
				{name: "macromark malformed URL", markup: "macromark", data: `[ex(https://%,unsafe)]`},
				{name: "macromark control character", markup: "macromark", data: "[ex(java\nscript:alert(1),unsafe)]"},
				{name: "markdown raw macro", markup: "markdown", data: `[ex(javascript:alert(1),unsafe)]`},
				{name: "markdown link syntax", markup: "markdown", data: `[unsafe](javascript:alert(1))`},
			} {
				t.Run(test_case.name, func(t *testing.T) {
					rendered := markup.Get_render_direct(db, "LinkDoc", test_case.data, test_case.markup, "test", "view")["data"]
					lower_rendered := strings.ToLower(rendered)
					if strings.Contains(lower_rendered, `href="javascript:`) || strings.Contains(lower_rendered, `href="vbscript:`) || strings.Contains(lower_rendered, `href="data:`) {
						t.Fatalf("%s emitted executable URL: %s", test_case.name, rendered)
					}
					if strings.Contains(rendered, `<a class="opennamu_link_out" href=`) {
						t.Fatalf("%s retained a clickable external link for an invalid URL: %s", test_case.name, rendered)
					}
					if !strings.Contains(rendered, "unsafe") {
						t.Fatalf("%s discarded the link label: %s", test_case.name, rendered)
					}
				})
			}

			implicit_label := markup.Get_render_direct(db, "LinkDoc", `[ex(https://example.invalid/<script>window.__retired_markup_xss__=true</script>)]`, "macromark", "test", "view")["data"]
			Markup_xss_assert_escaped(t, "Macromark implicit external-link label", implicit_label)
			for _, safe_link := range []string{"ftp://files.example.invalid/path", "tel:+18005550199"} {
				rendered := markup.Get_render_direct(db, "LinkDoc", `[ex(`+safe_link+`,valid)]`, "macromark", "test", "view")["data"]
				if !strings.Contains(rendered, `href="`+safe_link+`"`) {
					t.Errorf("previously safe external scheme %q was not preserved: %s", safe_link, rendered)
				}
			}
			valid_macromark := markup.Get_render_direct(db, "LinkDoc", `[b(working)] [ex(https://example.invalid/?q=target=_self,valid)]`, "macromark", "test", "view")["data"]
			if !strings.Contains(valid_macromark, "<b>working</b>") || !strings.Contains(valid_macromark, `href="https://example.invalid/?q=target=_self"`) {
				t.Fatalf("valid Macromark formatting/link was not preserved: %s", valid_macromark)
			}
			valid_macromark = route.Add_render_external_link_target(valid_macromark)
			if !strings.Contains(valid_macromark, `target="_blank"`) {
				t.Fatalf("quoted URL text target=_self suppressed actual target insertion: %s", valid_macromark)
			}
			actual_target := route.Add_render_external_link_target(`<a class="opennamu_link_out" href="https://example.invalid/" target="_self">link</a>`)
			if strings.Count(actual_target, `target=`) != 1 || !strings.Contains(actual_target, `target="_self"`) {
				t.Fatalf("actual target attribute was duplicated or changed: %s", actual_target)
			}

			valid_markdown := markup.Get_render_direct(db, "LinkDoc", "**bold** and [valid](https://example.invalid/path)", "markdown", "test", "view")["data"]
			if !strings.Contains(valid_markdown, "<b>bold</b>") || !strings.Contains(valid_markdown, `href="https://example.invalid/path"`) {
				t.Fatalf("valid Markdown formatting/link was not preserved: %s", valid_markdown)
			}
			namumark_input := "{{{#!html\n<b>NAMUMARK_SAFE</b><script>window.__retired_markup_xss__=true</script>\n}}}"
			namumark := markup.Get_render_direct(db, "LinkDoc", namumark_input, "namumark", "test", "view")["data"]
			if !strings.Contains(namumark, "<b>NAMUMARK_SAFE</b>") || strings.Contains(namumark, `<script>window.__retired_markup_xss__`) {
				t.Fatalf("existing Namumark safe HTML subset broke: %s", namumark)
			}
			namumark_beta := markup.Get_render_direct(db, "LinkDoc", `<script>window.__retired_markup_xss__ = true;</script>`, "namumark_beta", "test", "view")["data"]
			if strings.Contains(namumark_beta, `<script>window.__retired_markup_xss__`) || !strings.Contains(namumark_beta, `&lt;script&gt;`) {
				t.Fatalf("Namumark beta exposed literal script markup or stopped escaping it: %s", namumark_beta)
			}
		})
	}
}

func Test_acl_bulk_self_page_exception_does_not_authorize_other_documents(t *testing.T) {
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
			exec(`insert into data (title, data, type) values ('user:alice_ip', 'self profile', '')`)
			exec(`insert into data_set (doc_name, doc_rev, set_name, set_data) values ('user:alice_ip', '', 'document_markup', 'plain')`)
			exec(`insert into data (title, data, type) values ('BulkAclTarget', 'target body', '')`)
			exec(`insert into data_set (doc_name, doc_rev, set_name, set_data) values ('BulkAclTarget', '', 'document_markup', 'namumark')`)

			alice_config := tool.Config{IP: "alice_ip", Session: &emptySession{}}
			self_result := route.Api_acl_post(alice_config, "user:alice_ip", false, url.Values{
				"document_markup": {"raw"},
				"why":             {"self page"},
			})
			Markup_xss_response(t, self_result, "ok")
			if got := tool.Get_document_setting_value(db, "user:alice_ip", "document_markup", ""); got != "plain" {
				t.Fatalf("self-page ACL update stored markup %q, want plain", got)
			}

			bulk_result := route.Api_acl_post(alice_config, "user:alice_ip", true, url.Values{
				"title_name":      {"BulkAclTarget"},
				"document_markup": {"raw"},
			})
			Markup_xss_response(t, bulk_result, "require auth")
			if got := tool.Get_document_setting_value(db, "BulkAclTarget", "document_markup", ""); got != "namumark" {
				t.Fatalf("unauthorized bulk ACL changed target markup to %q", got)
			}
			cross_user_result := route.Api_acl_post(alice_config, "BulkAclTarget", false, url.Values{"document_markup": {"plain"}})
			Markup_xss_response(t, cross_user_result, "require auth")

			exec(`insert into alist (name, acl) values ('acl_manager_group', 'document_acl_manage')`)
			exec(`insert into user_set (name, id, data) values ('acl', 'acl_manager_ip', 'acl_manager_group')`)
			manager_result := route.Api_acl_post(tool.Config{IP: "acl_manager_ip", Session: &emptySession{}}, "", true, url.Values{
				"title_name":      {"BulkAclTarget"},
				"document_markup": {"plain"},
			})
			Markup_xss_response(t, manager_result, "ok")
			if got := tool.Get_document_setting_value(db, "BulkAclTarget", "document_markup", ""); got != "plain" {
				t.Fatalf("authorized bulk ACL did not update target markup: %q", got)
			}
		})
	}
}
