package route

import "opennamu/route/tool"

func Api_setting_main_logo_post(config tool.Config, form map[string]string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "setting_main_logo", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	for _, skin := range setting_logo_skins() {
		coverage := ""
		field_name := "main_css"
		if skin != "default" {
			coverage = skin
			field_name = skin
		}
		setting_save_value(db, "logo", coverage, setting_form_value(form, field_name, ""))
	}
	tool.Do_insert_auth_history(db, config.IP, "edit_set (logo)")
	return_data["response"] = "ok"
	return return_data
}
