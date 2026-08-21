package route

import (
	"strings"

	"opennamu/route/tool"
)

func Api_bbs_search(config tool.Config, keyword string, set_id string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	data_list := []map[string]string{}
	keyword = strings.TrimSpace(keyword)
	if keyword != "" {
		rows := tool.Query_DB(
			db,
			"select set_code, set_id from bbs_data where set_id = ? and set_name in ('title', 'prefix', 'tag') and set_data like ? group by set_code, set_id order by set_code + 0 desc limit 50",
			set_id,
			"%"+keyword+"%",
		)
		defer rows.Close()

		ip_parser_temp := map[string][]string{}
		for rows.Next() {
			var set_code_data string
			var set_id_data string
			if err := rows.Scan(&set_code_data, &set_id_data); err != nil {
				panic(err)
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
