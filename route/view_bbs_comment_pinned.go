package route

import "opennamu/route/tool"

func View_bbs_comment_pinned(config tool.Config, set_id string, set_code string, comment_code string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	api_data := Api_bbs_w_comment_pinned(config, set_id, set_code, comment_code, false)
	if api_data["response"] != "ok" {
		return tool.Get_redirect("/bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code))
	}

	state, _ := api_data["data"].(string)
	button_text := tool.Get_language(db, state, true)

	data_html := `<form method="post">
        <button class="__ON_BUTTON__" type="submit">` + button_text + `</button>
    </form>`

	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, "bbs_comment_pinned", true),
		data_html,
		[]any{"(" + tool.HTML_escape(set_code) + ") (#" + tool.HTML_escape(comment_code) + ")", 0},
		[][]any{
			{"bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code) + "#" + tool.Url_parser(comment_code), tool.Get_language(db, "return", true)},
		},
		map[string]string{},
	)
}
