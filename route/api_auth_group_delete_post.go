package route

import "opennamu/route/tool"

func Api_auth_group_delete_post(config tool.Config, name string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if tool.Auth_group_name_default(name) {
		return_data["response"] = "default"
		return return_data
	}
	if !tool.Check_permission(db, "auth_group_manage", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if tool.Auth_group_in_use(db, name) {
		return_data["response"] = "error"
		return return_data
	}
	tool.Exec_DB(db, "delete from alist where name = ?", name)
	tool.Do_insert_auth_history(db, config.IP, "auth_group_delete ("+name+")")

	return_data["response"] = "ok"
	return return_data
}
