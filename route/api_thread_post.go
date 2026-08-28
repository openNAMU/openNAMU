package route

import (
	"strings"

	"opennamu/route/tool"
)

func Api_thread_post(config tool.Config, topic_num string, doc_name string, content string, topic string, title string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	data := strings.ReplaceAll(content, "\r", "")
	if data == "" {
		return_data["response"] = "empty data"
		return return_data
	}
	if !tool.Do_edit_slow_check(db, config, "thread") {
		return_data["response"] = "error"
		return_data["data"] = "slow edit limit"
		return return_data
	}
	if !tool.Do_edit_filter(db, config, "", data) {
		return_data["response"] = "error"
		return_data["data"] = "edit filter (content)"
		return return_data
	}
	if !tool.Do_bbs_max_length_check(db, config, data) {
		return_data["response"] = "error"
		return_data["data"] = "bbs overflow max length"
		return return_data
	}

	name := doc_name
	sub := title
	if topic_num != "0" {
		if !tool.QueryRow_DB(db, "select title from rd where code = ?", []any{&name}, topic_num) {
			return_data["response"] = "not exist"
			return_data["data"] = "thread"
			return return_data
		}
		if !tool.Check_acl(db, name, topic_num, "topic_view", config.IP) || !tool.Check_acl(db, name, topic_num, "topic", config.IP) {
			return_data["response"] = "require auth"
			return return_data
		}
		if sub == "" {
			tool.QueryRow_DB(db, "select sub from rd where code = ?", []any{&sub}, topic_num)
		}
	} else {
		if topic != "" {
			name = topic
		}
		if name == "" {
			name = "Test"
		}
		if sub == "" {
			sub = tool.Get_language(db, "make_new_topic", true)
		}
		if !tool.Check_acl(db, name, "0", "topic_view", config.IP) || !tool.Check_acl(db, name, "0", "topic", config.IP) || !tool.Check_permission(db, "discuss_make_new_thread", config.IP) {
			return_data["response"] = "require auth"
			return return_data
		}
	}

	if !tool.Do_title_length_check(db, name, "document") {
		return_data["response"] = "error"
		return_data["data"] = "title length"
		return return_data
	}
	if !tool.Do_title_length_check(db, sub, "topic") {
		return_data["response"] = "error"
		return_data["data"] = "topic title length"
		return return_data
	}
	if !tool.Do_edit_filter(db, config, "", sub) {
		return_data["response"] = "error"
		return_data["data"] = "edit filter (title)"
		return return_data
	}

	if topic_num == "0" {
		topic_num = thread_next_code(db)
		tool.Exec_DB(
			db,
			"insert into rd (title, sub, code, date, band, stop, agree, acl) values (?, ?, ?, ?, '', '', '', '')",
			name,
			sub,
			topic_num,
			tool.Get_time(),
		)
	}

	id := thread_next_id(db, topic_num)
	thread_add(db, topic_num, id, data, config.IP, "")
	tool.Exec_DB(db, "update rd set date = ? where code = ?", tool.Get_time(), topic_num)
	topic_thread_notify(db, config, topic_num, id, name, sub)
	topic_reference_notify(db, config, data, id, topic_num, "", name, sub, "thread")

	return_data["response"] = "ok"
	return_data["topic_num"] = topic_num
	return_data["comment_num"] = id
	return return_data
}
