package route

import "opennamu/route/tool"

func Api_user_name_post(config tool.Config, target string, value string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if target == "" {
		target = config.IP
	}
	if target == config.IP && tool.IP_or_user(config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if target != config.IP && !tool.Check_permission(db, "user_manage", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	current := user_value(db, target, "user_name")
	if value == "" || (!tool.Get_user_name_check(db, value) && value != current) {
		return_data["response"] = "error"
		return_data["data"] = "user name error"
		return return_data
	}
	user_save(db, target, "user_name", value)
	return_data["response"] = "ok"
	return return_data
}
