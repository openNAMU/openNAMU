package route

import "opennamu/route/tool"

func View_setting_top_menu_post(config tool.Config, form map[string]string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	setting_save_value(db, "top_menu", "", setting_form_value(form, "content", ""))
	tool.Do_insert_auth_history(db, config.IP, "edit_set (top_menu)")

	return tool.Get_redirect("/setting/top_menu")
}
