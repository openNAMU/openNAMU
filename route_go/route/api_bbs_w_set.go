package route

import (
	"log"
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func BBS_w_set_list() map[string]string {
	setting_acl := map[string]string{}

	setting_acl["bbs_view_acl"] = ""
	setting_acl["bbs_acl"] = ""
	setting_acl["bbs_edit_acl"] = ""
	setting_acl["bbs_comment_acl"] = ""

	setting_acl["bbs_markup"] = ""
	setting_acl["bbs_name"] = ""

	return setting_acl
}

func Api_bbs_w_set(call_arg []string) string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &other_set)

	db := tool.DB_connect()
	defer db.Close()

	setting_acl := BBS_w_set_list()

	if val, ok := setting_acl[other_set["set_name"]]; ok {
		if val != "" {
			if tool.Check_acl(db, "", "", "owner_auth", other_set["ip"]) {
				return_data := make(map[string]interface{})
				return_data["response"] = "require auth"

				json_data, _ := json.Marshal(return_data)
				return string(json_data)
			}
		}

		stmt, err := db.Prepare(tool.DB_change("select set_data, set_code from bbs_set where set_name = ? and set_id = ?"))
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		rows, err := stmt.Query(other_set["set_name"], other_set["set_id"])
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		data_list := [][]string{}

		for rows.Next() {
			var set_data string
			var set_coverage string

			err := rows.Scan(&set_data, &set_coverage)
			if err != nil {
				log.Fatal(err)
			}

			data_list = append(data_list, []string{set_data, set_coverage})
		}

		return_data := make(map[string]interface{})
		return_data["response"] = "ok"
		return_data["data"] = data_list

		json_data, _ := json.Marshal(return_data)
		return string(json_data)
	} else {
		return_data := make(map[string]interface{})
		return_data["response"] = "not exist"

		json_data, _ := json.Marshal(return_data)
		return string(json_data)
	}
}
