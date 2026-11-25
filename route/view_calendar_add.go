package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_calendar_add(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if values == nil {
		if !tool.Check_permission(db, "calendar_create", config.IP) {
			return tool.Get_error_page(db, config, "auth")
		}

		data_html := `<form method="post">
            <div><label for="calendar_name">` + tool.Get_language(db, "name", true) + `</label><input id="calendar_name" class="__ON_INPUT__" name="name" required></div>
            <div><label for="calendar_subject">` + tool.Get_language(db, "explanation", true) + `</label><textarea id="calendar_subject" name="subject"></textarea></div>
            <hr class="main_hr">
            <button class="__ON_BUTTON__" type="submit">` + tool.Get_language(db, "send", true) + `</button>
        </form>`

		return tool.Get_template(
			db,
			config,
			tool.Get_language(db, "add_calendar", true),
			data_html,
			[]any{},
			[][]any{
				{"calendar", tool.Get_language(db, "return", true)},
			},
			map[string]string{},
		)
	}

	api_data := Api_calendar_add(config, values.Get("name"), values.Get("subject"))
	response, _ := api_data["response"].(string)
	if response == "ok" {
		id, _ := api_data["data"].(string)
		return tool.Get_redirect("/calendar/" + tool.Url_parser(id))
	}
	if response == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	}

	return tool.Get_error_page(db, config, "error")
}
