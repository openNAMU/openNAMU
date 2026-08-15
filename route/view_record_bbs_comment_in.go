package route

import "opennamu/route/tool"

func View_record_bbs_comment_in(config tool.Config, user_name string, set_id string, page string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	api_data := Api_record_bbs_comment_in(config, user_name, set_id, page)
	if api_data["response"].(string) != "ok" {
		return tool.Get_error_page(db, config, "auth")
	}

	data_list := api_data["data"].([][]string)
	data_html := ""
	bbs_name := Api_bbs_num_to_name(db, set_id)["data"].(string)

	for _, data := range data_list {
		bbs_id := data[0]
		post_id := data[1]
		comment_id := data[2] + "-" + data[3]
		api_data := Api_bbs_w(config, set_id, post_id)
		api_data_in := api_data["data"].(map[string]string)
		title := api_data_in["title"]
		date := api_data_in["date"]
		link := `<a href="/bbs/w/` + tool.Url_parser(bbs_id) + `/` + tool.Url_parser(post_id) + `">` + tool.HTML_escape(title) + `</a>`
		data_html += tool.Get_list_ui(link, date, comment_id, "")
	}
	data_html += tool.Get_page_control(db, tool.Str_to_int(page), len(data_list), 50, "/record_bbs_comment_in/"+tool.Url_parser(set_id)+"/"+tool.Url_parser(user_name)+"/{}")

	return tool.Get_template(
		db,
		config,
		user_name,
		data_html,
		[]any{"(" + bbs_name + ") (" + tool.Get_language(db, "bbs_comment_record", true) + ")"},
		[][]any{{"user/" + tool.Url_parser(user_name), tool.Get_language(db, "user_tool", true)}},
		map[string]string{},
	)
}
