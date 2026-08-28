package route

import "opennamu/route/tool"

func Api_setting_main_post(config tool.Config, form map[string]string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "setting_main", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if !acl_value_valid(db, form["user_document_view_acl_all"]) {
		return_data["response"] = "error"
		return_data["data"] = "error"
		return return_data
	}
	setting_save_fields(db, setting_main_fields(), form)
	tool.Do_insert_auth_history(db, config.IP, "edit_set (main)")
	return_data["response"] = "ok"
	return return_data
}
