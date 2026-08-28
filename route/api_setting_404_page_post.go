package route

import "opennamu/route/tool"

func Api_setting_404_page_post(config tool.Config, page string, content string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "setting_404", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	setting_save_value(db, "manage_404_page", "", page)
	setting_save_value(db, "manage_404_page_content", "", content)
	tool.Do_insert_auth_history(db, config.IP, "edit_set (404_page)")
	return_data["response"] = "ok"
	return return_data
}
