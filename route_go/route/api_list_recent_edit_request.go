package route

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"opennamu/route/tool"
	"strconv"
)

func Api_list_recent_edit_request(call_arg []string) {
	db_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &db_set)

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[1]), &other_set)

	db := tool.DB_connect(db_set)
	if db == nil {
		return
	}
	defer db.Close()

	limit_int, err := strconv.Atoi(other_set["limit"])
	if err != nil {
		log.Fatal(err)
	}

	if limit_int > 50 || limit_int < 0 {
		limit_int = 50
	}

	stmt, err := db.Prepare(tool.DB_change(db_set, "select doc_name, doc_rev, set_data from data_set where set_name = 'edit_request_doing' order by set_data desc limit ?"))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(limit_int)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var doc_name string
	var doc_rev string
	var date string

	var data_list [][]string
	for rows.Next() {
		err := rows.Scan(&doc_name, &doc_rev, &date)
		if err != nil {
			log.Fatal(err)
		}

		var ip string
		var send string
		var leng string

		stmt, err := db.Prepare(tool.DB_change(db_set, "select set_data from data_set where set_name = 'edit_request_user' and doc_rev = ? and doc_name = ?"))
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		err = stmt.QueryRow(doc_rev, doc_name).Scan(&ip)
		if err != nil {
			if err == sql.ErrNoRows {
				ip = ""
			} else {
				log.Fatal(err)
			}
		}

		stmt, err = db.Prepare(tool.DB_change(db_set, "select set_data from data_set where set_name = 'edit_request_send' and doc_rev = ? and doc_name = ?"))
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		err = stmt.QueryRow(doc_rev, doc_name).Scan(&send)
		if err != nil {
			if err == sql.ErrNoRows {
				send = ""
			} else {
				log.Fatal(err)
			}
		}

		stmt, err = db.Prepare(tool.DB_change(db_set, "select set_data from data_set where set_name = 'edit_request_leng' and doc_rev = ? and doc_name = ?"))
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		err = stmt.QueryRow(doc_rev, doc_name).Scan(&leng)
		if err != nil {
			if err == sql.ErrNoRows {
				leng = ""
			} else {
				log.Fatal(err)
			}
		}

		data_list = append(data_list, []string{
			doc_name,
			doc_rev,
			date,
			tool.IP_preprocess(db, db_set, ip, other_set["ip"])[0],
			send,
			leng,
			tool.IP_parser(db, db_set, ip, other_set["ip"]),
		})
	}

	if len(data_list) == 0 {
		fmt.Print("{}")
	} else {
		json_data, _ := json.Marshal(data_list)
		fmt.Print(string(json_data))
	}
}
