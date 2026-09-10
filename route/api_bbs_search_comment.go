package route

import (
	"database/sql"
	"strings"

	"opennamu/route/tool"
)

func bbs_search_comment_location(set_id string, set_code string, comment_code string) (string, string, bool) {
	parts := strings.Split(comment_code, "-")
	if !bbs_comment_code_regex.MatchString(comment_code) {
		return "", "", false
	}

	comment_set_id := set_id + "-" + set_code
	if len(parts) > 1 {
		comment_set_id += "-" + strings.Join(parts[:len(parts)-1], "-")
	}
	return comment_set_id, parts[len(parts)-1], true
}

func bbs_comment_storage_location(set_id string, set_code string) (string, string, string, bool) {
	parts := strings.Split(set_id, "-")
	root_id := ""
	root_code := ""
	comment_parts := []string{}

	if strings.HasPrefix(set_id, "-1-") {
		parts = strings.Split(strings.TrimPrefix(set_id, "-1-"), "-")
		if len(parts) < 1 || parts[0] == "" {
			return "", "", "", false
		}
		root_id = "-1"
		root_code = parts[0]
		comment_parts = parts[1:]
	} else {
		if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
			return "", "", "", false
		}
		root_id = parts[0]
		root_code = parts[1]
		comment_parts = parts[2:]
	}

	comment_code := set_code
	if len(comment_parts) > 0 {
		comment_code = strings.Join(append(comment_parts, set_code), "-")
	}
	if comment_code == "" {
		return "", "", "", false
	}

	return root_id, root_code, comment_code, true
}

func bbs_search_comment_item_data(db *sql.DB, config tool.Config, set_id string, set_code string, comment_code string, ip_parser_temp map[string][]string, auth_info map[string]bool, keyword string) (map[string]string, bool) {
	comment_manage := auth_info["bbs_comment_manage"]
	if !tool.Check_acl(db, set_id, "", "bbs_view", config.IP) && !comment_manage {
		return nil, false
	}

	data, visible := bbs_search_item_data(db, config, set_code, set_id, ip_parser_temp, auth_info, keyword, "title")
	if !visible {
		return nil, false
	}

	comment_set_id, comment_set_code, exists := bbs_search_comment_location(set_id, set_code, comment_code)
	if !exists {
		return nil, false
	}

	comment_data := ""
	comment_date := ""
	comment_user_id := ""
	blind := ""
	rows := tool.Query_DB(
		db,
		"select set_name, set_data from bbs_data where set_id = ? and set_code = ? and set_name in ('blind', 'comment', 'comment_date', 'comment_user_id')",
		comment_set_id,
		comment_set_code,
	)
	for rows.Next() {
		set_name := ""
		set_data := ""
		if rows.Scan(&set_name, &set_data) != nil {
			rows.Close()
			return nil, false
		}
		switch set_name {
		case "comment":
			comment_data = set_data
		case "comment_date":
			comment_date = set_data
		case "comment_user_id":
			comment_user_id = set_data
		case "blind":
			blind = set_data
		}
	}
	rows.Close()

	if comment_data == "" || comment_user_id == "" || (blind == "O" && !comment_manage) {
		return nil, false
	}

	if value, ok := ip_parser_temp[comment_user_id]; ok {
		data["user_id_render"] = value[1]
	} else {
		ip_render := tool.Get_user_profile_image_ui(db, comment_user_id) + tool.IP_parser(db, comment_user_id, config.IP)
		ip_parser_temp[comment_user_id] = []string{
			tool.IP_preprocess(db, comment_user_id, config.IP)[0],
			ip_render,
		}
		data["user_id_render"] = ip_render
	}

	data["comment_code"] = comment_code
	data["date"] = comment_date
	data["search_snippet_html"] = search_snippet(comment_data, keyword)
	return data, true
}

