package route

import "opennamu/route/tool"

func Api_history_reset_post(config tool.Config, doc_name string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "history_manage", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if doc_name == "" {
		return_data["response"] = "error"
		return_data["data"] = "invalid document"
		return return_data
	}

	tool.Exec_DB(db, "delete from history where title = ?", doc_name)
	tool.Do_insert_auth_history(db, config.IP, "history_reset ("+doc_name+")")

	return_data["response"] = "ok"
	return return_data
}
