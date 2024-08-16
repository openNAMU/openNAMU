package route

import (
	"log"
	"opennamu/route/tool"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

func Api_w_comment(call_arg []string) string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &other_set)

	db := tool.DB_connect()
	defer db.Close()

	return_data := make(map[string]interface{})

	comment_enable := tool.Get_setting(db, "enable_comment", "")
	if len(comment_enable) == 0 || comment_enable[0][0] == "" {
		return_data["response"] = "disable"
	} else if !tool.Check_acl(db, "", "", "comment_view", other_set["ip"]) {
		return_data["response"] = "require auth"
	} else {
		return_data["response"] = "ok"

		page_int, err := strconv.Atoi(other_set["num"])
		if err != nil {
			log.Fatal(err)
		}

		if page_int > 0 {
			page_int = (page_int * 50) - 50
		} else {
			page_int = 0
		}

		stmt, err := db.Prepare(tool.DB_change("select set_name, doc_rev, set_data from data_set where (set_name = 'comment' or set_name like 'comment%') and doc_name = ? order by doc_rev asc limit ?, 50"))
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		rows, err := stmt.Query(other_set["doc_name"], page_int)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		data_list := []map[string]string{}
		temp_dict := map[string]string{}
		before_set_code := ""

		for rows.Next() {
			var set_name string
			var set_code string
			var set_data string

			err := rows.Scan(&set_name, &set_code, &set_data)
			if err != nil {
				log.Fatal(err)
			}

			if before_set_code != set_code {
				if before_set_code != "" {
					data_list = append(data_list, temp_dict)
				}

				temp_dict = map[string]string{}
				temp_dict["id"] = set_code

				before_set_code = set_code
			}

			temp_dict[set_name] = set_data
		}

		if before_set_code != "" {
			data_list = append(data_list, temp_dict)
		}

		return_data["data"] = data_list
	}

	json_data, _ := json.Marshal(return_data)
	return string(json_data)
}
