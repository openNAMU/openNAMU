package route

import "opennamu/route/tool"

func View_list_view_page(config tool.Config, num string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	api_data := Api_list_view_page(config, num)
	data_list := api_data["data"].([][]string)
	data_html := ""

	for _, data := range data_list {
		if len(data) < 2 {
			continue
		}

		left := `<a href="/w/` + tool.Url_parser(data[0]) + `">` + tool.HTML_escape(data[0]) + `</a>`
		right := tool.Get_language(db, "page_view", true) + " : " + tool.HTML_escape(data[1])
		data_html += tool.Get_list_ui(left, right, "", "")
	}

	data_html += tool.Get_page_control(
		db,
		tool.Str_to_int(num),
		len(data_list),
		50,
		"/list/document/view/{}",
	)

	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, "page_view_sort", true),
		data_html,
		[]any{},
		[][]any{{"other", tool.Get_language(db, "return", true)}},
		map[string]string{},
	)
}
