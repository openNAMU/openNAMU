package route

import (
	"database/sql"
	"encoding/json"
	"log"
	"opennamu/route/tool"
)

func Api_bbs_w_comment(call_arg []string) string {
	db_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &db_set)

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[1]), &other_set)

	db := tool.DB_connect(db_set)
	defer db.Close()

	if other_set["tool"] == "length" {
		stmt, err := db.Prepare(tool.DB_change(db_set, "select count(*) from bbs_data where set_name = 'comment_date' and (set_id = ? or set_id like ?) order by set_code + 0 desc"))
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		var comment_length string
		bbs_and_post_num := other_set["bbs_num"] + "-" + other_set["post_num"]

		err = stmt.QueryRow(bbs_and_post_num, bbs_and_post_num + "-%").Scan(&comment_length)
		if err != nil {
			if err == sql.ErrNoRows {
				comment_length = "0"
			} else {
				log.Fatal(err)
			}
		}

		data_list := map[string]string{
			"data": comment_length,
		}

		json_data, _ := json.Marshal(data_list)
		return string(json_data)
	} else {
		return ""
	}
}
