package route

import (
	"log"
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func Api_bbs_w_set_put(call_arg []string) string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &other_set)

	db := tool.DB_connect()
	defer db.Close()

	auth_info := tool.Check_acl(db, "", "", "owner_auth", other_set["ip"])

	setting_acl := BBS_w_set_list()
	return_data := make(map[string]interface{})

	if _, ok := setting_acl[other_set["set_name"]]; ok {
		if auth_info {
			if _, ok := other_set["coverage"]; !ok {
				stmt, err := db.Prepare(tool.DB_change("delete from bbs_set where set_name = ? and set_id = ?"))
				if err != nil {
					log.Fatal(err)
				}
				defer stmt.Close()

				_, err = stmt.Exec(other_set["set_name"], other_set["set_id"])
				if err != nil {
					log.Fatal(err)
				}
			}

			stmt, err := db.Prepare(tool.DB_change("insert into bbs_set (set_name, set_code, set_id, set_data) values (?, '', ?, ?)"))
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()

			_, err = stmt.Exec(other_set["set_name"], other_set["set_id"], other_set["data"])
			if err != nil {
				log.Fatal(err)
			}

			return_data["response"] = "ok"
		} else {
			return_data["response"] = "require auth"
		}
	} else {
		return_data["response"] = "not exist"
	}

	json_data, _ := json.Marshal(return_data)
	return string(json_data)
}
