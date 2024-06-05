package route

import (
	"database/sql"
	"encoding/json"
	"log"
	"opennamu/route/tool"
	"strconv"
)

func Api_w_watch_list(call_arg []string) string {
	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &other_set)

	page, _ := strconv.Atoi(other_set["num"])
	num := 0
	if page*50 > 0 {
		num = page*50 - 50
	}

	db := tool.DB_connect()
	defer db.Close()

	if tool.Get_user_auth(db, other_set["ip"]) == "" {
		return "{}"
	}

	var stmt *sql.Stmt
	var err error
	if other_set["do_type"] == "star_doc" {
		stmt, err = db.Prepare(tool.DB_change("select id from user_set where name = 'star_doc' and data = ? limit ?, 50"))
	} else {
		stmt, err = db.Prepare(tool.DB_change("select id from user_set where name = 'watchlist' and data = ? limit ?, 50"))
	}
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(other_set["name"], num)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var data_list [][]string
	ip_parser_temp := map[string][]string{}

	for rows.Next() {
		var user_name string

		err := rows.Scan(&user_name)
		if err != nil {
			log.Fatal(err)
		}

		var ip_pre string
		var ip_render string

		if _, ok := ip_parser_temp[user_name]; ok {
			ip_pre = ip_parser_temp[user_name][0]
			ip_render = ip_parser_temp[user_name][1]
		} else {
			ip_pre = tool.IP_preprocess(db, user_name, other_set["ip"])[0]
			ip_render = tool.IP_parser(db, user_name, other_set["ip"])

			ip_parser_temp[user_name] = []string{ip_pre, ip_render}
		}

		data_list = append(data_list, []string{ip_pre, ip_render})
	}

	return_data := make(map[string]interface{})
	return_data["language"] = map[string]string{
		"watchlist": tool.Get_language(db, "watchlist", false),
		"star_doc":  tool.Get_language(db, "star_doc", false),
	}

	if len(data_list) == 0 {
		return_data["data"] = map[string]string{}
	} else {
		return_data["data"] = data_list
	}

	json_data, _ := json.Marshal(return_data)
	return string(json_data)
}
