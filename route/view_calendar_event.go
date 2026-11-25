package route

import (
	"strings"

	"opennamu/route/tool"
)

func View_calendar_event(config tool.Config, id string, code string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

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

	month := ""
	if len(event_data["date"]) >= 7 {
		month = event_data["date"][:7]
	}

	data_html := `<h2>` + tool.HTML_escape(event_data["title"]) + `</h2>`
	data_html += `<div>` + tool.HTML_escape(event_data["date"]) + `</div>`
	if event_data["memo"] != "" {
		data_html += `<hr class="main_hr"><div>` + strings.ReplaceAll(tool.HTML_escape(event_data["memo"]), "\n", "<br>") + `</div>`
	}
	if event_data["user_render"] != "" {
		data_html += `<hr class="main_hr"><div>` + tool.Get_language(db, "user", true) + `: ` + event_data["user_render"] + `</div>`
	}

	if editable {
		data_html += `<hr class="main_hr"><h3>` + tool.Get_language(db, "edit", true) + `</h3>
        <form method="post" action="/calendar/` + tool.Url_parser(id) + `/event/` + tool.Url_parser(code) + `">
            <div><label for="calendar_event_date">` + tool.Get_language(db, "date", true) + `</label> <input id="calendar_event_date" class="__ON_INPUT__" type="date" name="date" value="` + tool.HTML_escape(event_data["date"]) + `" required></div>
            <div><label for="calendar_event_title">` + tool.Get_language(db, "title", true) + `</label> <input id="calendar_event_title" class="__ON_INPUT__" name="title" value="` + tool.HTML_escape(event_data["title"]) + `" required></div>
            <div><label for="calendar_event_memo">` + tool.Get_language(db, "content", true) + `</label><textarea id="calendar_event_memo" name="memo">` + tool.HTML_escape(event_data["memo"]) + `</textarea></div>
            <hr class="main_hr">
            <button class="__ON_BUTTON__" type="submit">` + tool.Get_language(db, "save", true) + `</button>
        </form>
        <hr class="main_hr">
        <a href="/calendar/` + tool.Url_parser(id) + `/event/` + tool.Url_parser(code) + `/delete">` + tool.Get_language(db, "delete", true) + `</a>`
	}

	menu := [][]any{
		{"calendar/" + tool.Url_parser(id), tool.Get_language(db, "return", true)},
	}
	if month != "" {
		menu = [][]any{
			{"calendar/" + tool.Url_parser(id) + "/month/" + tool.Url_parser(month), tool.Get_language(db, "return", true)},
		}
	}

	return tool.Get_template(
		db,
		config,
		event_data["title"],
		data_html,
		[]any{event_data["date"]},
		menu,
		map[string]string{},
	)
}
