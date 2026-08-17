package route

import "opennamu/route/tool"

func View_bbs_edit(config tool.Config, set_id string, set_code string, comment_code string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	api_data := Api_bbs_w_edit_view(config, set_id, set_code, comment_code)
	response, _ := api_data["response"].(string)
	if response != "ok" {
		if response == "require auth" {
			return tool.Get_error_page(db, config, "auth")
		}

		return tool.Get_redirect("/bbs/main")
	}

	data, ok := api_data["data"].(map[string]string)
	if !ok {
		return tool.Get_redirect("/bbs/main")
	}

	path := "/bbs/edit/" + tool.Url_parser(set_id)
	if set_code != "" {
		path += "/" + tool.Url_parser(set_code)
	}
	if comment_code != "" {
		path += "/" + tool.Url_parser(comment_code)
	}

	title := tool.Get_language(db, "post_edit", true)
	if comment_code != "" {
		title = tool.Get_language(db, "bbs_comment_edit", true)
	} else if set_code == "" {
		title = tool.Get_language(db, "post_add", true)
	}

	title_style := ""
	if comment_code != "" {
		title_style = ` style="display: none;"`
	}

	data_html := `<a href="/filter/edit_filter">(` + tool.Get_language(db, "edit_filter_rule", true) + `)</a><hr class="main_hr">
        <form action="` + path + `" method="post">
            <input class="__ON_INPUT__"` + title_style + ` placeholder="` + tool.Get_language(db, "title", true) + `" name="title" value="` + tool.HTML_escape(data["title"]) + `">
            <hr` + title_style + ` class="main_hr">
            ` + tool.Get_editor_ui(db, config, data["data"], "bbs", "", "") + `
        </form>`

	return tool.Get_template(
		db,
		config,
		title,
		data_html,
		[]any{},
		[][]any{
			{"bbs/in/" + tool.Url_parser(set_id), tool.Get_language(db, "return", true)},
		},
		map[string]string{},
	)
}
