package route

import "opennamu/route/tool"

func Api_thread_comment_blind_post(config tool.Config, topic_num string, comment_num string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "bbs_comment_manage", config.IP) {
		return map[string]any{"response": "require auth"}
	}
	comment_set_id := thread_bbs_id + "-" + topic_num
	comment := ""
	if !tool.QueryRow_DB(db, "select set_data from bbs_data where set_name = 'comment' and set_id = ? and set_code = ?", []any{&comment}, comment_set_id, comment_num) {
		return map[string]any{"response": "not exist", "data": "comment"}
	}

	blind := ""
	tool.QueryRow_DB(db, "select set_data from bbs_data where set_name = 'blind' and set_id = ? and set_code = ?", []any{&blind}, comment_set_id, comment_num)
	if blind == "O" {
		tool.Exec_DB(db, "delete from bbs_data where set_name = 'blind' and set_id = ? and set_code = ?", comment_set_id, comment_num)
		blind = ""
	} else {
		if blind == "" {
			tool.Exec_DB(db, "insert into bbs_data (set_name, set_id, set_code, set_data) values ('blind', ?, ?, 'O')", comment_set_id, comment_num)
		} else {
			tool.Exec_DB(db, "update bbs_data set set_data = 'O' where set_name = 'blind' and set_id = ? and set_code = ?", comment_set_id, comment_num)
		}
		blind = "O"
	}
	return map[string]any{"response": "ok", "data": blind}
}
