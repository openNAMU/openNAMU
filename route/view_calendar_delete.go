package route

import "opennamu/route/tool"

func View_calendar_delete(config tool.Config, id string, do_delete bool) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !do_delete {
		api_data := Api_calendar_view(config, id, "")
		response, _ := api_data["response"].(string)
		if response == "require auth" {
			return tool.Get_error_page(db, config, "auth")
		}
		if response != "ok" {
			return tool.Get_redirect("/calendar")
		}

		calendar_data, ok := api_data["data"].(map[string]string)
		if !ok {
			return tool.Get_redirect("/calendar")
		}
		if calendar_data["owner"] != config.IP && !tool.Check_permission(db, "calendar_manage", config.IP) {
			return tool.Get_error_page(db, config, "auth")
		}

		data_html := `<form method="post">
            <span>` + tool.Get_language(db, "delete_warning", true) + `</span>
            <hr class="main_hr">
            <button class="__ON_BUTTON__" type="submit">` + tool.Get_language(db, "delete", true) + `</button>
        </form>`

		return tool.Get_template(
			db,
			config,
			tool.Get_language(db, "delete", true),
			data_html,
			[]any{calendar_data["name"]},
			[][]any{
				{"calendar/" + tool.Url_parser(id), tool.Get_language(db, "return", true)},
			},
			map[string]string{},
		)
	}

	api_data := Api_calendar_delete(config, id)
	response, _ := api_data["response"].(string)
	if response == "ok" {
		return tool.Get_redirect("/calendar")
	}
	if response == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	}
	if response == "not exist" {
		return tool.Get_redirect("/calendar")
	}

	return tool.Get_error_page(db, config, "error")
}
