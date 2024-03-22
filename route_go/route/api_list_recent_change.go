package route

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"opennamu/route/tool"
	"strconv"
)

func Api_list_recent_change(call_arg []string) {
	db_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &db_set)

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[1]), &other_set)

	db := tool.DB_connect(db_set)
	if db == nil {
		return
	}
	defer db.Close()

	set_type := other_set["set_type"]
	if set_type == "edit" {
		set_type = ""
	}

	limit_int, err := strconv.Atoi(other_set["limit"])
	if err != nil {
		return
	}

	if limit_int > 50 || limit_int < 0 {
		limit_int = 50
	}

	stmt, err := db.Prepare(tool.DB_change(db_set, "select id, title from rc where type = ? order by date desc limit ?"))
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(set_type, limit_int)
	if err != nil {
		return
	}
	defer rows.Close()

	var data_list [][]string
	admin_auth := tool.Get_user_auth(db, db_set, other_set["ip"])

	for rows.Next() {
		var id string
		var title string

		err := rows.Scan(&id, &title)
		if err != nil {
			return
		}

		var date string
		var ip string
		var send string
		var leng string
		var hide string
		var type_data string

		stmt, err := db.Prepare(tool.DB_change(db_set, "select date, ip, send, leng, hide, type from history where id = ? and title = ?"))
		if err != nil {
			return
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
				return
			}
		}

		if hide == "" || admin_auth != "" {
			data_list = append(data_list, []string{
				id,
				title,
				date,
				tool.IP_preprocess(db, db_set, ip, other_set["ip"])[0],
				send,
				leng,
				hide,
				tool.IP_parser(db, db_set, ip, other_set["ip"]),
				type_data,
			})
		} else {
			data_list = append(data_list, []string{"", "", "", "", "", "", hide, "", ""})
		}
	}

	if len(data_list) == 0 {
		fmt.Print("{}")
	} else {
		json_data, _ := json.Marshal(data_list)
		fmt.Print(string(json_data))
	}
}
