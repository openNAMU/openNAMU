package route

import (
	"strings"

	"opennamu/route/tool"
)

func Api_thread_setting_post(config tool.Config, topic_num string, stop string, agree string, why string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "thread_setting", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	name := ""
	old_stop := ""
	old_agree := ""
	if !tool.QueryRow_DB(db, "select title, stop, agree from rd where code = ?", []any{&name, &old_stop, &old_agree}, topic_num) {
		return_data["response"] = "not exist"
		return_data["data"] = "thread"
		return return_data
	}

	new_stop := stop
	if new_stop != "" && new_stop != "S" && new_stop != "O" {
		new_stop = ""
	}
	new_agree := ""
	if agree != "" {
		new_agree = "O"
	}
	tool.Exec_DB(db, "update rd set stop = ?, agree = ? where code = ?", new_stop, new_agree, topic_num)

	if old_stop != new_stop {
		state_key := "topic_state_change_normal"
		if new_stop == "S" {
			state_key = "topic_state_change_stop"
		} else if new_stop == "O" {
			state_key = "topic_state_change_close"
		}
		thread_add(db, topic_num, thread_next_id(db, topic_num), tool.Get_language(db, state_key, true), config.IP, "1")
	}
	if old_agree != new_agree {
		state_key := "topic_state_change_disagree"
		if new_agree == "O" {
			state_key = "topic_state_change_agree"
		}
		thread_add(db, topic_num, thread_next_id(db, topic_num), tool.Get_language(db, state_key, true), config.IP, "1")
	}
	if why = strings.TrimSpace(why); why != "" {
		thread_add(db, topic_num, thread_next_id(db, topic_num), tool.Get_language(db, "why", true)+" : "+why, config.IP, "1")
	}
	tool.Exec_DB(db, "update rd set date = ? where code = ?", tool.Get_time(), topic_num)
	tool.Do_insert_auth_history(db, config.IP, "change_topic_set (code "+topic_num+")")

	return_data["response"] = "ok"
	return return_data
}
