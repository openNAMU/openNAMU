package route

import "opennamu/route/tool"

func Api_setting_sitemap_set_post(config tool.Config, form map[string]string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "setting_sitemap", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	setting_save_fields(db, setting_sitemap_fields(), form)
	sync_indexnow_key(db)
	tool.Do_insert_auth_history(db, config.IP, "edit_set (sitemap)")
	return_data["response"] = "ok"
	return return_data
}
