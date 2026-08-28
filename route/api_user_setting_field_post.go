package route

import "opennamu/route/tool"

func Api_user_setting_field_post(config tool.Config, field string, value string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if tool.IP_or_user(config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if field == "user_name" && !tool.Get_user_name_check(db, value) && value != user_value(db, config.IP, "user_name") {
		return_data["response"] = "error"
		return_data["data"] = "user name error"
		return return_data
	}
	if field == "email" && value == "" {
		user_delete(db, config.IP, field)
	} else {
		user_save(db, config.IP, field, value)
	}
	return_data["response"] = "ok"
	return return_data
}
