package route

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"opennamu/route/tool"
)

func Api_w_raw(call_arg []string) {
	db_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &db_set)

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[1]), &other_set)

	db := tool.DB_connect(db_set)
	if db == nil {
		return
	}
	defer db.Close()

	if other_set["exist_check"] != "" {
		stmt, err := db.Prepare(tool.DB_change(db_set, "select title from data where title = ?"))
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		new_data := map[string]string{}
		var title string

		err = stmt.QueryRow(other_set["name"]).Scan(&title)
		if err != nil {
			if err == sql.ErrNoRows {
			} else {
				log.Fatal(err)
			}
		} else {
			new_data["exist"] = "1"
		}

		json_data, _ := json.Marshal(new_data)
		fmt.Print(string(json_data))
	} else {
		new_data := map[string]string{}
		var data string

		if other_set["rev"] != "" {
			stmt, err := db.Prepare(tool.DB_change(db_set, "select data from history where title = ? and id = ?"))
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()

			err = stmt.QueryRow(other_set["name"], other_set["rev"]).Scan(&data)
			if err != nil {
				if err == sql.ErrNoRows {
				} else {
					log.Fatal(err)
				}
			} else {
				new_data["title"] = other_set["name"]
				new_data["data"] = data
			}

			json_data, _ := json.Marshal(new_data)
			fmt.Print(string(json_data))
		} else {
			stmt, err := db.Prepare(tool.DB_change(db_set, "select data from data where title = ?"))
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()

			err = stmt.QueryRow(other_set["name"]).Scan(&data)
			if err != nil {
				if err == sql.ErrNoRows {
				} else {
					log.Fatal(err)
				}
			} else {
				new_data["title"] = other_set["name"]
				new_data["data"] = data
			}

			json_data, _ := json.Marshal(new_data)
			fmt.Print(string(json_data))
		}
	}
}
