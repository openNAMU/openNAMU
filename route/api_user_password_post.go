package route

import "opennamu/route/tool"

func Api_user_password_post(config tool.Config, current_password string, password string, password_repeat string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if tool.IP_or_user(config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if password == "" {
		return_data["response"] = "error"
		return_data["data"] = "password empty"
		return return_data
	}
	if password != password_repeat {
		return_data["response"] = "error"
		return_data["data"] = "password different"
		return return_data
	}
	minimum := user_other(db, "password_min_length")
	if minimum != "" && tool.Get_len(password) < tool.Str_to_int(minimum) {
		return_data["response"] = "error"
		return_data["data"] = "password too short"
		return return_data
	}
	if !tool.Password_check(db, config.IP, current_password) {
		return_data["response"] = "error"
		return_data["data"] = "password error"
		return return_data
	}
	user_save(db, config.IP, "pw", tool.Password_encode(db, password, tool.Get_user_encode(db, config.IP)))
	return_data["response"] = "ok"
	return return_data
}
