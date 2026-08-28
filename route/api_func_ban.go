package route

import (
	"opennamu/route/tool"
)

func Api_func_ban(config tool.Config, ip string, ban_type string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if ip == "" {
		ip = config.IP
	}

	auth_name := tool.Get_user_auth(db, ip)
	auth_info := tool.Get_auth_info(db, ip)
	blocked := tool.Auth_group_name_ban(auth_name)
	if ban_type == "login" {
		blocked = !auth_info["login_available"]
	} else if ban_type == "register" {
		blocked = !auth_info["register_available"]
	}

	new_data := make(map[string]any)
	new_data["response"] = "ok"
	new_data["ban"] = ""
	new_data["ban_type"] = ""
	if blocked {
		new_data["ban"] = "true"
		new_data["ban_type"] = auth_name
	}

	return new_data
}
