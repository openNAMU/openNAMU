package route

import "opennamu/route/tool"

func View_setting_phrase_post(config tool.Config, form map[string]string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	fields := setting_phrase_fields()
	for _, field := range fields {
		setting_save_value(db, field.name, "", setting_form_value(form, field.name, ""))
	}
	tool.Do_insert_auth_history(db, config.IP, "edit_set (phrase)")

	return tool.Get_redirect("/setting/phrase")
}
