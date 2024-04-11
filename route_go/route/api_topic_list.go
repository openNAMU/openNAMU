package route

import (
	"encoding/json"
	"log"
	"opennamu/route/tool"
	"strconv"
)

func Api_topic_list(call_arg []string) string {
	db_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &db_set)

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[1]), &other_set)

	db := tool.DB_connect(db_set)
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

	stmt, err := db.Prepare(tool.DB_change(db_set, "select code, sub, stop, agree from rd where title = ? order by sub asc limit ?, 50"))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(other_set["name"], page_int)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var data_list [][]string

	for rows.Next() {
		var code string
		var sub string
		var stop string
		var agree string

		err := rows.Scan(&code, &sub, &stop, &agree)
		if err != nil {
			log.Fatal(err)
		}

		data_list = append(data_list, []string{code, sub, stop, agree})
	}

	return_data := make(map[string]interface{})
	return_data["language"] = map[string]string{}

	if len(data_list) == 0 {
		return_data["data"] = map[string]string{}
	} else {
		return_data["data"] = data_list
	}

	json_data, _ := json.Marshal(return_data)
	return string(json_data)
}
