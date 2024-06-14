package route

import (
	"encoding/json"
	"log"
	"opennamu/route/tool"
	"strconv"
)

func Api_list_title_index(call_arg []string) string {
	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &other_set)

	db := tool.DB_connect()
	defer db.Close()

	page_int, err := strconv.Atoi(other_set["num"])
	if err != nil {
		log.Fatal(err)
	}

	if page_int > 0 {
		page_int = (page_int * 50) - 50
	} else {
		page_int = 0
	}

	stmt, err := db.Prepare(tool.DB_change("select title from data limit ?, 50"))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(page_int)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	data_list := []string{}

	for rows.Next() {
		var title string

		err := rows.Scan(&title)
		if err != nil {
			log.Fatal(err)
		}

		data_list = append(data_list, title)
	}

	return_data := make(map[string]interface{})
	return_data["data"] = data_list

	json_data, _ := json.Marshal(return_data)
	return string(json_data)
}
