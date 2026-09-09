package route

import "opennamu/route/tool"

func Api_thread_setting_post(config tool.Config, topic_num string, stop string, agree string, why string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "bbs_setting", config.IP) {
		return map[string]any{"response": "require auth"}
	}
	if !thread_bbs_root_exists(db, topic_num) {
		return map[string]any{"response": "not exist", "data": "thread"}
	}

	if stop != "" && stop != "S" && stop != "O" {
		stop = ""
	}
	if agree != "" {
		agree = "O"
	}

	prefix := "열림"
	if stop == "O" {
		prefix = "닫힘"
	}

	for _, value := range []struct {
		name string
		data string
	}{
		{"prefix", prefix},
	} {
		tool.Exec_DB(db, "delete from bbs_data where set_name = ? and set_id = ? and set_code = ?", value.name, thread_bbs_id, topic_num)
		if value.data != "" {
			tool.Exec_DB(db, "insert into bbs_data (set_name, set_id, set_code, set_data) values (?, ?, ?, ?)", value.name, thread_bbs_id, topic_num, value.data)
		}
	}

	tool.Exec_DB(db, "delete from bbs_data where set_name in ('topic_agree', 'topic_stop', 'comment_close') and set_id = ? and set_code = ?", thread_bbs_id, topic_num)
	date := tool.Get_time()
	tool.Exec_DB(db, "update bbs_data set set_data = ? where set_name = 'date' and set_id = ? and set_code = ?", date, thread_bbs_id, topic_num)
	bbs_post_last_activity_update(db, thread_bbs_id, topic_num, date)
	tool.Search_bbs_index_update(db, thread_bbs_id, topic_num)

	return map[string]any{"response": "ok"}
}
