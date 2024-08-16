package route

import (
	"log"
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func Api_w_comment_delete(call_arg []string) string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &other_set)

	db := tool.DB_connect()
	defer db.Close()

	return_data := make(map[string]interface{})

	pass := false
	if !tool.Check_acl(db, "", "", "comment_manager", other_set["ip"]) {
		return_data["response"] = "require auth"
	} else {
		pass = true
	}

	if pass {
		return_data["response"] = "ok"

		for _, v := range Comment_need_list() {
			stmt, err := db.Prepare(tool.DB_change("delete from data_set where doc_name = ? and doc_rev = ? and set_name = ?"))
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()

			_, err = stmt.Exec(other_set["doc_name"], other_set["num"], v)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	json_data, _ := json.Marshal(return_data)
	return string(json_data)
}
