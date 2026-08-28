package route

import "opennamu/route/tool"

func Api_user_head_reset_post(config tool.Config, skin_name string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if tool.IP_or_user(config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	user_save(db, config.IP, "custom_css", "")
	user_save(db, config.IP, "custom_css_"+skin_name, "")
	user_delete(db, config.IP, "head")
	return_data["response"] = "ok"
	return return_data
}
