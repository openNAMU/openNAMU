package route

import (
	"database/sql"
	"log"
	"opennamu/route/tool"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

func Api_bbs_w_comment(call_arg []string) string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &other_set)

	db := tool.DB_connect()
	defer db.Close()

	if other_set["tool"] == "length" {
		stmt, err := db.Prepare(tool.DB_change("select count(*) from bbs_data where set_name = 'comment_date' and set_id = ? order by set_code + 0 desc"))
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		var comment_length string
		bbs_and_post_num := other_set["bbs_num"] + "-" + other_set["post_num"]

		err = stmt.QueryRow(bbs_and_post_num).Scan(&comment_length)
		if err != nil {
			if err == sql.ErrNoRows {
				comment_length = "0"
			} else {
				log.Fatal(err)
			}
		}

		stmt, err = db.Prepare(tool.DB_change("select count(*) from bbs_data where set_name = 'comment_date' and set_id like ? order by set_code + 0 desc"))
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		var reply_length string

		err = stmt.QueryRow(bbs_and_post_num + "-%").Scan(&reply_length)
		if err != nil {
			if err == sql.ErrNoRows {
				reply_length = "0"
			} else {
				log.Fatal(err)
			}
		}

		comment_length_int, _ := strconv.Atoi(comment_length)
		reply_length_int, _ := strconv.Atoi(reply_length)

		length_int := comment_length_int + reply_length_int
		length_str := strconv.Itoa(length_int)

		data_list := map[string]string{
			"comment": comment_length,
			"reply":   reply_length,
			"data":    length_str,
		}

		json_data, _ := json.Marshal(data_list)
		return string(json_data)
	} else {
		return "{}"
	}
}
