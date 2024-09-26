package route

import (
	"log"
	"opennamu/route/tool"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

func Api_bbs_w_post(call_arg []string) string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &other_set)

	db := tool.DB_connect()
	defer db.Close()

	if tool.Check_acl(db, "", "", "bbs_comment", other_set["ip"]) {
		return_data := make(map[string]interface{})
		return_data["response"] = "require auth"

		json_data, _ := json.Marshal(return_data)
		return string(json_data)
	}
	
	stmt, err := db.Prepare(tool.DB_change("select set_code from bbs_data where set_name = 'title' and set_id = ? order by set_code + 0 desc"))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var set_code strong

	err := stmt.QueryRow(other_set["set_id"]).Scan(&set_code)
	if err != nil {
		log.Fatal(err)
	}

	set_code_int, _ := strconv.Atoi(set_code)
	set_code_int += 1

	return_data := make(map[string]interface{})
	return_data["response"] = "ok"

	json_data, _ := json.Marshal(return_data)
	return string(json_data)
}