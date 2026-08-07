package route

import (
	"opennamu/route/tool"
)

func Api_func_alarm_post(config tool.Config, from string, to string, data string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	tool.Send_alarm(db, from, to, data)

	return_data := make(map[string]any)
	return_data["response"] = "ok"

	return return_data
}
