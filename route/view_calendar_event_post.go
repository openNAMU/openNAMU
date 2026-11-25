package route

import "opennamu/route/tool"

func View_calendar_event_add_post(config tool.Config, id string, date string, title string, memo string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	api_data := Api_calendar_event_add(config, id, date, title, memo)
	response, _ := api_data["response"].(string)
	if response == "ok" {
		code, _ := api_data["data"].(string)
		month, _ := api_data["month"].(string)
		target := "/calendar/" + tool.Url_parser(id) + "/month/" + tool.Url_parser(month)
		if code != "" {
			target += "#event-" + tool.Url_parser(code)
		}

		return tool.Get_redirect(target)
	}
	if response == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	}
	if response == "not exist" {
		return tool.Get_redirect("/calendar")
	}

	return tool.Get_error_page(db, config, "error")
}

func View_calendar_event_edit_post(config tool.Config, id string, code string, date string, title string, memo string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	api_data := Api_calendar_event_edit(config, id, code, date, title, memo)
	response, _ := api_data["response"].(string)
	if response == "ok" {
		month, _ := api_data["month"].(string)
		target := "/calendar/" + tool.Url_parser(id) + "/month/" + tool.Url_parser(month)
		if code != "" {
			target += "#event-" + tool.Url_parser(code)
		}

		return tool.Get_redirect(target)
	}
	if response == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	}
	if response == "not exist" {
		return tool.Get_redirect("/calendar/" + tool.Url_parser(id))
	}

	return tool.Get_error_page(db, config, "error")
}
