package route

import "opennamu/route/tool"

func Api_thread_delete_post(config tool.Config, topic_num string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	title := ""
	if !tool.QueryRow_DB(db, "select title from rd where code = ?", []any{&title}, topic_num) {
		return_data["response"] = "not exist"
		return_data["data"] = "thread"
		return return_data
	}

	tool.Do_insert_auth_history(db, config.IP, "delete_topic (code "+topic_num+")")
	tool.Exec_DB(db, "delete from topic where code = ?", topic_num)
	tool.Exec_DB(db, "delete from topic_set where thread_code = ?", topic_num)
	tool.Exec_DB(db, "delete from rd where code = ?", topic_num)

	return_data["response"] = "ok"
	return_data["data"] = title
	return return_data
}
