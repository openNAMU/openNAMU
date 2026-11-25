package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_calendar_edit(config tool.Config, id string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if values == nil {
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
            <div><label for="calendar_name">` + tool.Get_language(db, "name", true) + `</label><input id="calendar_name" class="__ON_INPUT__" name="name" value="` + tool.HTML_escape(calendar_data["name"]) + `" required></div>
            <div><label for="calendar_subject">` + tool.Get_language(db, "explanation", true) + `</label><textarea id="calendar_subject" name="subject">` + tool.HTML_escape(calendar_data["subject"]) + `</textarea></div>
            <hr class="main_hr">
            <button class="__ON_BUTTON__" type="submit">` + tool.Get_language(db, "save", true) + `</button>
        </form>`

		return tool.Get_template(
			db,
			config,
			tool.Get_language(db, "edit", true),
			data_html,
			[]any{calendar_data["name"]},
			[][]any{
				{"calendar/" + tool.Url_parser(id), tool.Get_language(db, "return", true)},
			},
			map[string]string{},
		)
	}

	api_data := Api_calendar_edit(config, id, values.Get("name"), values.Get("subject"))
	response, _ := api_data["response"].(string)
	if response == "ok" {
		return tool.Get_redirect("/calendar/" + tool.Url_parser(id))
	}
	if response == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	}
	if response == "not exist" {
		return tool.Get_redirect("/calendar")
	}

	return tool.Get_error_page(db, config, "error")
}
