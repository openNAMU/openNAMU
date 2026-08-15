package route

import "opennamu/route/tool"

func Api_login_login(config tool.Config, id string, password string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)

	ban_data := tool.Get_user_ban(db, config.IP, "login")
	if ban_data[0] == "true" {
		return_data["response"] = "error"
		return_data["data"] = "ban"
		return_data["ban_type"] = ban_data[1]

		return return_data
	}

	if !tool.Password_check(db, id, password) {
		return_data["response"] = "error"
		return_data["data"] = "password error"

		return return_data
	}

	config.Session.Set("id", id)

	if err := config.Session.Save(); err != nil {
		return_data["response"] = "error"
		return_data["data"] = "session save error"

		return return_data
	}
	tool.Record_user_agent(db, id, config.IP, config.UserAgent, tool.Get_time())

	return_data["response"] = "ok"

	return return_data
}
