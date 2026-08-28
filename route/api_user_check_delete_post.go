package route

import "opennamu/route/tool"

func Api_user_check_delete_post(config tool.Config, user_name string, user_ip string, today string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	tool.Exec_DB(db, "delete from ua_d where name = ? and ip = ? and today = ?", user_name, user_ip, today)
	tool.Do_insert_auth_history(db, config.IP, "user_check_delete ("+user_name+")")
	return_data["response"] = "ok"
	return return_data
}
