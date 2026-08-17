package route

import "opennamu/route/tool"

func View_vote_list(config tool.Config, type_str string, num_str string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if type_str == "" {
		type_str = "open"
	}

	api_data := Api_vote_list(config, type_str, num_str)
	data_list := api_data["data"].([][]string)

	data_html := ""
	sub := any(0)
	menu := [][]any{
		{"other", tool.Get_language(db, "return", true)},
	}
	if type_str == "open" {
		data_html += `<a href="/vote/list/close">(` + tool.Get_language(db, "close_vote_list", true) + `)</a>`
		if tool.Check_acl(db, "", "", "vote", config.IP) {
			menu = append(menu, []any{"vote/add", tool.Get_language(db, "add_vote", true)})
		}
	} else {
		data_html += `<a href="/vote">(` + tool.Get_language(db, "open_vote_list", true) + `)</a>`
		sub = "(" + tool.Get_language(db, "closed", true) + ")"
	}

	data_html += "<ul>"
	for _, in_data := range data_list {
		open_select := "not_open_vote"
		if (type_str == "open" && in_data[2] == "open") || (type_str != "open" && in_data[2] == "close") {
			open_select = "open_vote"
		}
		data_html += `<li>` + tool.HTML_escape(in_data[1]) + `. <a href="/vote/` + tool.Url_parser(in_data[1]) + `">` + tool.HTML_escape(in_data[0]) + `</a> (` + tool.Get_language(db, open_select, true) + `)</li>`
	}
	data_html += "</ul>"

	page := tool.Str_to_int(num_str)
	if page < 1 {
		page = 1
	}
	page_url := "/vote/list/{}"
	if type_str != "open" {
		page_url = "/vote/list/close/{}"
	}
	if page > 1 || len(data_list) == 50 {
		data_html += tool.Get_page_control(db, page, len(data_list), 50, page_url)
	}

	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, "vote_list", true),
		data_html,
		[]any{sub},
		menu,
		map[string]string{},
	)
}
