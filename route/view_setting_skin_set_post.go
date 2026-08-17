package route

import "opennamu/route/tool"

func View_setting_skin_set_post(config tool.Config, form map[string]string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
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

	return tool.Get_redirect("/setting/skin_set")
}
