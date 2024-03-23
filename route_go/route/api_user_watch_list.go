package route

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"opennamu/route/tool"
	"strconv"
)

func Api_user_watch_list(call_arg []string) {
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
	name := other_set["name"]
	if ip != name && tool.Get_user_auth(db, db_set, ip) == "" {
		fmt.Print("{}")
		return
	}

	var stmt *sql.Stmt
	var err error
	if other_set["do_type"] == "star_doc" {
		stmt, err = db.Prepare(tool.DB_change(db_set, "select data from user_set where name = 'star_doc' and id = ? limit ?, 50"))
	} else {
		stmt, err = db.Prepare(tool.DB_change(db_set, "select data from user_set where name = 'watchlist' and id = ? limit ?, 50"))
	}
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(name, num)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var data_list []string

	for rows.Next() {
		var title_data string

		err := rows.Scan(&title_data)
		if err != nil {
			log.Fatal(err)
		}

		data_list = append(data_list, title_data)
	}

	if len(data_list) == 0 {
		fmt.Print("{}")
	} else {
		json_data, _ := json.Marshal(data_list)
		fmt.Print(string(json_data))
	}
}
