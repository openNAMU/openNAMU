package route

import "opennamu/route/tool"

func Api_setting_head_post(config tool.Config, name string, coverage string, content string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "setting_head", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if name == "" {
		return_data["response"] = "error"
		return return_data
	}
	setting_save_value(db, name, coverage, content)
	tool.Do_insert_auth_history(db, config.IP, "edit_set ("+name+")")
	return_data["response"] = "ok"
	return return_data
}
