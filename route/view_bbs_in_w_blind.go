package route

import "opennamu/route/tool"

func View_bbs_in_w_blind(config tool.Config, set_id string, set_code string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "bbs_post_manage", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	title, exists := tool.Get_bbs_data_value(db, set_id, set_code, "title")
	if !exists {
		return tool.Get_redirect("/bbs/main")
	}

	action_name := "blind_post"
	action_value := "1"
	if bbs_post_blind(db, set_id, set_code) {
		action_name = "blind_post_release"
		action_value = "0"
	}

	data_html := "<form method='post'>" +
		"<input type='hidden' name='blind' value='" + action_value + "'>" +
		"<button class='__ON_BUTTON__' type='submit'>" + tool.Get_language(db, action_name, true) + "</button>" +
		"</form>"

	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, action_name, true),
		data_html,
		[]any{"(" + tool.HTML_escape(title) + ")"},
		[][]any{
			{"bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code), tool.Get_language(db, "return", true)},
		},
		map[string]string{},
	)
}

func View_bbs_in_w_blind_post(config tool.Config, set_id string, set_code string, blind string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if blind != "1" && blind != "0" {
		return tool.Get_error_page(db, config, "error")
	}

	api_data := Api_bbs_w_blind_post(config, set_id, set_code, blind == "1")
	response, _ := api_data["response"].(string)
	if response == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	}
	if response != "ok" {
		return tool.Get_error_page(db, config, "not found")
	}

	return tool.Get_redirect("/bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code))
}
