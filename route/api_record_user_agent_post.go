package route

import "opennamu/route/tool"

func Api_record_user_agent_post(config tool.Config, user_name string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	tool.Record_user_agent(db, user_name, config.IP, config.UserAgent, tool.Get_time())

	return_data := make(map[string]any)
	return_data["response"] = "ok"
	return return_data
}
