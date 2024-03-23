package route

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"opennamu/route/tool"
	"strconv"
)

func Api_w_watch_list(call_arg []string) {
	db_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &db_set)

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[1]), &other_set)

	page, _ := strconv.Atoi(other_set["num"])
	num := 0
	if page*50 > 0 {
		num = page*50 - 50
	}

	db := tool.DB_connect(db_set)
	if db == nil {
		return
	}
	defer db.Close()

	ip := other_set["ip"]
	if tool.Get_user_auth(db, db_set, ip) == "" {
		fmt.Print("{}")
		return
	}

	var stmt *sql.Stmt
	var err error
	if other_set["do_type"] == "star_doc" {
		stmt, err = db.Prepare(tool.DB_change(db_set, "select id from user_set where name = 'star_doc' and data = ? limit ?, 50"))
	} else {
		stmt, err = db.Prepare(tool.DB_change(db_set, "select id from user_set where name = 'watchlist' and data = ? limit ?, 50"))
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

	var data_list []string

	for rows.Next() {
		var user_name string

		err := rows.Scan(&user_name)
		if err != nil {
			log.Fatal(err)
		}

		data_list = append(data_list, user_name)
	}

	if len(data_list) == 0 {
		fmt.Print("{}")
	} else {
		json_data, _ := json.Marshal(data_list)
		fmt.Print(string(json_data))
	}
}
