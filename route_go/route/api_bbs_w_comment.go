package route

import (
	"database/sql"
	"encoding/json"
	"log"
	"opennamu/route/tool"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

func Api_bbs_w_comment_all(sub_code string) []map[string]string {
	end_data := []map[string]string{}

	inter_other_set := map[string]string{}
	inter_other_set["sub_code"] = sub_code
	inter_other_set["tool"] = "around"
	inter_other_set["legacy"] = "on"

	json_data, _ := json.Marshal(inter_other_set)
	return_data := Api_bbs_w_comment_one([]string{string(json_data)})

	return_data_api := []map[string]string{}
	json.Unmarshal([]byte(return_data), &return_data_api)

	for for_a := 0; for_a < len(return_data_api); for_a++ {
		end_data = append(end_data, return_data_api[for_a])

		temp := Api_bbs_w_comment_all(sub_code + "-" + return_data_api[for_a]["code"])
		if len(temp) > 0 {
			for for_b := 0; for_b < len(temp); for_b++ {
				end_data = append(end_data, temp[for_b])
			}
		}
	}

	return end_data
}

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
		bbs_and_post_num := other_set["sub_code"]

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
		temp := Api_bbs_w_comment_all(other_set["sub_code"])

		if other_set["legacy"] != "" {
			json_data, _ := json.Marshal(temp)
			return string(json_data)
		} else {
			return_data := make(map[string]interface{})
			return_data["language"] = map[string]string{
				"normal":  tool.Get_language(db, "normal", false),
				"comment": tool.Get_language(db, "comment", false),
				"tool":    tool.Get_language(db, "tool", false),
				"return":  tool.Get_language(db, "return", false),
			}
			return_data["data"] = temp

			json_data, _ := json.Marshal(return_data)
			return string(json_data)
		}
	}
}
