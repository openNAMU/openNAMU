package route

import "opennamu/route/tool"

func Api_login_register(config tool.Config, id string, password string, password_check string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	error_name := user_register_validate(db, config, id, password, password_check)
	if error_name != "" {
		return_data["response"] = "error"
		return_data["data"] = error_name
		if error_name == "ban" {
			return_data["ban_type"] = tool.Get_user_auth(db, config.IP)
		}
		return return_data
	}

	Api_add_user(config, id, password, "", "")
	return_data["response"] = "ok"

	return return_data
}
