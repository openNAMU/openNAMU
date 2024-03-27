package route

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"opennamu/route/tool"
	"strconv"
)

func Api_list_recent_discuss(call_arg []string) {
	db_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &db_set)

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[1]), &other_set)

	db := tool.DB_connect(db_set)
	if db == nil {
		return
	}
	defer db.Close()

	limit_int, err := strconv.Atoi(other_set["limit"])
	if err != nil {
		log.Fatal(err)
	}

	if limit_int > 50 || limit_int < 0 {
		limit_int = 50
	}

	var stmt *sql.Stmt

	set_type := other_set["set_type"]
	if set_type == "normal" {
		stmt, err = db.Prepare(tool.DB_change(db_set, "select title, sub, date, code, stop from rd order by date desc limit ?"))
	} else if set_type == "close" {
		stmt, err = db.Prepare(tool.DB_change(db_set, "select title, sub, date, code, stop from rd where stop = 'O' order by date desc limit ?"))
	} else {
		stmt, err = db.Prepare(tool.DB_change(db_set, "select title, sub, date, code, stop from rd where stop != 'O' order by date desc limit ?"))
	}

	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(limit_int)
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

	if len(data_list) == 0 {
		fmt.Print("{}")
	} else {
		json_data, _ := json.Marshal(data_list)
		fmt.Print(string(json_data))
	}
}
