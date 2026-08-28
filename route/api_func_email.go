package route

import (
	"opennamu/route/tool"
)

func Api_func_email_post(config tool.Config, who string, title string, data string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return map[string]any{"response": "require auth"}
	}

	err := tool.Send_email(db, config.IP, who, title, data)
	if err == nil {
		new_data := make(map[string]any)
		new_data["response"] = "ok"

		return new_data
	}

	new_data := make(map[string]any)
	new_data["response"] = "err"
	new_data["data"] = err.Error()

	return new_data
}
