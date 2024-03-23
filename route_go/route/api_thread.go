package route

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"opennamu/route/tool"
)

func Api_thread(call_arg []string) {
	db_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &db_set)

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[1]), &other_set)

	db := tool.DB_connect(db_set)
	if db == nil {
		return
	}
	defer db.Close()

	if other_set["tool"] == "length" {
		stmt, err := db.Prepare(tool.DB_change(db_set, "select id from topic where code = ? order by id + 0 desc limit 1"))
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		var length string
		err = stmt.QueryRow(other_set["topic_num"]).Scan(&length)
		if err != nil {
			if err == sql.ErrNoRows {
				length = "0"
			} else {
				log.Fatal(err)
			}
		}

		new_data := map[string]string{}
		new_data["length"] = length

		json_data, _ := json.Marshal(new_data)
		fmt.Print(string(json_data))
	} else {
		var rows *sql.Rows

		if other_set["tool"] == "top" {
			stmt, err := db.Prepare(tool.DB_change(db_set, "select id, data, date, ip, block, top from topic where code = ? and top = 'O' order by id + 0 asc"))
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()

			rows, err = stmt.Query(other_set["topic_num"])
			if err != nil {
				log.Fatal(err)
			}
		} else {
			if other_set["s_num"] != "" && other_set["e_num"] != "" {
				stmt, err := db.Prepare(tool.DB_change(db_set, "select id, data, date, ip, block, top from topic where code = ? and ? + 0 <= id + 0 and id + 0 <= ? + 0 order by id + 0 asc"))
				if err != nil {
					log.Fatal(err)
				}
				defer stmt.Close()

				rows, err = stmt.Query(other_set["topic_num"], other_set["s_num"], other_set["e_num"])
				if err != nil {
					log.Fatal(err)
				}
			} else {
				stmt, err := db.Prepare(tool.DB_change(db_set, "select id, data, date, ip, block, top from topic where code = ? order by id + 0 asc"))
				if err != nil {
					log.Fatal(err)
				}
				defer stmt.Close()

				rows, err = stmt.Query(other_set["topic_num"])
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		defer rows.Close()

		var id, data, date, ip, block, top string
		var data_list [][]string

		for rows.Next() {
			err := rows.Scan(&id, &data, &date, &ip, &block, &top)
			if err != nil {
				log.Fatal(err)
			}

			data_list = append(data_list, []string{id, data, date, ip, block, top})
		}

		new_data := make(map[string]interface{})
		new_data["data"] = []map[string]string{}
		data_slice := []map[string]string{}

		admin_auth := tool.Get_user_auth(db, db_set, other_set["ip"])

		for for_a := 0; for_a < len(data_list); for_a++ {
			data := ""
			if data_list[for_a][4] != "O" || admin_auth != "" {
				data = data_list[for_a][1]
			}

			data_slice = append(data_slice, map[string]string{
				"id":        data_list[for_a][0],
				"data":      data,
				"date":      data_list[for_a][2],
				"ip":        tool.IP_preprocess(db, db_set, data_list[for_a][3], other_set["ip"])[0],
				"ip_render": tool.IP_parser(db, db_set, data_list[for_a][3], other_set["ip"]),
				"blind":     data_list[for_a][4],
			})
		}

		new_data["data"] = data_slice
		new_data["language"] = map[string]string{
			"tool":   tool.Get_language(db, db_set, "tool", false),
			"render": tool.Get_language(db, db_set, "render", false),
		}

		json_data, _ := json.Marshal(new_data)
		fmt.Print(string(json_data))
	}
}
