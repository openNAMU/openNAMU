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
	error_name, invite_hash := user_register_invite(db, config, "")
	if error_name != "" {
		return_data["response"] = "error"
		return_data["data"] = error_name
		return return_data
	}

	result := Api_add_user_invite(config, id, password, "", "", invite_hash)
	if result["response"] != "ok" {
		return_data["response"] = "error"
		return_data["data"] = result["data"]
		return return_data
	}
	return_data["response"] = "ok"

	return return_data
}
