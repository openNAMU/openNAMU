package route

import (
	"database/sql"
	"encoding/json"
	"log"
	"opennamu/route/tool"
	"strconv"
)

func Api_list_recent_change(call_arg []string) string {
	db_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &db_set)

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[1]), &other_set)

	db := tool.DB_connect(db_set)
	defer db.Close()

	set_type := other_set["set_type"]
	if set_type == "edit" {
		set_type = ""
	}

	limit_int, err := strconv.Atoi(other_set["limit"])
	if err != nil {
		log.Fatal(err)
	}

	if limit_int > 50 || limit_int < 0 {
		limit_int = 50
	}

	page_int, err := strconv.Atoi(other_set["num"])
	if err != nil {
		log.Fatal(err)
	}

	if page_int > 0 {
		page_int = (page_int * limit_int) - limit_int
	} else {
		page_int = 0
	}

	stmt, err := db.Prepare(tool.DB_change(db_set, "select id, title from rc where type = ? order by date desc limit ?, ?"))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(set_type, page_int, limit_int)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var data_list [][]string

	admin_auth := tool.Get_user_auth(db, db_set, other_set["ip"])
	ip_parser_temp := map[string][]string{}

	for rows.Next() {
		var id string
		var title string

		err := rows.Scan(&id, &title)
		if err != nil {
			log.Fatal(err)
		}

		var date string
		var ip string
		var send string
		var leng string
		var hide string
		var type_data string

		stmt, err := db.Prepare(tool.DB_change(db_set, "select date, ip, send, leng, hide, type from history where id = ? and title = ?"))
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		err = stmt.QueryRow(id, title).Scan(&date, &ip, &send, &leng, &hide, &type_data)
		if err != nil {
			if err == sql.ErrNoRows {
				date = ""
				ip = ""
				send = ""
				leng = ""
				hide = ""
				type_data = ""
			} else {
				log.Fatal(err)
			}
		}

		var ip_pre string
		var ip_render string

		if _, ok := ip_parser_temp[ip]; ok {
			ip_pre = ip_parser_temp[ip][0]
			ip_render = ip_parser_temp[ip][1]
		} else {
			ip_pre = tool.IP_preprocess(db, db_set, ip, other_set["ip"])[0]
			ip_render = tool.IP_parser(db, db_set, ip, other_set["ip"])

			ip_parser_temp[ip] = []string{ip_pre, ip_render}
		}

		if hide == "" || admin_auth != "" {
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

	if other_set["legacy"] != "" {
		if len(data_list) == 0 {
			return "{}"
		} else {
			json_data, _ := json.Marshal(data_list)
			return string(json_data)
		}
	} else {
		auth_name := tool.Get_user_auth(db, db_set, other_set["ip"])
		auth_info := tool.Get_auth_group_info(db, db_set, auth_name)

		return_data := make(map[string]interface{})
		return_data["language"] = map[string]string{
			"tool":           tool.Get_language(db, db_set, "tool", false),
			"normal":         tool.Get_language(db, db_set, "normal", false),
			"edit":           tool.Get_language(db, db_set, "edit", false),
			"move":           tool.Get_language(db, db_set, "move", false),
			"delete":         tool.Get_language(db, db_set, "delete", false),
			"revert":         tool.Get_language(db, db_set, "revert", false),
			"new_doc":        tool.Get_language(db, db_set, "new_doc", false),
			"edit_request":   tool.Get_language(db, db_set, "edit_request", false),
			"user_document":  tool.Get_language(db, db_set, "user_document", false),
			"raw":            tool.Get_language(db, db_set, "raw", false),
			"compare":        tool.Get_language(db, db_set, "compare", false),
			"history":        tool.Get_language(db, db_set, "history", false),
			"hide":           tool.Get_language(db, db_set, "hide", false),
			"history_delete": tool.Get_language(db, db_set, "history_delete", false),
			"send_edit":      tool.Get_language(db, db_set, "send_edit", false),
			"file":           tool.Get_language(db, db_set, "file", false),
			"category":       tool.Get_language(db, db_set, "category", false),
			"setting":        tool.Get_language(db, db_set, "setting", false),
		}
		return_data["auth"] = auth_info

		if len(data_list) == 0 {
			return_data["data"] = map[string]string{}
		} else {
			return_data["data"] = data_list
		}

		json_data, _ := json.Marshal(return_data)
		return string(json_data)
	}
}
