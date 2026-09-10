package route

import "opennamu/route/tool"

func View_bbs_in_w_comment_blind(config tool.Config, set_id string, set_code string, comment_code string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "bbs_comment_manage", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	comment_set_id, comment_set_code, exists := bbs_search_comment_location(set_id, set_code, comment_code)
	if !exists {
		return tool.Get_redirect("/bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code))
	}
	comment, comment_exists := tool.Get_bbs_data_value(db, comment_set_id, comment_set_code, "comment")
	if !comment_exists || comment == "" {
		return tool.Get_redirect("/bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code))
	}

	blind_data, hidden := tool.Get_bbs_data_value(db, comment_set_id, comment_set_code, "blind")
	hidden = hidden && blind_data == "O"
	action_name := "blind"
	action_value := "1"
	if hidden {
		action_name = "blind_release"
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
		[]any{"(" + tool.HTML_escape(set_code) + ") (#" + tool.HTML_escape(comment_code) + ")", 0},
		[][]any{
			{"bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code) + "#" + tool.Url_parser(comment_code), tool.Get_language(db, "return", true)},
		},
		map[string]string{},
	)
}

func View_bbs_in_w_comment_blind_post(config tool.Config, set_id string, set_code string, comment_code string, blind string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if blind != "1" && blind != "0" {
		return tool.Get_error_page(db, config, "error")
	}

	api_data := Api_bbs_w_comment_blind_post(config, set_id, set_code, comment_code, blind == "1")
	response, _ := api_data["response"].(string)
	if response == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	}
	if response != "ok" {
		return tool.Get_error_page(db, config, "not found")
	}

	return tool.Get_redirect("/bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code) + "#" + tool.Url_parser(comment_code))
}
