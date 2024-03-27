package route

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"opennamu/route/tool"
)

func bbs_list(db *sql.DB, db_set map[string]string) map[string]string {
	rows, err := db.Query(tool.DB_change(db_set, "select set_data, set_id from bbs_set where set_name = 'bbs_name'"))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	data_list := map[string]string{}

	for rows.Next() {
		var name string
		var id string

		err := rows.Scan(&name, &id)
		if err != nil {
			log.Fatal(err)
		}

		data_list[name] = id
	}

	return data_list
}

func Api_bbs_list(call_arg []string) {
	db_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &db_set)

	db := tool.DB_connect(db_set)
	if db == nil {
		return
	}
	defer db.Close()

	data_list := bbs_list(db, db_set)
	data_list_sub := map[string][]string{}

	for k, v := range data_list {
		stmt, err := db.Prepare(tool.DB_change(db_set, "select set_data from bbs_set where set_name = 'bbs_type' and set_id = ?"))
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		var bbs_type string

		err = stmt.QueryRow(v).Scan(&bbs_type)
		if err != nil {
			if err == sql.ErrNoRows {
				bbs_type = ""
			} else {
				log.Fatal(err)
			}
		}

		stmt, err = db.Prepare(tool.DB_change(db_set, "select set_data from bbs_data where set_id = ? and set_name = 'date' order by set_code + 0 desc limit 1"))
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		var bbs_date string

		err = stmt.QueryRow(v).Scan(&bbs_date)
		if err != nil {
			if err == sql.ErrNoRows {
				bbs_date = ""
			} else {
				log.Fatal(err)
			}
		}

		data_list_sub[k] = []string{v, bbs_type, bbs_date}
	}

	if len(data_list_sub) == 0 {
		fmt.Print("{}")
	} else {
		json_data, _ := json.Marshal(data_list_sub)
		fmt.Print(string(json_data))
	}
}
