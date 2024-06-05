package route

import (
	"encoding/json"
	"opennamu/route/tool"
)

func Api_func_auth_list(call_arg []string) string {
	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &other_set)

	db := tool.DB_connect()
	defer db.Close()

	auth_name := tool.Get_user_auth(db, other_set["ip"])
	auth_info := tool.Get_auth_group_info(db, auth_name)

	json_data, _ := json.Marshal(auth_info)
	return string(json_data)
}
