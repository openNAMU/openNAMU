package route

import (
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func Api_func_auth(call_arg []string) string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &other_set)

	db := tool.DB_connect()
	defer db.Close()

	auth_name := tool.Get_user_auth(db, other_set["ip"])
	auth_info := tool.Get_auth_group_info(db, auth_name)

	return_data := make(map[string]interface{})
	return_data["response"] = "ok"
	return_data["name"] = auth_name
	return_data["info"] = auth_info

	json_data, _ := json.Marshal(return_data)
	return string(json_data)
}
