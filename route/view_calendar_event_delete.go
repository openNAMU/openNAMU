package route

import "opennamu/route/tool"

func View_calendar_event_delete(config tool.Config, id string, code string, do_delete bool) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !do_delete {
		api_data := Api_calendar_event(config, id, code)
		response, _ := api_data["response"].(string)
		if response == "require auth" {
			return tool.Get_error_page(db, config, "auth")
		}
		if response != "ok" {
			return tool.Get_redirect("/calendar/" + tool.Url_parser(id))
		}

		event_data, ok := api_data["data"].(map[string]string)
		if !ok {
			return tool.Get_redirect("/calendar/" + tool.Url_parser(id))
		}
		editable, _ := api_data["editable"].(bool)
		if !editable {
			return tool.Get_error_page(db, config, "auth")
		}

		data_html := `<div>` + tool.HTML_escape(event_data["date"]) + ` ` + tool.HTML_escape(event_data["title"]) + `</div>
        <hr class="main_hr">
        <form method="post">
            <span>` + tool.Get_language(db, "delete_warning", true) + `</span>
            <hr class="main_hr">
            <button class="__ON_BUTTON__" type="submit">` + tool.Get_language(db, "delete", true) + `</button>
        </form>`

		return tool.Get_template(
			db,
			config,
			tool.Get_language(db, "delete", true),
			data_html,
			[]any{event_data["title"]},
			[][]any{
				{"calendar/" + tool.Url_parser(id) + "/event/" + tool.Url_parser(code), tool.Get_language(db, "return", true)},
			},
			map[string]string{},
		)
	}

	api_data := Api_calendar_event_delete(config, id, code)
	response, _ := api_data["response"].(string)
	if response == "ok" {
		month, _ := api_data["month"].(string)
		if month != "" {
			return tool.Get_redirect("/calendar/" + tool.Url_parser(id) + "/month/" + tool.Url_parser(month))
		}

		return tool.Get_redirect("/calendar/" + tool.Url_parser(id))
	}
	if response == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	}
	if response == "not exist" {
		return tool.Get_redirect("/calendar/" + tool.Url_parser(id))
	}

	return tool.Get_error_page(db, config, "error")
}
