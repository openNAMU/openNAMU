package route

import (
	"encoding/json"
	"fmt"

	"opennamu/route/tool"
)

func Api_func_ip(call_arg []string) {
	db_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &db_set)

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[1]), &other_set)

	db := tool.DB_connect(db_set)
	if db == nil {
		return
	}
	defer db.Close()

	ip_data := tool.IP_parser(db, db_set, other_set["data"], other_set["ip"])

	new_data := map[string]string{}
	new_data["data"] = ip_data

	json_data, _ := json.Marshal(new_data)
	fmt.Print(string(json_data))
}
