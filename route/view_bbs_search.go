package route

import "opennamu/route/tool"

func View_bbs_search(config tool.Config, set_id string, keyword string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	bbs_name := Api_bbs_num_to_name(db, set_id)["data"].(string)
	if bbs_name == "" {
		return tool.Get_redirect("/bbs/main")
	}

	data_html := `<form method="post" action="/bbs/search/` + tool.Url_parser(set_id) + `">
        <input class="__ON_INPUT__" name="keyword" value="` + tool.HTML_escape(keyword) + `" placeholder="` + tool.Get_language(db, "search", true) + `">
        <button class="__ON_BUTTON__" type="submit">` + tool.Get_language(db, "search", true) + `</button>
    </form><hr class="main_hr">`

	if keyword != "" {
		data_api := Api_bbs_search(config, keyword, set_id)
		data_html += Get_bbs_list_ui(config, data_api["data"].([]map[string]string), map[string]string{set_id: bbs_name})
	}

	return tool.Get_template(
		db,
		config,
		bbs_name,
		data_html,
		[]any{},
		[][]any{
			{"bbs/in/" + tool.Url_parser(set_id), tool.Get_language(db, "return", true)},
		},
		map[string]string{},
	)
}
