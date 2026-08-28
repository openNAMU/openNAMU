package route

import (
	"opennamu/route/tool"
)

func Api_list_acl(config tool.Config, func_type string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	data := tool.List_acl(func_type)
	for _, group := range tool.List_auth(db) {
		if !tool.Arr_in_str(data, group) {
			data = append(data, group)
		}
	}

	return_data := make(map[string]any)
	return_data["response"] = "ok"
	return_data["data"] = data

	return return_data
}
