package route

import (
	"encoding/json"
	"log"
	"opennamu/route/tool"
)

func Api_bbs_comment(call_arg []string) string {
	db_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &db_set)

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[1]), &other_set)

	db := tool.DB_connect(db_set)
	defer db.Close()

	if other_set["tool"] == "length" {
		stmt, err := db.Prepare(tool.DB_change(db_set, "select set_name, set_data, set_code, set_id from bbs_data where (set_name = 'comment' or set_name = 'comment_date' or set_name = 'comment_user_id') and set_id = ? order by set_code + 0 asc"))
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		rows, err := stmt.Query(other_set["bbs_num"])
		if err != nil {
			log.Fatal(err)
		}

	} else {
		return ""
	}
}
