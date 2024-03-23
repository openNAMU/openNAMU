package route

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"opennamu/route/tool"
)

func Api_w_random(call_arg []string) {
	db_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &db_set)

	db := tool.DB_connect(db_set)
	if db == nil {
		return
	}
	defer db.Close()

	var title string

	err := db.QueryRow(tool.DB_change(db_set, "select title from data where title not like 'user:%' and title not like 'category:%' and title not like 'file:%' order by random() limit 1")).Scan(&title)
	if err != nil {
		if err == sql.ErrNoRows {
			title = ""
		} else {
			log.Fatal(err)
		}
	}

	new_data := map[string]string{}
	new_data["data"] = title

	json_data, _ := json.Marshal(new_data)
	fmt.Print(string(json_data))
}
