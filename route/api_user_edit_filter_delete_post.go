package route

import "opennamu/route/tool"

func Api_user_edit_filter_delete_post(config tool.Config, user_name string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	tool.Exec_DB(db, "delete from user_set where name = 'edit_filter' and id = ?", user_name)
	return_data["response"] = "ok"
	return return_data
}
