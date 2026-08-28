package route

import "opennamu/route/tool"

func Api_user_top_menu_post(config tool.Config, content string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if tool.IP_or_user(config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	user_save(db, config.IP, "top_menu", content)
	return_data["response"] = "ok"
	return return_data
}
