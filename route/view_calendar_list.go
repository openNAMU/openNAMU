package route

import "opennamu/route/tool"

func View_calendar_list(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	api_data := Api_calendar_list(config)
	response, _ := api_data["response"].(string)
	if response == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	}

	data_list, _ := api_data["data"].([][]string)

	data_html := ""
	if tool.Check_permission(db, "calendar_create", config.IP) {
		data_html += `<a href="/calendar/add">(` + tool.Get_language(db, "add_calendar", true) + `)</a>`
	}

	data_html += "<ul>"
	for _, calendar := range data_list {
		if len(calendar) < 3 {
			continue
		}

		data_html += `<li><a href="/calendar/` + tool.Url_parser(calendar[0]) + `">` + tool.HTML_escape(calendar[1]) + `</a>`
		if calendar[2] != "" {
			data_html += ` - ` + tool.HTML_escape(calendar[2])
		}
		data_html += `</li>`
	}
	data_html += "</ul>"

	if len(data_list) == 0 {
		data_html += `<div>` + tool.Get_language(db, "calendar_list_empty", true) + `</div>`
	}

	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, "calendar", true),
		data_html,
		[]any{},
		[][]any{
			{"other", tool.Get_language(db, "return", true)},
		},
		map[string]string{},
	)
}
