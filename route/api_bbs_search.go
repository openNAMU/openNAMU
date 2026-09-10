package route

import (
	"database/sql"
	"strings"

	"opennamu/route/tool"
)

func bbs_search_item_data(db *sql.DB, config tool.Config, set_code string, set_id string, ip_parser_temp map[string][]string, auth_info map[string]bool, keyword string, search_type string) (map[string]string, bool) {
	if !tool.Check_acl(db, set_id, "", "bbs_view", config.IP) {
		return nil, false
	}

	temp_data := map[string]string{
		"set_code": set_code,
		"set_id":   set_id,
		"pinned":   "0",
	}
	user_id := ""
	content_data := ""
	item_rows := tool.Query_DB(
		db,
		"select set_name, set_data from bbs_data where set_code = ? and set_id = ?",
		set_code,
		set_id,
	)
	for item_rows.Next() {
		set_name := ""
		set_data := ""
		if err := item_rows.Scan(&set_name, &set_data); err != nil {
			item_rows.Close()
			return nil, false
		}

		if set_name == "user_id" {
			user_id = set_data
			ip_pre := ""
			ip_render := ""
			if value, ok := ip_parser_temp[set_data]; ok {
				ip_pre = value[0]
				ip_render = value[1]
			} else {
				ip_pre = tool.IP_preprocess(db, set_data, config.IP)[0]
				ip_render = tool.Get_user_profile_image_ui(db, set_data) + tool.IP_parser(db, set_data, config.IP)
				ip_parser_temp[set_data] = []string{ip_pre, ip_render}
			}
			set_data = ip_pre
			temp_data["user_id_render"] = ip_render
		}

		if set_name == "tag" {
			if temp_data["tags"] != "" {
				temp_data["tags"] += ", "
			}
			temp_data["tags"] += set_data
			continue
		}

		if set_name == "pinned" {
			temp_data["pinned"] = "1"
			continue
		}

		if set_name == "data" {
			content_data = set_data
			continue
		}

		temp_data[set_name] = set_data
	}
	item_rows.Close()

	if !bbs_post_view_allowed(db, set_id, set_code, user_id, config.IP, auth_info) {
		return nil, false
	}
	temp_data["title_html"] = search_highlight(temp_data["title"], keyword)
	temp_data["prefix_html"] = search_highlight(temp_data["prefix"], keyword)
	temp_data["tags_html"] = bbs_tags_html(temp_data["set_id"], temp_data["tags"], keyword)
	if search_type == "data" {
		temp_data["search_snippet_html"] = search_snippet(content_data, keyword)
	}
	return temp_data, true
}

func bbs_search_index_data(db *sql.DB, config tool.Config, keyword string, set_id string, page int, search_type string) ([]map[string]string, bool) {
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
		var candidate_list []string
		var ok bool
		if search_type == "data" {
			candidate_list, ok = tool.Search_bbs_index_search_data(keyword, set_id, candidate_offset, candidate_limit)
		} else {
			candidate_list, ok = tool.Search_bbs_index_search(keyword, set_id, candidate_offset, candidate_limit)
		}
		if !ok {
			return nil, false
		}
		if len(candidate_list) == 0 {
			break
		}
		candidate_offset += len(candidate_list)

		for _, key := range candidate_list {
			set_id_data, set_code_data, valid := tool.Search_bbs_index_key_data(key)
			if !valid || (set_id == "" && set_id_data == "0") {
				continue
			}
			item_data, visible := bbs_search_item_data(db, config, set_code_data, set_id_data, ip_parser_temp, auth_info, keyword, search_type)
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

func Api_bbs_search(config tool.Config, keyword string, set_id string, page string) map[string]any {
	return api_bbs_search(config, keyword, set_id, page, "title")
}

func Api_bbs_search_data(config tool.Config, keyword string, set_id string, page string) map[string]any {
	return api_bbs_search(config, keyword, set_id, page, "data")
}

func api_bbs_search(config tool.Config, keyword string, set_id string, page string, search_type string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if set_id == "" && !tool.Check_permission(db, "bbs_main_view", config.IP) {
		return map[string]any{"response": "require auth", "data": []map[string]string{}}
	}

	if set_id != "" && !tool.Check_acl(db, set_id, "", "bbs_view", config.IP) {
		return map[string]any{"response": "require auth", "data": []map[string]string{}}
	}

	data_list := []map[string]string{}
	keyword = strings.TrimSpace(keyword)
	if keyword != "" {
		page_num := tool.Str_to_int(page)
		if page_num < 1 {
			page_num = 1
		}
		offset := (page_num - 1) * 50
		if data_list, ok := bbs_search_index_data(db, config, keyword, set_id, page_num, search_type); ok {
			return map[string]any{
				"response": "ok",
				"data":     data_list,
			}
		}

		where_data := "b.set_id = ?"
		values := []any{set_id}
		if set_id == "" {
			where_data = "not b.set_id = \"0\""
			values = []any{}
		}
		view_sql, view_values := bbs_post_view_sql(db, set_id, config.IP, "b")
		if view_sql != "" {
			where_data += " and " + view_sql
			values = append(values, view_values...)
		}
		search_set_name := "'title', 'prefix', 'tag'"
		if search_type == "data" {
			search_set_name = "'data'"
		}
		values = append(values, "%"+keyword+"%", offset)

		rows := tool.Query_DB(
			db,
			"select b.set_code, b.set_id from bbs_data b left join bbs_data d on d.set_code = b.set_code and d.set_id = b.set_id and d.set_name = 'date' where "+where_data+" and b.set_name in ("+search_set_name+") and b.set_data like ? group by b.set_code, b.set_id order by max(d.set_data) desc limit ?, 50",
			values...,
		)
		defer rows.Close()

		ip_parser_temp := map[string][]string{}
		auth_info := tool.Get_auth_info(db, config.IP)
		for rows.Next() {
			var set_code_data string
			var set_id_data string
			if err := rows.Scan(&set_code_data, &set_id_data); err != nil {
				panic(err)
			}

			if !tool.Check_acl(db, set_id_data, "", "bbs_view", config.IP) {
				continue
			}

			item_data, visible := bbs_search_item_data(db, config, set_code_data, set_id_data, ip_parser_temp, auth_info, keyword, search_type)
			if visible {
				data_list = append(data_list, item_data)
			}
		}
	}

	return_data := make(map[string]any)
	return_data["response"] = "ok"
	return_data["data"] = data_list

	return return_data
}
