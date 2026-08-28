package route

import "opennamu/route/tool"

func Api_thread_change_post(config tool.Config, topic_num string, new_title string, new_sub string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "thread_change", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	title := ""
	sub := ""
	if !tool.QueryRow_DB(db, "select title, sub from rd where code = ?", []any{&title, &sub}, topic_num) {
		return_data["response"] = "not exist"
		return_data["data"] = "thread"
		return return_data
	}
	if new_title == "" {
		new_title = title
	}
	if new_sub == "" {
		new_sub = sub
	}
	if !tool.Do_title_length_check(db, new_title, "document") {
		return_data["response"] = "error"
		return_data["data"] = "title length"
		return return_data
	}
	if !tool.Do_title_length_check(db, new_sub, "topic") {
		return_data["response"] = "error"
		return_data["data"] = "topic title length"
		return return_data
	}

	tool.Exec_DB(db, "update rd set title = ?, sub = ?, date = ? where code = ?", new_title, new_sub, tool.Get_time(), topic_num)
	thread_add(db, topic_num, thread_next_id(db, topic_num), tool.Get_language(db, "topic_name_change", true)+" : "+sub+" ("+title+") → "+new_sub+" ("+new_title+")", config.IP, "1")
	tool.Do_insert_auth_history(db, config.IP, "change_topic_name (code "+topic_num+")")

	return_data["response"] = "ok"
	return return_data
}
