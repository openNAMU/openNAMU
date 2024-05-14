package route

import (
	"encoding/json"

	"opennamu/route/tool"
)

func Api_func_ip_menu(call_arg []string) string {
	db_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &db_set)

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[1]), &other_set)

	db := tool.DB_connect(db_set)
	defer db.Close()

	ip_data := tool.IP_menu(db, db_set, other_set["ip"], other_set["my_ip"], other_set["option"])

	new_data := make(map[string]interface{})
	new_data["data"] = ip_data

	json_data, _ := json.Marshal(new_data)
	return string(json_data)
}
