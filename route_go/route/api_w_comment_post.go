package route

import (
	"database/sql"
	"log"
	"opennamu/route/tool"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

func Comment_need_list() []string {
	return []string{"comment", "comment_date", "comment_user_id"}
}

func Api_w_comment_post(call_arg []string) string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &other_set)

	db := tool.DB_connect()
	defer db.Close()

	ip := other_set["ip"]

	stmt, err := db.Prepare(tool.DB_change("select doc_rev from data_set where set_name = 'comment' and doc_name = ? order by doc_rev asc limit 1"))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var view_count string

	err = stmt.QueryRow(other_set["doc_name"]).Scan()
	if err != nil {
		if err == sql.ErrNoRows {
			view_count = "0"
		} else {
			log.Fatal(err)
		}
	}

	view_count_int, _ := strconv.Atoi(view_count)
	view_count_int += 1

	stmt, err = db.Prepare(tool.DB_change("insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, ?, ?, ?)"))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(other_set["doc_name"], view_count)
	if err != nil {
		log.Fatal(err)
	}

	return_data := make(map[string]interface{})
	return_data["response"] = "ok"

	json_data, _ := json.Marshal(return_data)
	return string(json_data)
}
