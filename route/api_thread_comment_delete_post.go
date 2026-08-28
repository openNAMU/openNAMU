package route

import "opennamu/route/tool"

func Api_thread_comment_delete_post(config tool.Config, topic_num string, comment_num string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "thread_comment_delete", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	exists_id := ""
	if !tool.QueryRow_DB(db, "select id from topic where code = ? and id = ?", []any{&exists_id}, topic_num, comment_num) {
		return_data["response"] = "not exist"
		return_data["data"] = "comment"
		return return_data
	}
	tool.Exec_DB(db, "delete from topic where code = ? and id = ?", topic_num, comment_num)
	tool.Exec_DB(db, "update rd set date = ? where code = ?", tool.Get_time(), topic_num)
	tool.Do_insert_auth_history(db, config.IP, "delete_topic_comment (code "+topic_num+"#"+comment_num+")")

	return_data["response"] = "ok"
	return return_data
}
