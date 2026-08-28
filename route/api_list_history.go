package route

import (
	"database/sql"
	"opennamu/route/tool"
)

func api_history_language(db *sql.DB) map[string]string {
	language_name_list := []string{
		"normal",
		"edit",
		"move",
		"delete",
		"revert",
		"r1",
		"setting",
		"user",
		"user_document",
		"file",
		"category",
		"view",
		"raw",
		"compare",
		"history",
		"hide",
		"history_delete",
		"send_edit",
	}
	language_data := map[string]string{}
	for _, name := range language_name_list {
		language_data[name] = tool.Get_language(db, name, false)
	}
	return language_data
}

func api_history_auth(db *sql.DB, ip string) map[string]bool {
	auth_info := tool.Get_auth_info(db, ip)
	return map[string]bool{
		"owner": auth_info["owner"],
		"hidel": auth_info["hidel"],
	}
}

func Api_list_history(config tool.Config, doc_name string, set_type string, num string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, doc_name, "", "render", config.IP) {
		return map[string]any{"response": "require auth", "data": [][]string{}}
	}

	page_int := tool.Str_to_int(num)
	if page_int > 0 {
		page_int = (page_int * 50) - 50
	} else {
		page_int = 0
	}

	var rows *sql.Rows

	if set_type == "edit" {
		set_type = ""
	}

	if set_type == "normal" {
		rows = tool.Query_DB(
			db,
			"select id, title, date, ip, send, leng, hide, type from history where title = ? order by id + 0 desc limit ?, 50",
			doc_name, page_int,
		)
	} else {
		rows = tool.Query_DB(
			db,
			"select id, title, date, ip, send, leng, hide, type from history where title = ? and type = ? order by id + 0 desc limit ?, 50",
			doc_name, set_type, page_int,
		)
	}
	defer rows.Close()

	data_list := [][]string{}

	admin_auth := tool.Check_acl(db, "", "", "hidel_auth", config.IP)
	ip_parser_temp := map[string][]string{}

	for rows.Next() {
		var id string
		var title string
		var date string
		var ip string
		var send string
		var leng string
		var hide string
		var type_data string

		err := rows.Scan(&id, &title, &date, &ip, &send, &leng, &hide, &type_data)
		if err != nil {
			panic(err)
		}

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
