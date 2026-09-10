package route

import "opennamu/route/tool"

func View_bbs_in_w_comment_tabom(config tool.Config, set_id string, set_code string, comment_code string, vote_type string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, set_id, "", "bbs_comment", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	if vote_type != "down" {
		vote_type = "up"
	}
	vote_language := "upvote"
	if vote_type == "down" {
		vote_language = "downvote"
	}

	data_html := `<form method="post" action="/bbs/w/` + tool.Url_parser(set_id) + `/` + tool.Url_parser(set_code) + `/comment/` + tool.Url_parser(comment_code) + `/tabom"><input type="hidden" name="vote_type" value="` + vote_type + `"><button class="__ON_BUTTON__" type="submit">` + tool.Get_language(db, vote_language, true) + `</button></form>`
	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, vote_language, true),
		data_html,
		[]any{"(#" + tool.HTML_escape(comment_code) + ")", 0},
		[][]any{{"bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code) + "#" + tool.Url_parser(comment_code), tool.Get_language(db, "return", true)}},
		map[string]string{},
	)
}

func View_bbs_in_w_comment_tabom_post(config tool.Config, set_id string, set_code string, comment_code string, vote_type string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	data := Api_bbs_w_comment_tabom_post(config, set_id, set_code, comment_code, vote_type)
	response, _ := data["response"].(string)
	if response == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	}
	if response == "not exist" {
		return tool.Get_error_page(db, config, "not found")
	}

	return tool.Get_redirect("/bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code) + "#" + tool.Url_parser(comment_code))
}
