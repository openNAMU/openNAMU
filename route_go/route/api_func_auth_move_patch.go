package route

import (
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func Api_func_auth_move_patch(call_arg []string) string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &other_set)

	db := tool.DB_connect()
	defer db.Close()

	if tool.Check_acl() {
		new_data := make(map[string]interface{})
		new_data["response"] = "ok"

		json_data, _ := json.Marshal(new_data)
		return string(json_data)
	} else {

	}
}
