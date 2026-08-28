package route

import (
	"strings"

	"opennamu/route/tool"
)

func Api_bbs_search(config tool.Config, keyword string, set_id string, page string) map[string]any {
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
		values = append(values, "%"+keyword+"%", offset)

		rows := tool.Query_DB(
			db,
			"select b.set_code, b.set_id from bbs_data b left join bbs_data d on d.set_code = b.set_code and d.set_id = b.set_id and d.set_name = 'date' where "+where_data+" and b.set_name in ('title', 'prefix', 'tag') and b.set_data like ? group by b.set_code, b.set_id order by max(d.set_data) desc limit ?, 50",
			values...,
		)
		defer rows.Close()

		ip_parser_temp := map[string][]string{}
		for rows.Next() {
			var set_code_data string
			var set_id_data string
			if err := rows.Scan(&set_code_data, &set_id_data); err != nil {
				panic(err)
			}

			if !tool.Check_acl(db, set_id_data, "", "bbs_view", config.IP) {
				continue
			}

			temp_data := map[string]string{
				"set_code": set_code_data,
				"set_id":   set_id_data,
				"pinned":   "0",
			}

			item_rows := tool.Query_DB(
				db,
				"select set_name, set_data from bbs_data where set_code = ? and set_id = ?",
				set_code_data,
				set_id_data,
			)

			for item_rows.Next() {
				var set_name string
				var set_data string
				if err := item_rows.Scan(&set_name, &set_data); err != nil {
					item_rows.Close()
					panic(err)
				}

				if set_name == "user_id" {
					ip_pre := ""
					ip_render := ""
					if value, ok := ip_parser_temp[set_data]; ok {
						ip_pre = value[0]
						ip_render = value[1]
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

				if set_name == "pinned" {
					temp_data["pinned"] = "1"
					continue
				}

				if set_name != "data" {
					temp_data[set_name] = set_data
				}
			}
			item_rows.Close()
			data_list = append(data_list, temp_data)
		}
	}

	return_data := make(map[string]any)
	return_data["response"] = "ok"
	return_data["data"] = data_list

	return return_data
}
