package route

import (
	"database/sql"
	"encoding/json"
	"opennamu/route/tool"
	"strconv"
)

func Api_list_recent_discuss(call_arg []string) {
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
		return
	}

	if limit_int > 50 || limit_int < 0 {
		limit_int = 50
	}

	var stmt *sql.Stmt

	set_type := other_set["set_type"]
	if set_type == "normal" {
		stmt, err = db.Prepare(tool.DB_change(db_set, "select title, sub, date, code from rd where not stop = 'O' order by date desc limit ?"))
	} else if set_type == "close" {
		stmt, err = db.Prepare(tool.DB_change(db_set, "select title, sub, date, code from rd where stop = 'O' order by date desc limit ?"))
	} else {
		stmt, err = db.Prepare(tool.DB_change(db_set, "select title, sub, date, code from rd where stop != 'O' order by date asc limit ?"))
	}

	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(limit_int)
	if err != nil {
		return
	}
	defer rows.Close()

	// var data_list [][]string
	// admin_auth := tool.Get_admin_auth(db, db_set, other_set["ip"])

	for rows.Next() {
		var title string
		var sub string
		var date string
		var code string

		err := rows.Scan(&title, &sub, &date, &code)
		if err != nil {
			return
		}

		stmt, err := db.Prepare(tool.DB_change(db_set, "select ip from topic where code = ? order by id + 0 desc limit 1"))
		if err != nil {
			return
		}
		defer stmt.Close()

		var ip string

		err = stmt.QueryRow(code).Scan(&ip)
		if err != nil {
			if err == sql.ErrNoRows {
				ip = ""
			} else {
				return
			}
		}
	}
}
