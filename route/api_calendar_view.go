package route

import (
	"opennamu/route/tool"
)

func Api_calendar_view(config tool.Config, id string, month string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	return_data["events"] = [][]string{}
	if !Calendar_view_allowed(db, config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	id, id_ok := Calendar_id_normalize(id)
	if !id_ok {
		return_data["response"] = "not exist"
		return_data["data"] = "calendar"
		return return_data
	}

	calendar_data, calendar_exists := Calendar_meta_get(db, id)
	if !calendar_exists {
		return_data["response"] = "not exist"
		return_data["data"] = "calendar"
		return return_data
	}

	if month == "" {
		month = tool.Get_month()
	}
	month, month_ok := tool.Calendar_month_normalize(month)
	if !month_ok {
		return_data["response"] = "error"
		return_data["data"] = "invalid date"
		return return_data
	}

	month_start, next_month_start := tool.Calendar_month_range(month)
	if month_start == "" || next_month_start == "" {
		return_data["response"] = "error"
		return_data["data"] = "invalid date"
		return return_data
	}

	rows := tool.Query_DB(
		db,
		`select d.set_code, d.set_data, coalesce(t.set_data, ''), coalesce(m.set_data, ''), coalesce(u.set_data, '')
		from calendar d
		left join calendar t on t.set_id = d.set_id and t.set_code = d.set_code and t.set_name = 'title'
		left join calendar m on m.set_id = d.set_id and m.set_code = d.set_code and m.set_name = 'memo'
		left join calendar u on u.set_id = d.set_id and u.set_code = d.set_code and u.set_name = 'user'
		where d.set_id = ? and d.set_name = 'date' and d.set_data >= ? and d.set_data < ?
		order by d.set_data asc, d.set_code + 0 asc`,
		id,
		month_start,
		next_month_start,
	)
	defer rows.Close()

	events := [][]string{}
	for rows.Next() {
		code := ""
		date := ""
		title := ""
		memo := ""
		user := ""
		if err := rows.Scan(&code, &date, &title, &memo, &user); err != nil {
			panic(err)
		}

		user_render := ""
		if user != "" {
			user_render = tool.Get_user_profile_image_ui(db, user) + tool.IP_parser(db, user, config.IP)
		}

		events = append(events, []string{code, date, title, memo, user, user_render})
	}

	data := map[string]string{
		"name":       calendar_data["name"],
		"subject":    calendar_data["subject"],
		"owner":      calendar_data["owner"],
		"month":      month,
		"prev_month": tool.Calendar_month_shift(month, -1),
		"next_month": tool.Calendar_month_shift(month, 1),
	}

	return_data["response"] = "ok"
	return_data["data"] = data
	return_data["events"] = events
	return return_data
}
