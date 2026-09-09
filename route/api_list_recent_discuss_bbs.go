package route

import "opennamu/route/tool"

func api_list_recent_discuss_bbs(config tool.Config, limit string, num string, set_type string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, thread_bbs_id, "", "bbs_view", config.IP) {
		return map[string]any{"response": "require auth", "data": [][]string{}}
	}

	limit_int := tool.Str_to_int(limit)
	if limit_int > 50 || limit_int < 1 {
		limit_int = 50
	}
	page := tool.Str_to_int(num)
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit_int

	condition := ""
	switch set_type {
	case "close":
		condition = " and coalesce((select set_data from bbs_data prefix_data where prefix_data.set_name = 'prefix' and prefix_data.set_id = document_data.set_id and prefix_data.set_code = document_data.set_code limit 1), '') = '닫힘'"
	case "open":
		condition = " and coalesce((select set_data from bbs_data prefix_data where prefix_data.set_name = 'prefix' and prefix_data.set_id = document_data.set_id and prefix_data.set_code = document_data.set_code limit 1), '') != '닫힘'"
	}

	rows := tool.Query_DB(
		db,
		"select document_data.set_data, coalesce((select set_data from bbs_data title_data where title_data.set_name = 'title' and title_data.set_id = document_data.set_id and title_data.set_code = document_data.set_code limit 1), ''), coalesce((select set_data from bbs_data date_data where date_data.set_name = 'date' and date_data.set_id = document_data.set_id and date_data.set_code = document_data.set_code limit 1), ''), document_data.set_code, case when coalesce((select set_data from bbs_data prefix_data where prefix_data.set_name = 'prefix' and prefix_data.set_id = document_data.set_id and prefix_data.set_code = document_data.set_code limit 1), '') = '닫힘' then 'O' else '' end, coalesce((select set_data from bbs_data agree_data where agree_data.set_name = 'topic_agree' and agree_data.set_id = document_data.set_id and agree_data.set_code = document_data.set_code limit 1), '') from bbs_data document_data where document_data.set_id = ? and document_data.set_name = 'document'"+condition+" order by (select set_data from bbs_data date_data where date_data.set_name = 'date' and date_data.set_id = document_data.set_id and date_data.set_code = document_data.set_code limit 1) desc limit ?, ?",
		thread_bbs_id,
		offset,
		limit_int,
	)
	defer rows.Close()

	data_list := [][]string{}
	for rows.Next() {
		var title string
		var sub string
		var date string
		var code string
		var stop string
		var agree string
		if rows.Scan(&title, &sub, &date, &code, &stop, &agree) != nil {
			continue
		}

		comment_set_id := thread_bbs_id + "-" + code
		ip := ""
		id := ""
		tool.QueryRow_DB(db, "select set_data from bbs_data where set_name = 'comment_user_id' and set_id = ? order by set_code + 0 desc limit 1", []any{&ip}, comment_set_id)
		tool.QueryRow_DB(db, "select set_code from bbs_data where set_name = 'comment' and set_id = ? order by set_code + 0 desc limit 1", []any{&id}, comment_set_id)

		ip_pre := ""
		ip_render := ""
		if ip != "" {
			ip_pre = tool.IP_preprocess(db, ip, config.IP)[0]
			ip_render = tool.IP_parser(db, ip, config.IP)
		}
		data_list = append(data_list, []string{title, sub, date, code, stop, ip_pre, ip_render, id, agree})
	}

	return map[string]any{"response": "ok", "data": data_list}
}
