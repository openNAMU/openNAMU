package route

import (
	"database/sql"
	"encoding/json"
	"log"
	"opennamu/route/tool"
)

func Api_setting(call_arg []string) string {
	db_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &db_set)

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[1]), &other_set)

	db := tool.DB_connect(db_set)
	defer db.Close()

	setting_acl := map[string]string{}

	setting_acl["manage_404_page"] = ""
	setting_acl["manage_404_page_content"] = ""

	if val, ok := setting_acl[other_set["set_name"]]; ok {
		if val != "" {
			auth_name := tool.Get_user_auth(db, db_set, other_set["ip"])
			auth_info := tool.Get_auth_group_info(db, db_set, auth_name)

			if _, ok := auth_info["owner"]; !ok {
				return_data := make(map[string]interface{})
				return_data["response"] = "require auth"

				json_data, _ := json.Marshal(return_data)
				return string(json_data)
			}
		}

		stmt, err := db.Prepare(tool.DB_change(db_set, "select data from other where name = ?"))
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		var set_data string

		err = stmt.QueryRow(other_set["set_name"]).Scan(&set_data)
		if err != nil {
			if err == sql.ErrNoRows {
				set_data = ""
			} else {
				log.Fatal(err)
			}
		}

		return_data := make(map[string]interface{})
		return_data["response"] = "ok"
		return_data["data"] = set_data

		json_data, _ := json.Marshal(return_data)
		return string(json_data)
	} else {
		return_data := make(map[string]interface{})
		return_data["response"] = "not exist"

		json_data, _ := json.Marshal(return_data)
		return string(json_data)
	}
}
