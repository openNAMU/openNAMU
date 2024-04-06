package route

import (
	"database/sql"
	"encoding/json"
	"log"
	"opennamu/route/tool"
	"strconv"
)

func Api_list_old_page(call_arg []string) {
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

	page_int, err := strconv.Atoi(other_set["num"])
	if err != nil {
		log.Fatal(err)
	}

	if page_int > 0 {
		page_int = (page_int * limit_int) - limit_int
	} else {
		page_int = 0
	}

	var stmt *sql.Stmt

	if other_set["set_type"] == "old" {
		stmt, err = db.Prepare(tool.DB_change(db_set, "select doc_name, set_data from data_set where set_name = 'last_edit' and doc_rev = '' order by set_data desc limit ?, 50"))
	} else {
		stmt, err = db.Prepare(tool.DB_change(db_set, "select doc_name, set_data from data_set where set_name = 'last_edit' and doc_rev = '' order by set_data asc limit ?, 50"))
	}

	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(page_int, limit_int)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
}
