package route

import (
	"opennamu/route/tool"
)

func View_edit_delete(config tool.Config, doc_name string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, doc_name, "", "document_delete", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	raw_data := Api_w_raw(config, doc_name, "true", "")
	response := raw_data["response"].(string)
	if response == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	} else if response != "ok" || !raw_data["exist"].(bool) {
		return tool.Get_redirect("/w/" + tool.Url_parser(doc_name))
	}

	data_html := `<form method="post">
        <input class="__ON_INPUT__" placeholder="` + tool.Get_language(db, "why", true) + `" name="send">
        <hr class="main_hr">
        ` + tool.Get_captcha_ui(db, config) + tool.Get_IP_warning_ui(db, config) + tool.Get_edit_check_box_ui(db) + tool.Get_edit_bottom_text_ui(db, "delete") + `
        <button class="__ON_BUTTON__" type="submit">` + tool.Get_language(db, "delete", true) + `</button>
    </form>`

	return tool.Get_template(
		db,
		config,
		doc_name,
		data_html,
		[]any{"(" + tool.Get_language(db, "delete", true) + ")"},
		[][]any{
			{"w/" + tool.Url_parser(doc_name), tool.Get_language(db, "return", true)},
		},
		map[string]string{},
	)
}
