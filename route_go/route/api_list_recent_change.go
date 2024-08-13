package route

import (
	"database/sql"
	"log"
	"opennamu/route/tool"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

func Api_list_recent_change(call_arg []string) string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &other_set)

	db := tool.DB_connect()
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

	stmt, err := db.Prepare(tool.DB_change("select id, title from rc where type = ? order by date desc limit ?, ?"))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(set_type, page_int, limit_int)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	data_list := [][]string{}

	admin_auth := tool.Check_acl(db, "", "", "hidel_auth", other_set["ip"])
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

		stmt, err := db.Prepare(tool.DB_change("select date, ip, send, leng, hide, type from history where id = ? and title = ?"))
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
			ip_pre = tool.IP_preprocess(db, ip, other_set["ip"])[0]
			ip_render = tool.IP_parser(db, ip, other_set["ip"])

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

	if other_set["legacy"] != "" {
		json_data, _ := json.Marshal(data_list)
		return string(json_data)
	} else {
		auth_name := tool.Get_user_auth(db, other_set["ip"])
		auth_info := tool.Get_auth_group_info(db, auth_name)

		return_data := make(map[string]interface{})
		return_data["language"] = map[string]string{
			"tool":           tool.Get_language(db, "tool", false),
			"normal":         tool.Get_language(db, "normal", false),
			"edit":           tool.Get_language(db, "edit", false),
			"move":           tool.Get_language(db, "move", false),
			"delete":         tool.Get_language(db, "delete", false),
			"revert":         tool.Get_language(db, "revert", false),
			"new_doc":        tool.Get_language(db, "new_doc", false),
			"edit_request":   tool.Get_language(db, "edit_request", false),
			"user_document":  tool.Get_language(db, "user_document", false),
			"raw":            tool.Get_language(db, "raw", false),
			"compare":        tool.Get_language(db, "compare", false),
			"history":        tool.Get_language(db, "history", false),
			"hide":           tool.Get_language(db, "hide", false),
			"history_delete": tool.Get_language(db, "history_delete", false),
			"send_edit":      tool.Get_language(db, "send_edit", false),
			"file":           tool.Get_language(db, "file", false),
			"category":       tool.Get_language(db, "category", false),
			"setting":        tool.Get_language(db, "setting", false),
			"remove_hidden":  tool.Get_language(db, "remove_hidden", false),
			"admin_tool":     tool.Get_language(db, "admin_tool", false),
		}
		return_data["auth"] = auth_info
		return_data["data"] = data_list

		json_data, _ := json.Marshal(return_data)
		return string(json_data)
	}
}
