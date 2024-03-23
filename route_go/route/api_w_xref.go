package route

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"opennamu/route/tool"
	"strconv"
)

func Api_w_xref(call_arg []string) {
	db_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &db_set)

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[1]), &other_set)

	db := tool.DB_connect(db_set)
	if db == nil {
		return
	}
	defer db.Close()

	page, _ := strconv.Atoi(other_set["page"])
	num := 0
	if page*50 > 0 {
		num = page*50 - 50
	}

	var link_case_insensitive string

	err := db.QueryRow(tool.DB_change(db_set, "select data from other where name = 'link_case_insensitive'")).Scan(&link_case_insensitive)
	if err != nil {
		if err == sql.ErrNoRows {
			link_case_insensitive = ""
		} else {
			log.Fatal(err)
		}
	}

	if link_case_insensitive != "" {
		link_case_insensitive = " collate nocase"
	}

	var stmt *sql.Stmt
	if other_set["do_type"] == "1" {
		stmt, err = db.Prepare(tool.DB_change(db_set, "select distinct link, type from back where title"+link_case_insensitive+" = ? and not type = 'no' and not type = 'nothing' order by type asc, link asc limit ?, 50"))
	} else {
		stmt, err = db.Prepare(tool.DB_change(db_set, "select distinct title, type from back where link"+link_case_insensitive+" = ? and not type = 'no' and not type = 'nothing' order by type asc, title asc limit ?, 50"))
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

	var name string
	var type_data string

	var data_list [][]string

	for rows.Next() {
		err := rows.Scan(&name, &type_data)
		if err != nil {
			log.Fatal(err)
		}

		data_list = append(data_list, []string{name, type_data})
	}

	if len(data_list) == 0 {
		fmt.Print("{}")
	} else {
		json_data, _ := json.Marshal(data_list)
		fmt.Print(string(json_data))
	}
}
