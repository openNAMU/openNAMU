package route

import "opennamu/route/tool"

func View_edit_revert(config tool.Config, doc_name string, rev string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, doc_name, "", "document_edit", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	data, hide, exists := tool.Get_history_content(db, doc_name, rev)
	if !exists {
		return tool.Get_redirect("/w/" + tool.Url_parser(doc_name))
	}
	if hide != "" && !tool.Check_permission(db, "hidel", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	data_html := `<form method="post">
        <input class="__ON_INPUT__" placeholder="` + tool.Get_language(db, "why", true) + `" name="send">
        <hr class="main_hr">
        ` + tool.Get_captcha_ui(db, config) + tool.Get_IP_warning_ui(db, config) + tool.Get_edit_check_box_ui(db) + tool.Get_edit_bottom_text_ui(db, "revert") + `
        <button class="__ON_BUTTON__" type="submit">` + tool.Get_language(db, "revert", true) + `</button>
    </form>
    <hr class="main_hr"><pre>` + tool.HTML_escape(data) + `</pre>`

	return tool.Get_template(
		db,
		config,
		doc_name,
		data_html,
		[]any{"(r" + tool.HTML_escape(rev) + ") (" + tool.Get_language(db, "revert", true) + ")"},
		[][]any{
			{"history/" + tool.Url_parser(doc_name), tool.Get_language(db, "history", true)},
			{"recent_changes", tool.Get_language(db, "recent_change", true)},
		},
		map[string]string{},
	)
}
