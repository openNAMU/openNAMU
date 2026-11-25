package route

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"opennamu/route/tool"
)

func Calendar_grid_html(db *sql.DB, id string, month string, events [][]string) string {
	leading, days := tool.Calendar_month_grid(month)
	if days == 0 {
		return ""
	}

	events_by_date := map[string][][]string{}
	for _, event := range events {
		if len(event) < 3 {
			continue
		}
		events_by_date[event[1]] = append(events_by_date[event[1]], event)
	}

	weekday_keys := []string{
		"weekday_monday",
		"weekday_tuesday",
		"weekday_wednesday",
		"weekday_thursday",
		"weekday_friday",
		"weekday_saturday",
		"weekday_sunday",
	}

	data_html := `<table class="calendar_table" style="width: 100%; table-layout: fixed;"><tr>`
	for _, weekday_key := range weekday_keys {
		data_html += `<td style="text-align: center; padding: 2px 4px;"><b>` + tool.Get_language(db, weekday_key, true) + `</b></td>`
	}
	data_html += `</tr><tr>`

	cell_count := 0
	for i := 0; i < leading; i++ {
		data_html += `<td style="vertical-align: top; overflow: hidden; padding: 2px 4px;"></td>`
		cell_count++
	}

	today := tool.Get_date()
	for day := 1; day <= days; day++ {
		date := fmt.Sprintf("%s-%02d", month, day)
		day_html := strconv.Itoa(day)
		if date == today {
			day_html = `<b>` + day_html + `</b>`
		}

		data_html += `<td style="vertical-align: top; overflow: hidden; padding: 2px 4px;"><div>` + day_html + `</div>`

		day_events := events_by_date[date]
		event_limit := len(day_events)
		if event_limit > 3 {
			event_limit = 3
		}
		for i := 0; i < event_limit; i++ {
			event := day_events[i]
			data_html += `<div style="overflow: hidden; text-overflow: ellipsis; white-space: nowrap;"><a href="/calendar/` + tool.Url_parser(id) + `/event/` + tool.Url_parser(event[0]) + `">` + tool.HTML_escape(event[2]) + `</a></div>`
		}
		if len(day_events) > event_limit {
			data_html += `<div>+` + strconv.Itoa(len(day_events)-event_limit) + `</div>`
		}

		data_html += `</td>`
		cell_count++
		if cell_count%7 == 0 && day < days {
			data_html += `</tr><tr>`
		}
	}

	for cell_count%7 != 0 {
		data_html += `<td style="vertical-align: top; overflow: hidden; padding: 2px 4px;"></td>`
		cell_count++
	}
	data_html += `</tr></table>`

	return data_html
}

func View_calendar(config tool.Config, id string, month string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	api_data := Api_calendar_view(config, id, month)
	response, _ := api_data["response"].(string)
	if response == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	}
	if response != "ok" {
		if response == "not exist" {
			return tool.Get_redirect("/calendar")
		}
		return tool.Get_redirect("/calendar/" + tool.Url_parser(id))
	}

	calendar_data, ok := api_data["data"].(map[string]string)
	if !ok {
		return tool.Get_redirect("/calendar")
	}
	events, _ := api_data["events"].([][]string)

	month_value := calendar_data["month"]
	can_manage := calendar_data["owner"] == config.IP || tool.Check_permission(db, "calendar_manage", config.IP)

	data_html := `<h2>` + tool.HTML_escape(calendar_data["name"]) + `</h2>`
	if calendar_data["subject"] != "" {
		data_html += `<div>` + tool.HTML_escape(calendar_data["subject"]) + `</div>`
	}
	data_html += `<div>(<a href="/calendar/` + tool.Url_parser(id) + `/month/` + tool.Url_parser(calendar_data["prev_month"]) + `">` + tool.Get_language(db, "previous", true) + `</a>) `
	data_html += `(<a href="/calendar/` + tool.Url_parser(id) + `/month/` + tool.Url_parser(calendar_data["next_month"]) + `">` + tool.Get_language(db, "next", true) + `</a>) `
	data_html += `(<a href="/calendar/` + tool.Url_parser(id) + `">` + tool.Get_language(db, "record_this_month", true) + `</a>)</div>`
	if can_manage {
		data_html += `<div>(<a href="/calendar/` + tool.Url_parser(id) + `/edit">` + tool.Get_language(db, "edit", true) + `</a>) `
		data_html += `(<a href="/calendar/` + tool.Url_parser(id) + `/delete">` + tool.Get_language(db, "delete", true) + `</a>)</div>`
	}

	data_html += `<hr class="main_hr">` + Calendar_grid_html(db, id, month_value, events)

	data_html += `<hr class="main_hr"><h3>` + tool.Get_language(db, "calendar_event", true) + `</h3>`
	if len(events) == 0 {
		data_html += `<div>` + tool.Get_language(db, "calendar_empty", true) + `</div>`
	} else {
		data_html += `<ul>`
		for _, event := range events {
			if len(event) < 6 {
				continue
			}

			data_html += `<li id="event-` + tool.Url_parser(event[0]) + `"><b>` + tool.HTML_escape(event[1]) + `</b> <a href="/calendar/` + tool.Url_parser(id) + `/event/` + tool.Url_parser(event[0]) + `">` + tool.HTML_escape(event[2]) + `</a>`
			if event[3] != "" {
				data_html += `<div>` + strings.ReplaceAll(tool.HTML_escape(event[3]), "\n", "<br>") + `</div>`
			}
			if event[5] != "" {
				data_html += `<div>` + event[5] + `</div>`
			}
			data_html += `</li>`
		}
		data_html += `</ul>`
	}

	if tool.Check_permission(db, "calendar_create", config.IP) {
		default_date := tool.Get_date()
		if month_value != tool.Get_month() {
			default_date = month_value + "-01"
		}

		data_html += `<hr class="main_hr"><h3>` + tool.Get_language(db, "calendar_event_add", true) + `</h3>
        <form method="post" action="/calendar/` + tool.Url_parser(id) + `/event/add">
            <div><label for="calendar_event_date">` + tool.Get_language(db, "date", true) + `</label> <input id="calendar_event_date" class="__ON_INPUT__" type="date" name="date" value="` + tool.HTML_escape(default_date) + `" required></div>
            <div><label for="calendar_event_title">` + tool.Get_language(db, "title", true) + `</label> <input id="calendar_event_title" class="__ON_INPUT__" name="title" required></div>
            <div><label for="calendar_event_memo">` + tool.Get_language(db, "content", true) + `</label><textarea id="calendar_event_memo" name="memo"></textarea></div>
            <hr class="main_hr">
            <button class="__ON_BUTTON__" type="submit">` + tool.Get_language(db, "add", true) + `</button>
        </form>`
	}

	return tool.Get_template(
		db,
		config,
		calendar_data["name"],
		data_html,
		[]any{month_value},
		[][]any{
			{"calendar", tool.Get_language(db, "return", true)},
		},
		map[string]string{},
	)
}
