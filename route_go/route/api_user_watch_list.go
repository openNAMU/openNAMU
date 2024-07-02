package route

import (
	"database/sql"
	"log"
	"opennamu/route/tool"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

func Api_user_watch_list(call_arg []string) string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &other_set)

	db := tool.DB_connect()
	defer db.Close()

	page, _ := strconv.Atoi(other_set["num"])
	num := 0
	if page*50 > 0 {
		num = page*50 - 50
	}

	ip := other_set["ip"]
	name := other_set["name"]
	if ip != name && tool.Get_user_auth(db, ip) == "" {
		return "{}"
	}

	var stmt *sql.Stmt
	var err error
	if other_set["do_type"] == "star_doc" {
		stmt, err = db.Prepare(tool.DB_change("select data from user_set where name = 'star_doc' and id = ? limit ?, 50"))
	} else {
		stmt, err = db.Prepare(tool.DB_change("select data from user_set where name = 'watchlist' and id = ? limit ?, 50"))
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

	data_list := []string{}

	for rows.Next() {
		var title_data string

		err := rows.Scan(&title_data)
		if err != nil {
			log.Fatal(err)
		}

		data_list = append(data_list, title_data)
	}

	return_data := make(map[string]interface{})
	return_data["language"] = map[string]string{
		"watchlist": tool.Get_language(db, "watchlist", false),
		"star_doc":  tool.Get_language(db, "star_doc", false),
	}
	return_data["data"] = data_list

	json_data, _ := json.Marshal(return_data)
	return string(json_data)
}
