package route

import "opennamu/route/tool"

func Api_thread_comment_notice_post(config tool.Config, topic_num string, comment_num string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "thread_comment_manage", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	top := ""
	if !tool.QueryRow_DB(db, "select top from topic where code = ? and id = ?", []any{&top}, topic_num, comment_num) {
		return_data["response"] = "not exist"
		return_data["data"] = "comment"
		return return_data
	}
	if top == "O" {
		top = ""
	} else {
		top = "O"
	}
	tool.Exec_DB(db, "update topic set top = ? where code = ? and id = ?", top, topic_num, comment_num)
	tool.Exec_DB(db, "update rd set date = ? where code = ?", tool.Get_time(), topic_num)
	tool.Do_insert_auth_history(db, config.IP, "notice (code "+topic_num+"#"+comment_num+")")

	return_data["response"] = "ok"
	return_data["data"] = top
	return return_data
}
