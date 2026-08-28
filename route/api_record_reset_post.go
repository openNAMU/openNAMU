package route

import "opennamu/route/tool"

func Api_record_reset_post(config tool.Config, user_name string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if user_name == "" {
		return_data["response"] = "error"
		return return_data
	}
	tool.Exec_DB(db, "delete from history where ip = ?", user_name)
	tool.Do_insert_auth_history(db, config.IP, "record_reset ("+user_name+")")
	return_data["response"] = "ok"
	return return_data
}
