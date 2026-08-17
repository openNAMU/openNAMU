package route

import "opennamu/route/tool"

func View_setting_external_post(config tool.Config, form map[string]string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	setting_save_fields(db, setting_external_fields(), form)
	tool.Do_insert_auth_history(db, config.IP, "edit_set (external)")

	return tool.Get_redirect("/setting/external")
}
