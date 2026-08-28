package route

import "opennamu/route/tool"

func Api_server_action_post(config tool.Config, action string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if action != "restart" && action != "shutdown" && action != "update" {
		return_data["response"] = "error"
		return return_data
	}
	tool.Do_insert_auth_history(db, config.IP, "server_"+action)

	return_data["response"] = "ok"
	return return_data
}
