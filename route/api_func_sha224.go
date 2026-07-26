package route

import (
	"opennamu/route/tool"
)

func Api_func_sha224(config tool.Config, data string) map[string]any {
    hash_str := tool.Sha224(data)

    return_data := make(map[string]any)
    return_data["response"] = "ok"
    return_data["data"] = hash_str

    return return_data
}
