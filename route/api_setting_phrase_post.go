package route

import "opennamu/route/tool"

func Api_setting_phrase_post(config tool.Config, form map[string]string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "setting_phrase", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	for _, field := range setting_phrase_fields() {
		setting_save_value(db, field.name, "", setting_form_value(form, field.name, ""))
	}
	tool.Do_insert_auth_history(db, config.IP, "edit_set (phrase)")
	return_data["response"] = "ok"
	return return_data
}
