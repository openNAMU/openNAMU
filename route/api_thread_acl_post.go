package route

import "opennamu/route/tool"

func Api_thread_acl_post(config tool.Config, topic_num string, acl string, acl_view string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "thread_acl", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	thread_acl := ""
	if !tool.QueryRow_DB(db, "select acl from rd where code = ?", []any{&thread_acl}, topic_num) {
		return_data["response"] = "not exist"
		return_data["data"] = "thread"
		return return_data
	}
	if !acl_value_valid(db, acl) || !acl_value_valid(db, acl_view) {
		return_data["response"] = "error"
		return_data["data"] = "error"
		return return_data
	}

	tool.Exec_DB(db, "update rd set acl = ?, date = ? where code = ?", acl, tool.Get_time(), topic_num)
	tool.Exec_DB(db, "delete from topic_set where thread_code = ? and set_name = 'thread_view_acl'", topic_num)
	tool.Exec_DB(db, "insert into topic_set (thread_code, set_name, set_id, set_data) values (?, 'thread_view_acl', '1', ?)", topic_num, acl_view)
	thread_add(db, topic_num, thread_next_id(db, topic_num), tool.Get_language(db, "acl_thread_change", true)+" : "+acl, config.IP, "1")
	tool.Do_insert_auth_history(db, config.IP, "change_topic_acl (code "+topic_num+")")

	return_data["response"] = "ok"
	return return_data
}
