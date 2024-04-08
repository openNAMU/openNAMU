package route

import (
	"encoding/json"
	"opennamu/route/tool"
)

func Api_func_auth_list(call_arg []string) string {
	db_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &db_set)

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[1]), &other_set)

	db := tool.DB_connect(db_set)
	defer db.Close()

	auth_name := tool.Get_user_auth(db, db_set, other_set["ip"])
	auth_info := tool.Get_auth_group_info(db, db_set, auth_name)

	json_data, _ := json.Marshal(auth_info)
	return string(json_data)
}
