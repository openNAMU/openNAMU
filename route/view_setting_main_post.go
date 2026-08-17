package route

import "opennamu/route/tool"

func View_setting_main_post(config tool.Config, form map[string]string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	fields := setting_main_fields()
	setting_save_fields(db, fields, form)
	tool.Do_insert_auth_history(db, config.IP, "edit_set (main)")

	return tool.Get_redirect("/setting/main")
}
