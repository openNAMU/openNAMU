package route

import "opennamu/route/tool"

func Api_list_recent_change(config tool.Config, set_type string, limit string, num string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if set_type == "edit" {
		set_type = ""
	}

	limit_int := tool.Str_to_int(limit)
	if limit_int > 50 || limit_int < 0 {
		limit_int = 50
	}

	page_int := tool.Str_to_int(num)
	if page_int > 0 {
		page_int = (page_int * limit_int) - limit_int
	} else {
		page_int = 0
	}

	var rc_count int
	tool.QueryRow_DB(
		db,
		"select count(*) from rc where type = ?",
		[]any{&rc_count},
		set_type,
	)
	history_auth := tool.Check_permission(db, "history_view", config.IP)

	id_title_list := [][]string{}
	if page_int < rc_count {
		rc_limit := limit_int
		if page_int+rc_limit > rc_count {
			rc_limit = rc_count - page_int
		}

		rows := tool.Query_DB(
			db,
			"select id, title from rc where type = ? order by date desc limit ?, ?",
			set_type,
			page_int,
			rc_limit,
		)
		for rows.Next() {
			var id string
			var title string
			if err := rows.Scan(&id, &title); err != nil {
				panic(err)
			}
			id_title_list = append(id_title_list, []string{id, title})
		}
		rows.Close()
	}

	history_offset := page_int - rc_count
	if history_offset < 0 {
		history_offset = 0
	}
	history_limit := limit_int - len(id_title_list)
	if history_auth && history_limit > 0 {
		history_query := `select h.id, h.title from history h where `
		history_values := []any{}
		if set_type != "normal" {
			history_query += "h.type = ? and "
			history_values = append(history_values, set_type)
		}
		history_query += `not exists (
			select 1 from rc r
			where r.type = ? and r.id = h.id and r.title = h.title
		) order by h.date desc limit ?, ?`
		history_values = append(history_values, set_type, history_offset, history_limit)

		rows := tool.Query_DB(db, history_query, history_values...)
		for rows.Next() {
			var id string
			var title string
			if err := rows.Scan(&id, &title); err != nil {
				panic(err)
			}
			id_title_list = append(id_title_list, []string{id, title})
		}
		rows.Close()
	}

	data_list := [][]string{}

	admin_auth := tool.Check_permission(db, "hidel", config.IP)
	ip_parser_temp := map[string][]string{}

	for _, id_title := range id_title_list {
		id := id_title[0]
		title := id_title[1]

		date := ""
		ip := ""
		send := ""
		leng := ""
		hide := ""
		type_data := ""
		tool.QueryRow_DB(
			db,
			"select date, ip, send, leng, hide, type from history where id = ? and title = ?",
			[]any{&date, &ip, &send, &leng, &hide, &type_data},
			id, title,
		)

		var ip_pre string
		var ip_render string

		if _, ok := ip_parser_temp[ip]; ok {
			ip_pre = ip_parser_temp[ip][0]
			ip_render = ip_parser_temp[ip][1]
		} else {
			ip_pre = tool.IP_preprocess(db, ip, config.IP)[0]
			ip_render = tool.IP_parser(db, ip, config.IP)

			ip_parser_temp[ip] = []string{ip_pre, ip_render}
		}

		if hide == "" || admin_auth {
			data_list = append(data_list, []string{
				id,
				title,
				date,
				ip_pre,
				send,
				leng,
				hide,
				ip_render,
				type_data,
			})
		} else {
			data_list = append(data_list, []string{"", "", "", "", "", "", hide, "", ""})
		}
	}

	return_data := make(map[string]any)
	return_data["response"] = "ok"
	return_data["language"] = api_history_language(db)
	return_data["auth"] = api_history_auth(db, config.IP)
	return_data["data"] = data_list

	return return_data
}
