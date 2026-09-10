package route

import (
	"database/sql"
	"strconv"
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

type bbs_filter struct {
	comment_min          int
	commented            int
	tabom_min            int
	mine                 bool
	participate          bool
	comment_user         string
	comment_user_invalid bool
	tabom_user           bool
	author               string
	author_invalid       bool
	prefix               string
	tag                  string
}

func bbs_filter_number(data string) int {
	num := tool.Str_to_int(data)
	if num < 0 {
		return 0
	}

	return num
}

func bbs_filter_parse(data string) bbs_filter {
	data = strings.Trim(data, "/")
	parts := strings.Split(data, "/")
	filter := bbs_filter{}

	for i := 0; i+1 < len(parts); {
		if parts[i] == "tag" {
			filter.tag = strings.TrimSpace(strings.Join(parts[i+1:], "/"))
			break
		}

		value := bbs_filter_number(parts[i+1])
		switch parts[i] {
		case "comment":
			filter.comment_min = value
		case "commented":
			if parts[i+1] == "1" {
				filter.commented = 1
			} else if parts[i+1] == "0" {
				filter.commented = 2
			}
		case "tabom":
			filter.tabom_min = value
		case "mine":
			filter.mine = parts[i+1] == "1"
		case "participate":
			filter.participate = parts[i+1] == "1"
		case "comment_user":
			decoded, err := tool.Get_base64_decode(parts[i+1])
			if err == nil {
				filter.comment_user = decoded
			} else {
				filter.comment_user_invalid = true
			}
		case "tabom_user":
			filter.tabom_user = parts[i+1] == "1"
		case "user":
			decoded, err := tool.Get_base64_decode(parts[i+1])
			if err == nil {
				filter.author = decoded
			} else {
				filter.author_invalid = true
			}
		case "prefix":
			filter.prefix = strings.TrimSpace(parts[i+1])
		}
		i += 2
	}

	return filter
}

func bbs_filter_path(filter bbs_filter) string {
	path := []string{}
	if filter.comment_min > 0 {
		path = append(path, "comment", strconv.Itoa(filter.comment_min))
	}
	if filter.commented == 1 {
		path = append(path, "commented", "1")
	} else if filter.commented == 2 {
		path = append(path, "commented", "0")
	}
	if filter.comment_user != "" {
		path = append(path, "comment_user", tool.Base64_encode(filter.comment_user))
	}
	if filter.tabom_min > 0 {
		path = append(path, "tabom", strconv.Itoa(filter.tabom_min))
	}
	if filter.mine {
		path = append(path, "mine", "1")
	}
	if filter.participate {
		path = append(path, "participate", "1")
	}
	if filter.tabom_user {
		path = append(path, "tabom_user", "1")
	}
	if filter.author != "" {
		path = append(path, "user", tool.Base64_encode(filter.author))
	}
	if filter.prefix != "" {
		path = append(path, "prefix", tool.Url_parser(filter.prefix))
	}
	if filter.tag != "" {
		path = append(path, "tag", tool.Url_parser(filter.tag))
	}

	return strings.Join(path, "/")
}

func bbs_filter_path_data(data string) (string, string) {
	parts := strings.Split(strings.Trim(data, "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		return "1", ""
	}

	page := parts[len(parts)-1]
	if tool.Str_to_int(page) < 1 {
		page = "1"
	}

	return page, strings.Join(parts[:len(parts)-1], "/")
}

func bbs_post_comment_count_sql(row_alias string) string {
	return "coalesce((select set_data from bbs_data comment_count_data where comment_count_data.set_name = 'comment_count' and comment_count_data.set_id = " + row_alias + ".set_id and comment_count_data.set_code = " + row_alias + ".set_code limit 1), '0') + 0"
}

func bbs_post_comment_count_update(db *sql.DB, set_id string, set_code string, change int) {
	if change == 0 {
		return
	}

	comment_count := "0"
	exists := tool.QueryRow_DB(
		db,
		"select set_data from bbs_data where set_name = 'comment_count' and set_id = ? and set_code = ? limit 1",
		[]any{&comment_count},
		set_id,
		set_code,
	)

	if !exists {
		if change > 0 {
			tool.Exec_DB(
				db,
				"insert into bbs_data (set_name, set_id, set_code, set_data) values ('comment_count', ?, ?, ?)",
				set_id,
				set_code,
				strconv.Itoa(change),
			)
		}
		return
	}

	comment_count_int := tool.Str_to_int(comment_count) + change
	if comment_count_int < 0 {
		comment_count_int = 0
	}

	tool.Exec_DB(
		db,
		"update bbs_data set set_data = ? where set_name = 'comment_count' and set_id = ? and set_code = ?",
		strconv.Itoa(comment_count_int),
		set_id,
		set_code,
	)
}

func bbs_post_tabom_count_sql(row_alias string) string {
	return "coalesce((select set_data from bbs_data tabom_data where tabom_data.set_name = 'tabom_count' and tabom_data.set_id = " + row_alias + ".set_id and tabom_data.set_code = " + row_alias + ".set_code limit 1), '0') + 0"
}

func bbs_comment_set_id_sql(row_alias string, suffix string) string {
	if tool.Get_DB_type() == "mysql" {
		return "concat(" + row_alias + ".set_id, '-', " + row_alias + ".set_code, '" + suffix + "')"
	}

	return row_alias + ".set_id || '-' || " + row_alias + ".set_code || '" + suffix + "'"
}

func bbs_post_last_activity_sql(row_alias string) string {
	comment_set_id := bbs_comment_set_id_sql(row_alias, "")
	comment_set_id_nested := bbs_comment_set_id_sql(row_alias, "-%")

	return "coalesce((select nullif(last_activity_data.set_data, '') from bbs_data last_activity_data where last_activity_data.set_name = 'last_activity' and last_activity_data.set_id = " + row_alias + ".set_id and last_activity_data.set_code = " + row_alias + ".set_code limit 1), (select max(comment_date_data.set_data) from bbs_data comment_date_data where comment_date_data.set_name = 'comment_date' and (comment_date_data.set_id = " + comment_set_id + " or comment_date_data.set_id like " + comment_set_id_nested + ")), (select nullif(date_data.set_data, '') from bbs_data date_data where date_data.set_name = 'date' and date_data.set_id = " + row_alias + ".set_id and date_data.set_code = " + row_alias + ".set_code limit 1))"
}

func bbs_post_last_activity_update(db *sql.DB, set_id string, set_code string, date string) {
	if date == "" {
		return
	}

	last_activity := ""
	if tool.QueryRow_DB(db, "select set_data from bbs_data where set_name = 'last_activity' and set_id = ? and set_code = ? limit 1", []any{&last_activity}, set_id, set_code) {
		tool.Exec_DB(db, "update bbs_data set set_data = ? where set_name = 'last_activity' and set_id = ? and set_code = ?", date, set_id, set_code)
		return
	}

	tool.Exec_DB(db, "insert into bbs_data (set_name, set_id, set_code, set_data) values ('last_activity', ?, ?, ?)", set_id, set_code, date)
}

func bbs_post_last_activity_rebuild(db *sql.DB, set_id string, set_code string) {
	comment_set_id := set_id + "-" + set_code
	last_activity := ""
	tool.QueryRow_DB(
		db,
		"select coalesce(max(set_data), '') from bbs_data where set_name = 'comment_date' and (set_id = ? or set_id like ?)",
		[]any{&last_activity},
		comment_set_id,
		comment_set_id+"-%",
	)
	if last_activity == "" {
		tool.QueryRow_DB(
			db,
			"select coalesce(set_data, '') from bbs_data where set_name = 'date' and set_id = ? and set_code = ? limit 1",
			[]any{&last_activity},
			set_id,
			set_code,
		)
	}
	if last_activity == "" {
		tool.Exec_DB(db, "delete from bbs_data where set_name = 'last_activity' and set_id = ? and set_code = ?", set_id, set_code)
		return
	}
	bbs_post_last_activity_update(db, set_id, set_code, last_activity)
}

func bbs_filter_sql(filter bbs_filter, row_alias string, user_id string) (string, []any) {
	filter_sql := ""
	filter_values := []any{}

	if filter.comment_min > 0 {
		filter_sql += " and " + bbs_post_comment_count_sql(row_alias) + " >= ?"
		filter_values = append(filter_values, filter.comment_min)
	}
	if filter.commented == 1 {
		filter_sql += " and " + bbs_post_comment_count_sql(row_alias) + " >= 1"
	} else if filter.commented == 2 {
		filter_sql += " and " + bbs_post_comment_count_sql(row_alias) + " = 0"
	}
	if filter.tabom_min > 0 {
		filter_sql += " and " + bbs_post_tabom_count_sql(row_alias) + " >= ?"
		filter_values = append(filter_values, filter.tabom_min)
	}
	if filter.mine {
		filter_sql += " and exists (select 1 from bbs_data author_data where author_data.set_name = 'user_id' and author_data.set_id = " + row_alias + ".set_id and author_data.set_code = " + row_alias + ".set_code and author_data.set_data = ?)"
		filter_values = append(filter_values, user_id)
	}
	if filter.participate {
		comment_set_id := bbs_comment_set_id_sql(row_alias, "")
		comment_set_id_nested := bbs_comment_set_id_sql(row_alias, "-%")
		filter_sql += " and (exists (select 1 from bbs_data author_data where author_data.set_name = 'user_id' and author_data.set_id = " + row_alias + ".set_id and author_data.set_code = " + row_alias + ".set_code and author_data.set_data = ?) or exists (select 1 from bbs_data comment_user_data where comment_user_data.set_name = 'comment_user_id' and comment_user_data.set_data = ? and (comment_user_data.set_id = " + comment_set_id + " or comment_user_data.set_id like " + comment_set_id_nested + ")))"
		filter_values = append(filter_values, user_id, user_id)
	}
	if filter.comment_user_invalid || filter.author_invalid {
		filter_sql += " and 1 = 0"
	}
	if filter.comment_user != "" {
		comment_set_id := bbs_comment_set_id_sql(row_alias, "")
		comment_set_id_nested := bbs_comment_set_id_sql(row_alias, "-%")
		filter_sql += " and exists (select 1 from bbs_data comment_user_data where comment_user_data.set_name = 'comment_user_id' and comment_user_data.set_data = ? and (comment_user_data.set_id = " + comment_set_id + " or comment_user_data.set_id like " + comment_set_id_nested + "))"
		filter_values = append(filter_values, filter.comment_user)
	}
	if filter.tabom_user {
		filter_sql += " and exists (select 1 from bbs_data tabom_user_data where tabom_user_data.set_name = 'tabom_list' and tabom_user_data.set_id = " + row_alias + ".set_id and tabom_user_data.set_code = " + row_alias + ".set_code and tabom_user_data.set_data = ?)"
		filter_values = append(filter_values, user_id)
	}
	if filter.author != "" {
		filter_sql += " and exists (select 1 from bbs_data author_data where author_data.set_name = 'user_id' and author_data.set_id = " + row_alias + ".set_id and author_data.set_code = " + row_alias + ".set_code and author_data.set_data = ?)"
		filter_values = append(filter_values, filter.author)
	}
	if filter.prefix != "" {
		filter_sql += " and exists (select 1 from bbs_data prefix_data where prefix_data.set_name = 'prefix' and prefix_data.set_id = " + row_alias + ".set_id and prefix_data.set_code = " + row_alias + ".set_code and prefix_data.set_data = ?)"
		filter_values = append(filter_values, filter.prefix)
	}
	if filter.tag != "" {
		filter_sql += " and exists (select 1 from bbs_data tag_data where tag_data.set_name = 'tag' and tag_data.set_id = " + row_alias + ".set_id and tag_data.set_code = " + row_alias + ".set_code and tag_data.set_data = ?)"
		filter_values = append(filter_values, filter.tag)
	}

	return filter_sql, filter_values
}

func Api_bbs(config tool.Config, bbs_num string, page string, sort_type string) map[string]any {
	return api_bbs(config, bbs_num, page, sort_type, bbs_filter{})
}

func Api_bbs_filter(config tool.Config, bbs_num string, filter_data string) map[string]any {
	page, filter_path := bbs_filter_path_data(filter_data)
	return api_bbs(config, bbs_num, page, "", bbs_filter_parse(filter_path))
}

func api_bbs(config tool.Config, bbs_num string, page string, sort_type string, filter bbs_filter) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if bbs_num == "" && !tool.Check_permission(db, "bbs_main_view", config.IP) {
		return map[string]any{"response": "require auth", "data": []map[string]string{}}
	}

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
		filter_sql, filter_values := bbs_filter_sql(filter, "bbs_data", config.IP)
		query := "select set_code, set_id, '1'"
		if sort_type == "activity" {
			query += ", " + bbs_post_last_activity_sql("bbs_data")
		}
		query += " from bbs_data where set_name = 'pinned' and set_id like ?"
		values := []any{bbs_num}
		if view_sql != "" {
			query += " and " + view_sql
			values = append(values, view_values...)
		}
		query += filter_sql
		values = append(values, filter_values...)
		query += " order by set_data desc"
		rows := tool.Query_DB(db, query, values...)

		rows_arr = append(rows_arr, rows)

		if sort_type == "activity" {
			view_sql, view_values = bbs_post_view_sql(db, bbs_num, config.IP, "title")
			filter_sql, filter_values := bbs_filter_sql(filter, "title", config.IP)
			query = "select title.set_code, title.set_id, '0', " + bbs_post_last_activity_sql("title") + " from bbs_data title where title.set_name = 'title' and title.set_id like ?"
			values = []any{bbs_num}
			if view_sql != "" {
				query += " and " + view_sql
				values = append(values, view_values...)
			}
			query += filter_sql
			values = append(values, filter_values...)
			query += " order by " + bbs_post_last_activity_sql("title") + " desc, title.set_code + 0 desc limit ?, 50"
			values = append(values, num)
			rows = tool.Query_DB(db, query, values...)
		} else if sort_type == "view" {
			view_sql, view_values = bbs_post_view_sql(db, bbs_num, config.IP, "title")
			filter_sql, filter_values := bbs_filter_sql(filter, "title", config.IP)
			query = "select title.set_code, title.set_id, '0' from bbs_data title left join bbs_data view_data on view_data.set_name = 'view_count' and view_data.set_id = title.set_id and view_data.set_code = title.set_code where title.set_name = 'title' and title.set_id like ?"
			values = []any{bbs_num}
			if view_sql != "" {
				query += " and " + view_sql
				values = append(values, view_values...)
			}
			query += filter_sql
			values = append(values, filter_values...)
			query += " order by coalesce(view_data.set_data, '0') + 0 desc, title.set_code + 0 desc limit ?, 50"
			values = append(values, num)
			rows = tool.Query_DB(db, query, values...)
		} else if sort_type == "comment" {
			view_sql, view_values = bbs_post_view_sql(db, bbs_num, config.IP, "title")
			filter_sql, filter_values := bbs_filter_sql(filter, "title", config.IP)
			query = "select title.set_code, title.set_id, '0' from bbs_data title where title.set_name = 'title' and title.set_id like ?"
			values = []any{bbs_num}
			if view_sql != "" {
				query += " and " + view_sql
				values = append(values, view_values...)
			}
			query += filter_sql
			values = append(values, filter_values...)
			query += " order by " + bbs_post_comment_count_sql("title") + " desc, title.set_code + 0 desc limit ?, 50"
			values = append(values, num)
			rows = tool.Query_DB(db, query, values...)
		} else if sort_type == "tabom" {
			view_sql, view_values = bbs_post_view_sql(db, bbs_num, config.IP, "title")
			filter_sql, filter_values := bbs_filter_sql(filter, "title", config.IP)
			query = "select title.set_code, title.set_id, '0' from bbs_data title where title.set_name = 'title' and title.set_id like ?"
			values = []any{bbs_num}
			if view_sql != "" {
				query += " and " + view_sql
				values = append(values, view_values...)
			}
			query += filter_sql
			values = append(values, filter_values...)
			query += " order by " + bbs_post_tabom_count_sql("title") + " desc, title.set_code + 0 desc limit ?, 50"
			values = append(values, num)
			rows = tool.Query_DB(db, query, values...)
		} else {
			view_sql, view_values = bbs_post_view_sql(db, bbs_num, config.IP, "bbs_data")
			filter_sql, filter_values := bbs_filter_sql(filter, "bbs_data", config.IP)
			query = "select set_code, set_id, '0' from bbs_data where set_name = 'title' and set_id like ?"
			values = []any{bbs_num}
			if view_sql != "" {
				query += " and " + view_sql
				values = append(values, view_values...)
			}
			query += filter_sql
			values = append(values, filter_values...)
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
			activity_date := ""
			var err error
			if sort_type == "activity" {
				err = rows_arr[for_a].Scan(&set_code, &set_id, &pinned, &activity_date)
			} else {
				err = rows_arr[for_a].Scan(&set_code, &set_id, &pinned)
			}
			if err != nil {
				panic(err)
			}

			if bbs_num == "" && !tool.Check_acl(db, set_id, "", "bbs_view", config.IP) {
				continue
			}

			temp_data["set_code"] = set_code
			temp_data["set_id"] = set_id
			temp_data["pinned"] = pinned
			if activity_date != "" {
				temp_data["activity_date"] = activity_date
			}

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
