package route

import (
	"opennamu/route/tool"
)

func Api_func_sha224(config tool.Config) string {
    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    hash_str := tool.Sha224(other_set["data"])

    return_data := make(map[string]any)
    return_data["response"] = "ok"
    return_data["data"] = hash_str

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}
