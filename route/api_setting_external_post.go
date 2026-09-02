package route

import "opennamu/route/tool"

func Api_setting_external_post(config tool.Config, form map[string]string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "setting_external", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if form["recaptcha_ver"] == "" {
		form["recaptcha_ver"] = "altcha_high"
	}
	setting_save_fields(db, setting_external_fields(), form)
	tool.Do_insert_auth_history(db, config.IP, "edit_set (external)")
	return_data["response"] = "ok"
	return return_data
}
