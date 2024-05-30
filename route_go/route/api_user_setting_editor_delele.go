package route

import (
	"encoding/json"
	"log"
	"opennamu/route/tool"
)

func Api_user_setting_editor_delete(call_arg []string) string {
	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &other_set)

	db := tool.DB_connect()
	defer db.Close()

	ip := other_set["ip"]
	if !tool.IP_or_user(ip) {
		stmt, err := db.Prepare(tool.DB_change("delete from user_set where id = ? and name = 'user_editor_top' and data = ?"))
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(ip, other_set["data"])
		if err != nil {
			log.Fatal(err)
		}

		return_data := make(map[string]interface{})
		return_data["response"] = "ok"
		return_data["language"] = map[string]string{
			"delete": tool.Get_language(db, "delete", false),
		}

		json_data, _ := json.Marshal(return_data)
		return string(json_data)
	} else {
		return_data := make(map[string]interface{})
		return_data["response"] = "require auth"
		return_data["language"] = map[string]string{
			"authority_error": tool.Get_language(db, "authority_error", false),
		}

		json_data, _ := json.Marshal(return_data)
		return string(json_data)
	}
}
