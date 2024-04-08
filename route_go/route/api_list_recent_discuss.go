package route

import (
	"database/sql"
	"encoding/json"
	"log"
	"opennamu/route/tool"
	"strconv"
)

func Api_list_recent_discuss(call_arg []string) string {
	db_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &db_set)

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[1]), &other_set)

	db := tool.DB_connect(db_set)
	defer db.Close()

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

	var stmt *sql.Stmt

	set_type := other_set["set_type"]
	if set_type == "normal" {
		stmt, err = db.Prepare(tool.DB_change(db_set, "select title, sub, date, code, stop from rd order by date desc limit ?, ?"))
	} else if set_type == "close" {
		stmt, err = db.Prepare(tool.DB_change(db_set, "select title, sub, date, code, stop from rd where stop = 'O' order by date desc limit ?, ?"))
	} else {
		stmt, err = db.Prepare(tool.DB_change(db_set, "select title, sub, date, code, stop from rd where stop != 'O' order by date desc limit ?, ?"))
	}

	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(page_int, limit_int)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var data_list [][]string
	ip_parser_temp := map[string][]string{}

	for rows.Next() {
		var title string
		var sub string
		var date string
		var code string
		var stop string

		err := rows.Scan(&title, &sub, &date, &code, &stop)
		if err != nil {
			log.Fatal(err)
		}

		stmt, err := db.Prepare(tool.DB_change(db_set, "select ip, id from topic where code = ? order by id + 0 desc limit 1"))
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		var ip string
		var id string

		err = stmt.QueryRow(code).Scan(&ip, &id)
		if err != nil {
			if err == sql.ErrNoRows {
				ip = ""
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

		data_list = append(data_list, []string{
			title,
			sub,
			date,
			code,
			stop,
			ip_pre,
			ip_render,
			id,
		})
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
			"tool":             tool.Get_language(db, db_set, "tool", false),
			"normal":           tool.Get_language(db, db_set, "normal", false),
			"close_discussion": tool.Get_language(db, db_set, "close_discussion", false),
			"open_discussion":  tool.Get_language(db, db_set, "open_discussion", false),
			"closed":           tool.Get_language(db, db_set, "closed", false),
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