func bbs_search_comment_index_data(db *sql.DB, config tool.Config, keyword string, set_id string, page int) ([]map[string]string, bool) {
	target_count := page * 50
	candidate_limit := 500
	if target_count < candidate_limit {
		candidate_limit = target_count
	}
	if candidate_limit < 50 {
		candidate_limit = 50
	}

	auth_info := tool.Get_auth_info(db, config.IP)
	ip_parser_temp := map[string][]string{}
	data_list := []map[string]string{}
	candidate_offset := 0
	max_candidate := target_count + 1000

	for candidate_offset < max_candidate && len(data_list) < target_count {
		candidate_list, ok := tool.Search_bbs_index_search_comment(keyword, set_id, candidate_offset, candidate_limit)
		if !ok {
			return nil, false
		}
		if len(candidate_list) == 0 {
			break
		}
		candidate_offset += len(candidate_list)

		for _, key := range candidate_list {
			comment_set_id, set_code, comment_code, valid := tool.Search_bbs_index_comment_key_data(key)
			if !valid || (set_id == "" && comment_set_id == "0") {
				continue
			}
			item_data, visible := bbs_search_comment_item_data(
				db,
				config,
				comment_set_id,
				set_code,
				comment_code,
				ip_parser_temp,
				auth_info,
				keyword,
			)
			if visible {
				data_list = append(data_list, item_data)
			}
			if len(data_list) >= target_count {
				break
			}
		}

		if len(candidate_list) < candidate_limit {
			break
		}
	}

	start := (page - 1) * 50
	if start >= len(data_list) {
		return []map[string]string{}, true
	}
	end := start + 50
	if end > len(data_list) {
		end = len(data_list)
	}
	return data_list[start:end], true
}

func bbs_search_comment_sql_data(db *sql.DB, config tool.Config, keyword string, set_id string, page int) []map[string]string {
	where_data := "b.set_name = 'comment' and b.set_data like ?"
	values := []any{"%" + keyword + "%"}
	if set_id != "" {
		where_data += " and b.set_id like ?"
		values = append(values, set_id+"-%")
	}

	if !tool.Check_permission(db, "bbs_comment_manage", config.IP) {
		where_data += " and not exists (select 1 from bbs_data blind_data where blind_data.set_name = 'blind' and blind_data.set_data = 'O' and blind_data.set_id = b.set_id and blind_data.set_code = b.set_code)"
	}

	offset := (page - 1) * 50
	values = append(values, offset)
	rows := tool.Query_DB(
		db,
		"select b.set_id, b.set_code from bbs_data b left join bbs_data d on d.set_id = b.set_id and d.set_code = b.set_code and d.set_name = 'comment_date' where "+where_data+" order by d.set_data desc limit ?, 50",
		values...,
	)
	defer rows.Close()

	auth_info := tool.Get_auth_info(db, config.IP)
	ip_parser_temp := map[string][]string{}
	data_list := []map[string]string{}
	for rows.Next() {
		comment_set_id := ""
		comment_set_code := ""
		if rows.Scan(&comment_set_id, &comment_set_code) != nil {
			continue
		}

		post_set_id, post_set_code, comment_code, valid := bbs_comment_storage_location(comment_set_id, comment_set_code)
		if !valid || (set_id == "" && post_set_id == "0") {
			continue
		}

		item_data, visible := bbs_search_comment_item_data(
			db,
			config,
			post_set_id,
			post_set_code,
			comment_code,
			ip_parser_temp,
			auth_info,
			keyword,
		)
		if visible {
			data_list = append(data_list, item_data)
		}
	}
	return data_list
}

func Api_bbs_search_comment(config tool.Config, keyword string, set_id string, page string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if set_id == "" && !tool.Check_permission(db, "bbs_main_view", config.IP) {
		return map[string]any{"response": "require auth", "data": []map[string]string{}}
	}
	if set_id != "" && !tool.Check_acl(db, set_id, "", "bbs_view", config.IP) {
		return map[string]any{"response": "require auth", "data": []map[string]string{}}
	}

	keyword = strings.TrimSpace(keyword)
	data_list := []map[string]string{}
	if keyword != "" {
		page_num := tool.Str_to_int(page)
		if page_num < 1 {
			page_num = 1
		}

		if indexed_data, ok := bbs_search_comment_index_data(db, config, keyword, set_id, page_num); ok {
			data_list = indexed_data
		} else {
			data_list = bbs_search_comment_sql_data(db, config, keyword, set_id, page_num)
		}
	}

	return map[string]any{
		"response": "ok",
		"data":     data_list,
	}
}
