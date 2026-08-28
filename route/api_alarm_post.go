package route

import "opennamu/route/tool"

func Api_alarm_delete_post(config tool.Config, user_name string, id string, all bool) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if user_name == "" {
		user_name = config.IP
	}
	if user_name != config.IP && !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if all {
		tool.Exec_DB(db, "delete from user_notice where name = ?", user_name)
	} else if id != "" {
		tool.Exec_DB(db, "delete from user_notice where name = ? and id = ?", user_name, id)
	}
	return_data["response"] = "ok"
	return return_data
}

func Api_alarm_read_post(config tool.Config, user_name string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if user_name == "" {
		user_name = config.IP
	}
	if user_name != config.IP && !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	tool.Exec_DB(db, "update user_notice set readme = '1' where name = ?", user_name)
	return_data["response"] = "ok"
	return return_data
}
