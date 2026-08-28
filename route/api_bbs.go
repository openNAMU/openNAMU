package route

import (
	"database/sql"
	"strings"

	"opennamu/route/tool"
)

func bbs_post_view_allowed(db *sql.DB, set_id string, user_id string, ip string, auth_info map[string]bool) bool {
	acl_data := bbs_set_value(db, set_id, "bbs_only_my_data_view_acl")
	if acl_data == "" || acl_data == "normal" || user_id == ip {
		return true
	}
	if auth_info == nil {
		auth_info = tool.Get_auth_info(db, ip)
	}
	if auth_info["bbs"] {
		return true
	}
	return tool.Check_acl_group(db, acl_data, auth_info)
}

func bbs_post_view_sql(db *sql.DB, set_id string, ip string, row_alias string) (string, []any) {
	auth_info := tool.Get_auth_info(db, ip)
	if set_id != "" {
		acl_data := bbs_set_value(db, set_id, "bbs_only_my_data_view_acl")
		if acl_data == "" || acl_data == "normal" || auth_info["bbs"] || tool.Check_acl_group(db, acl_data, auth_info) {
			return "", nil
		}

		return "exists (select 1 from bbs_data author where author.set_name = 'user_id' and author.set_id = " + row_alias + ".set_id and author.set_code = " + row_alias + ".set_code and author.set_data = ?)", []any{ip}
	}

	rows := tool.Query_DB(
		db,
		"select set_id, set_data from bbs_set where set_name = 'bbs_only_my_data_view_acl' and set_code = ''",
	)
	defer rows.Close()

	allowed_set_id := []string{}
	private_count := 0
	for rows.Next() {
		setting_set_id := ""
		acl_data := ""
		if rows.Scan(&setting_set_id, &acl_data) != nil || acl_data == "" || acl_data == "normal" {
			continue
		}

		private_count += 1
		if auth_info["bbs"] || tool.Check_acl_group(db, acl_data, auth_info) {
			if !tool.Arr_in_str(allowed_set_id, setting_set_id) {
				allowed_set_id = append(allowed_set_id, setting_set_id)
			}
		}
	}

	if private_count == 0 {
		return "", nil
	}

	view_sql := "(not exists (select 1 from bbs_set only_view where only_view.set_name = 'bbs_only_my_data_view_acl' and only_view.set_code = '' and only_view.set_id = " + row_alias + ".set_id and only_view.set_data != '' and only_view.set_data != 'normal')"
	view_values := []any{}
	if len(allowed_set_id) > 0 {
		placeholders := strings.Repeat("?,", len(allowed_set_id))
		placeholders = strings.TrimSuffix(placeholders, ",")
		view_sql += " or " + row_alias + ".set_id in (" + placeholders + ")"
		for _, allowed_id := range allowed_set_id {
			view_values = append(view_values, allowed_id)
		}
	}
	view_sql += " or exists (select 1 from bbs_data author where author.set_name = 'user_id' and author.set_id = " + row_alias + ".set_id and author.set_code = " + row_alias + ".set_code and author.set_data = ?))"
	view_values = append(view_values, ip)

	return view_sql, view_values
}

func Api_bbs(config tool.Config, bbs_num string, page string, sort_type string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if bbs_num != "" && !tool.Check_acl(db, bbs_num, "", "bbs_view", config.IP) {
		return map[string]any{"response": "require auth", "data": []map[string]string{}}
	}

	rows_arr := []*sql.Rows{}
	if bbs_num == "" {
		view_sql, view_values := bbs_post_view_sql(db, bbs_num, config.IP, "bbs_data")
		query := "select set_code, set_id, '0' from bbs_data where set_name = 'date' and " + tool.Get_except_set_id_SQL()
		if view_sql != "" {
			query += " and " + view_sql
		}
		query += " order by set_data desc limit 50"
		rows := tool.Query_DB(db, query, view_values...)

		rows_arr = append(rows_arr, rows)
	} else {
		page := tool.Str_to_int(page)
		num := 0
		if page*50 > 0 {
			num = page*50 - 50
		}

		view_sql, view_values := bbs_post_view_sql(db, bbs_num, config.IP, "bbs_data")
		query := "select set_code, set_id, '1' from bbs_data where set_name = 'pinned' and set_id like ?"
		values := []any{bbs_num}
		if view_sql != "" {
			query += " and " + view_sql
			values = append(values, view_values...)
		}
		query += " order by set_data desc"
		rows := tool.Query_DB(db, query, values...)

		rows_arr = append(rows_arr, rows)

		if sort_type == "view" {
			view_sql, view_values = bbs_post_view_sql(db, bbs_num, config.IP, "title")
			query = "select title.set_code, title.set_id, '0' from bbs_data title left join bbs_data view_data on view_data.set_name = 'view_count' and view_data.set_id = title.set_id and view_data.set_code = title.set_code where title.set_name = 'title' and title.set_id like ?"
			values = []any{bbs_num}
			if view_sql != "" {
				query += " and " + view_sql
				values = append(values, view_values...)
			}
			query += " order by coalesce(view_data.set_data, '0') + 0 desc, title.set_code + 0 desc limit ?, 50"
			values = append(values, num)
			rows = tool.Query_DB(db, query, values...)
		} else {
			view_sql, view_values = bbs_post_view_sql(db, bbs_num, config.IP, "bbs_data")
			query = "select set_code, set_id, '0' from bbs_data where set_name = 'title' and set_id like ?"
			values = []any{bbs_num}
			if view_sql != "" {
				query += " and " + view_sql
				values = append(values, view_values...)
			}
			query += " order by set_code + 0 desc limit ?, 50"
			values = append(values, num)
			rows = tool.Query_DB(db, query, values...)
		}

		rows_arr = append(rows_arr, rows)
	}

	data_list := []map[string]string{}
	ip_parser_temp := map[string][]string{}

	for for_a := 0; for_a < len(rows_arr); for_a++ {
		defer rows_arr[for_a].Close()

		for rows_arr[for_a].Next() {
			temp_data := make(map[string]string)

			var set_code string
			var set_id string
			var pinned string

			err := rows_arr[for_a].Scan(&set_code, &set_id, &pinned)
			if err != nil {
				panic(err)
			}

			temp_data["set_code"] = set_code
			temp_data["set_id"] = set_id
			temp_data["pinned"] = pinned

			rows := tool.Query_DB(
				db,
				"select set_name, set_data, set_code, set_id from bbs_data where set_code = ? and set_id = ?",
				set_code,
				set_id,
			)
			defer rows.Close()

			for rows.Next() {
				var set_name string
				var set_data string

				err := rows.Scan(&set_name, &set_data, &set_code, &set_id)
				if err != nil {
					panic(err)
				}

				if set_name == "user_id" {
					var ip_pre string
					var ip_render string

					if _, ok := ip_parser_temp[set_data]; ok {
						ip_pre = ip_parser_temp[set_data][0]
						ip_render = ip_parser_temp[set_data][1]
					} else {
						ip_pre = tool.IP_preprocess(db, set_data, config.IP)[0]
						ip_render = tool.IP_parser(db, set_data, config.IP)

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

				if set_name != "data" && set_name != "pinned" {
					temp_data[set_name] = set_data
				}
			}
			data_list = append(data_list, temp_data)
		}
	}

	return_data := make(map[string]any)
	return_data["response"] = "ok"
	return_data["data"] = data_list

	return return_data
}
