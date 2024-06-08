package route

import (
	"encoding/json"

	"opennamu/route/tool"
)

func Api_func_sha224(call_arg []string) string {
	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &other_set)

	hash_str := tool.Sha224(other_set["data"])

	return_data := make(map[string]interface{})
	return_data["response"] = "ok"
	return_data["data"] = hash_str

	json_data, _ := json.Marshal(return_data)
	return string(json_data)
}
