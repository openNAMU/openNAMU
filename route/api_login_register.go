package route

import "opennamu/route/tool"

func Api_login_register(config tool.Config, id string, password string, password_check string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	error_name := User_register_validate(db, config, id, password, password_check)
	if error_name != "" {
		return_data["response"] = "error"
		return_data["data"] = error_name
		if error_name == "ban" {
			return_data["ban_type"] = tool.Get_user_auth(db, config.IP)
		}
		return return_data
	}
	invite_error, _ := User_register_invite(db, config, "")
	if invite_error != "" {
		return_data["response"] = "error"
		return_data["data"] = invite_error
		return return_data
	}

	result := Api_add_user(config, id, password, "", "")
	if result["response"] != "ok" {
		return_data["response"] = "error"
		return_data["data"] = result["data"]
		return return_data
	}
	return_data["response"] = "ok"

	return return_data
}
