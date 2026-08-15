package route

import "opennamu/route/tool"

func View_record_bbs(config tool.Config, user_name string, page string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	api_data := Api_record_bbs(config, user_name, page)
	if api_data["response"].(string) != "ok" {
		return tool.Get_error_page(db, config, "auth")
	}

	data_list := api_data["data"].([][]string)
	data_html := ""

	for _, data := range data_list {
		bbs_name := Api_bbs_num_to_name(db, data[0])["data"].(string)
		set_id := data[0]
		date := data[1]
		link := `<a href="/record_bbs_in/` + tool.Url_parser(set_id) + `/` + tool.Url_parser(user_name) + `">` + tool.HTML_escape(bbs_name) + `</a>`
		data_html += tool.Get_list_ui(link, date, "", "")
	}
	data_html += tool.Get_page_control(db, tool.Str_to_int(page), len(data_list), 50, "/record/bbs/"+tool.Url_parser(user_name)+"/{}")

	return tool.Get_template(
		db,
		config,
		user_name,
		data_html,
		[]any{"(" + tool.Get_language(db, "bbs_record", true) + ")"},
		[][]any{{"user/" + tool.Url_parser(user_name), tool.Get_language(db, "user_tool", true)}},
		map[string]string{},
	)
}
