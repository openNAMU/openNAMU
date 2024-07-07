package route

import (
	"encoding/json"
	"opennamu/route/tool"
)

func Api_list_auth(call_arg []string) string {
	db := tool.DB_connect()
	defer db.Close()

	data := tool.List_auth(db)

	return_data := make(map[string]interface{})
	return_data["response"] = "ok"
	return_data["data"] = data

	json_data, _ := json.Marshal(return_data)
	return string(json_data)
}
