package route

import (
	"opennamu/route/tool"
)

func Api_func_acl(config tool.Config, name string, topic_number string, tool_name string, ip string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if ip == "" {
		ip = config.IP
	}

	new_data := make(map[string]any)
	new_data["response"] = "ok"
	new_data["data"] = tool.Check_acl(db, name, topic_number, tool_name, ip)

	return new_data
}
