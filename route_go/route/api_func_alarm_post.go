package route

import (
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func Api_func_alarm_post(call_arg []string) string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &other_set)

	db := tool.DB_connect()
	defer db.Close()

	tool.Send_alarm(db, other_set["from"], other_set["to"], other_set["data"])

	return_data := make(map[string]interface{})
	return_data["response"] = "ok"

	json_data, _ := json.Marshal(return_data)
	return string(json_data)
}
