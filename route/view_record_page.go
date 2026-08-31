package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func View_record_page(config tool.Config, user_name string, record_type string, page string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if user_name == "" {
		user_name = config.IP
	}
	if user_name != config.IP && !tool.Check_permission(db, "hidel", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	if record_type == "" {
		record_type = "edit"
	}

	page_num := tool.Str_to_int(page)
	if page_num < 1 {
		page_num = 1
	}
	offset := (page_num - 1) * 50
	data_html := ""
	count := 0

	if record_type == "topic" {
		rows := tool.Query_DB(db, "select code, id, date from topic where ip = ? order by date desc limit ?, 50", user_name, offset)
		for rows.Next() {
			code, comment_id, date := "", "", ""
			if rows.Scan(&code, &comment_id, &date) != nil {
				continue
			}

			topic_title, topic_sub := "", ""
			tool.QueryRow_DB(db, "select title, sub from rd where code = ?", []any{&topic_title, &topic_sub}, code)
			left := `<a href="/thread/` + tool.Url_parser(code) + `#` + tool.Url_parser(comment_id) + `">` + tool.HTML_escape(topic_sub+"#"+comment_id) + `</a>`
			left += " (" + tool.HTML_escape(topic_title) + ")"
			right := tool.IP_parser(db, user_name, config.IP) + " | " + tool.HTML_escape(date)
			data_html += tool.Get_list_ui(left, right, "", "")
			count++
		}
		rows.Close()
		data_html += tool.Get_page_control(db, page_num, count, 50, "/record/topic/"+tool.Url_parser(user_name)+"/{}")

		return tool.Get_template(
			db,
			config,
			user_name,
			data_html,
			[]any{"(" + tool.Get_language(db, "discussion_record", true) + ")"},
			[][]any{
				{"other", tool.Get_language(db, "other", true)},
				{"user/" + tool.Url_parser(user_name), tool.Get_language(db, "user_tool", true)},
			},
			map[string]string{},
		)
	}

	var rows *sql.Rows
	if record_type == "normal" || record_type == "edit" {
		rows = tool.Query_DB(db, "select id, title, date, ip, send, leng, hide, type from history where ip = ? order by date desc limit ?, 50", user_name, offset)
	} else {
		rows = tool.Query_DB(db, "select id, title, date, ip, send, leng, hide, type from history where ip = ? and type = ? order by date desc limit ?, 50", user_name, record_type, offset)
	}

	data_list := [][]string{}
	ip_cache := map[string][]string{}
	can_view_hidden := tool.Check_permission(db, "hidel", config.IP)
	for rows.Next() {
		id, title, date, ip, send, leng, hide, type_data := "", "", "", "", "", "", "", ""
		if rows.Scan(&id, &title, &date, &ip, &send, &leng, &hide, &type_data) != nil {
			continue
		}
		if hide != "" && !can_view_hidden {
			data_list = append(data_list, []string{"", "", "", "", "", "", hide, "", ""})
			continue
		}

		ip_data, ok := ip_cache[ip]
		if !ok {
			ip_pre := tool.IP_preprocess(db, ip, config.IP)
			ip_data = []string{"", tool.IP_parser(db, ip, config.IP)}
			if len(ip_pre) > 0 {
				ip_data[0] = ip_pre[0]
			}
			ip_cache[ip] = ip_data
		}
		data_list = append(data_list, []string{id, title, date, ip_data[0], send, leng, hide, ip_data[1], type_data})
	}
	rows.Close()

	data_html, _ = Get_ui_history(db, data_list)
	count = len(data_list)
	data_html += tool.Get_page_control(db, page_num, count, 50, "/record/{}/"+tool.Url_parser(record_type)+"/"+tool.Url_parser(user_name))

	menu_html := ""
	for _, option := range []string{"normal", "edit", "move", "delete", "revert", "r1", "file", "category"} {
		menu_html += `<a href="/record/1/` + option + `/` + tool.Url_parser(user_name) + `">(` + tool.Get_language(db, option, true) + `)</a> `
	}
	data_html = menu_html + data_html

	menu := [][]any{{"user/" + tool.Url_parser(user_name), tool.Get_language(db, "user_tool", true)}}
	if tool.Check_permission(db, "record_manage", config.IP) {
		menu = append(menu, []any{"record/reset/" + tool.Url_parser(user_name), tool.Get_language(db, "record_reset", true)})
	}

	return tool.Get_template(
		db,
		config,
		user_name,
		data_html,
		[]any{"(" + tool.Get_language(db, "edit_record", true) + ")"},
		menu,
		map[string]string{},
	)
}
