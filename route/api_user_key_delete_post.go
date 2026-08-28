package route

import "opennamu/route/tool"

func Api_user_key_delete_post(config tool.Config) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if tool.IP_or_user(config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	user_delete(db, config.IP, "random_key")
	return_data["response"] = "ok"
	return return_data
}
