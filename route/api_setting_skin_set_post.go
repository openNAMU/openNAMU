package route

import "opennamu/route/tool"

func Api_setting_skin_set_post(config tool.Config, form map[string]string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "setting_skin", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	fields, set_list := setting_skin_fields(db)
	for _, field := range fields {
		default_value := ""
		if choices := set_list[field.name]; len(choices) > 0 {
			default_value = choices[0][0]
		}
		setting_save_value(db, field.name, "", setting_form_value(form, field.name, default_value))
	}
	tool.Do_insert_auth_history(db, config.IP, "edit_set (skin_set)")
	return_data["response"] = "ok"
	return return_data
}
