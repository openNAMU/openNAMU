package route

import "opennamu/route/tool"

func View_setting_main_logo_post(config tool.Config, form map[string]string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
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

	return tool.Get_redirect("/setting/main/logo")
}
