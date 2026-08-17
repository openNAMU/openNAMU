package route

import "opennamu/route/tool"

func View_bbs_pinned(config tool.Config, set_id string, set_code string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	api_data := Api_bbs_w_pinned(config, set_id, set_code, false)
	if api_data["response"] != "ok" {
		if api_data["response"] == "require auth" {
			return tool.Get_redirect("/bbs/in/" + tool.Url_parser(set_id))
		}
		return tool.Get_redirect("/bbs/main")
	}

	state, _ := api_data["data"].(string)
	button_text := tool.Get_language(db, state, true)
	bbs_name_data := Api_bbs_num_to_name(db, set_id)
	bbs_name, _ := bbs_name_data["data"].(string)

	data_html := `<form method="post">
        <button class="__ON_BUTTON__" type="submit">` + button_text + `</button>
    </form>`

	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, "bbs_post_pinned", true),
		data_html,
		[]any{"(" + bbs_name + ") (" + tool.HTML_escape(set_code) + ")", 0},
		[][]any{
			{"bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code), tool.Get_language(db, "return", true)},
		},
		map[string]string{},
	)
}
