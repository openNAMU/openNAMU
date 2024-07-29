package route

import (
	"database/sql"
	"log"
	"strconv"

	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func Api_func_search(call_arg []string) string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &other_set)

	page, _ := strconv.Atoi(other_set["num"])
	num := 0
	if page*50 > 0 {
		num = page*50 - 50
	}

	db := tool.DB_connect()
	defer db.Close()

	var stmt *sql.Stmt
	var err error
	if other_set["search_type"] == "title" {
		stmt, err = db.Prepare(tool.DB_change("select title from data where title collate nocase like ? order by title limit ?, 50"))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		stmt, err = db.Prepare(tool.DB_change("select title from data where data collate nocase like ? order by title limit ?, 50"))
		if err != nil {
			log.Fatal(err)
		}
	}
	defer stmt.Close()

	title_list := []string{}

	rows, err := stmt.Query("%"+other_set["name"]+"%", num)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var title string

		err := rows.Scan(&title)
		if err != nil {
			log.Fatal(err)
		}

		title_list = append(title_list, title)
	}

	json_data, _ := json.Marshal(title_list)
	return string(json_data)
}
