package route

import (
	"database/sql"
	"log"
	"opennamu/route/tool"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

func Api_w_page_view(call_arg []string) string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &other_set)

	db := tool.DB_connect()
	defer db.Close()

	pv_continue := tool.Get_setting(db, "not_use_view_count", "")
	if len(pv_continue) == 0 || pv_continue[0][0] == "" {
		stmt, err := db.Prepare(tool.DB_change("select set_data from data_set where doc_name = ? and set_name = 'view_count'"))
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		var view_count string

		err = stmt.QueryRow(other_set["doc_name"]).Scan(&view_count)
		if err != nil {
			if err == sql.ErrNoRows {
				view_count = "0"
			} else {
				log.Fatal(err)
			}
		}

		if view_count == "0" {
			stmt, err := db.Prepare(tool.DB_change("insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, '', 'view_count', '1')"))
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()

			_, err = stmt.Exec(other_set["doc_name"])
			if err != nil {
				log.Fatal(err)
			}
		} else {
			view_count_int, _ := strconv.Atoi(view_count)

			stmt, err := db.Prepare(tool.DB_change("update data_set set set_data = ? where doc_name = ? and set_name = 'view_count'"))
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()

			_, err = stmt.Exec(view_count_int+1, other_set["doc_name"])
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	return_data := make(map[string]interface{})
	return_data["response"] = "ok"

	json_data, _ := json.Marshal(return_data)
	return string(json_data)
}
