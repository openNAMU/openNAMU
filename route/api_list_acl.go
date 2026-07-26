package route

import (
	"opennamu/route/tool"
)

func Api_list_acl(config tool.Config, func_type string) map[string]any {
    data := tool.List_acl(func_type)

    return_data := make(map[string]any)
    return_data["response"] = "ok"
    return_data["data"] = data

    return return_data
}
