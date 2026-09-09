package route

import (
	"database/sql"
	"strconv"
	"strings"

	"opennamu/route/tool"
)

func api_thread_bbs_post(config tool.Config, topic_num string, doc_name string, content string, topic string, title string) map[string]any {
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
		if !tool.QueryRow_DB(
			db,
			"select set_data from bbs_data where set_name = 'document' and set_id = ? and set_code = ?",
			[]any{&name},
			thread_bbs_id,
			topic_num,
		) {
			return_data["response"] = "not exist"
			return_data["data"] = "thread"
			return return_data
		}
		if sub == "" {
			tool.QueryRow_DB(
				db,
				"select set_data from bbs_data where set_name = 'title' and set_id = ? and set_code = ?",
				[]any{&sub},
				thread_bbs_id,
				topic_num,
			)
		}
		if !tool.Check_acl(db, thread_bbs_id, "", "bbs_comment", config.IP) {
			return_data["response"] = "require auth"
			return return_data
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
		if !tool.Check_acl(db, thread_bbs_id, "", "bbs_edit", config.IP) {
			return_data["response"] = "require auth"
			return return_data
		}
	}

	if !tool.Check_acl(db, name, "", "render", config.IP) {
		return_data["response"] = "require auth"
		return return_data
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
		last_code := ""
		tool.QueryRow_DB(
			db,
			"select set_code from bbs_data where set_name = 'title' and set_id = ? order by set_code + 0 desc limit 1",
			[]any{&last_code},
			thread_bbs_id,
		)
		topic_num = strconv.Itoa(tool.Str_to_int(last_code) + 1)
		if err := thread_bbs_insert_post(db, topic_num, name, sub, data, config.IP); err != nil {
			return_data["response"] = "error"
			return_data["data"] = "thread"
			return return_data
		}
		topic_reference_notify(db, config, data, "1", topic_num, thread_bbs_id, name, sub, "bbs")
		bbs_watch_notify(db, config, thread_bbs_id, topic_num, "1", "thread", sub, config.IP, "")
	} else {
		id := thread_bbs_next_comment(db, topic_num)
		date := tool.Get_time()
		comment_set_id := thread_bbs_id + "-" + topic_num
		for _, value := range [][]string{
			{"comment", data},
			{"comment_date", date},
			{"comment_user_id", config.IP},
		} {
			tool.Exec_DB(
				db,
				"insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)",
				value[0],
				id,
				comment_set_id,
				value[1],
			)
		}
		bbs_post_comment_count_update(db, thread_bbs_id, topic_num, 1)
		tool.Exec_DB(
			db,
			"update bbs_data set set_data = ? where set_name = 'date' and set_id = ? and set_code = ?",
			date,
			thread_bbs_id,
			topic_num,
		)
		tool.Search_bbs_index_update_comment(db, thread_bbs_id, topic_num, id)
		topic_reference_notify(db, config, data, id, topic_num, thread_bbs_id, name, sub, "bbs")
		bbs_watch_notify(db, config, thread_bbs_id, topic_num, id, "thread", sub, thread_bbs_post_user(db, topic_num), "")
		return_data["comment_num"] = id
	}

	return_data["response"] = "ok"
	return_data["topic_num"] = topic_num
	if topic_num == "0" {
		return_data["comment_num"] = "1"
	}
	if _, ok := return_data["comment_num"]; !ok {
		return_data["comment_num"] = "1"
	}
	return return_data
}

func thread_bbs_post_user(db *sql.DB, topic_num string) string {
	user_id := ""
	tool.QueryRow_DB(
		db,
		"select set_data from bbs_data where set_name = 'user_id' and set_id = ? and set_code = ?",
		[]any{&user_id},
		thread_bbs_id,
		topic_num,
	)
	return user_id
}

func thread_bbs_next_comment(db *sql.DB, topic_num string) string {
	last_code := ""
	tool.QueryRow_DB(
		db,
		"select set_code from bbs_data where set_name = 'comment' and set_id = ? order by set_code + 0 desc limit 1",
		[]any{&last_code},
		thread_bbs_id+"-"+topic_num,
	)
	return strconv.Itoa(tool.Str_to_int(last_code) + 1)
}

func thread_bbs_insert_post(db *sql.DB, topic_num string, name string, sub string, data string, user_id string) error {
	date := tool.Get_time()
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for _, value := range [][]string{
		{"title", sub},
		{"data", ""},
		{"date", date},
		{"user_id", user_id},
		{"comment_count", "1"},
		{"document", name},
		{"tag", name},
		{"topic_source", "bbs"},
		{"prefix", "열림"},
	} {
		if _, err = tx.Exec(
			tool.DB_change("insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)"),
			value[0],
			topic_num,
			thread_bbs_id,
			value[1],
		); err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	for _, value := range [][]string{
		{"comment", data},
		{"comment_date", date},
		{"comment_user_id", user_id},
	} {
		if _, err = tx.Exec(
			tool.DB_change("insert into bbs_data (set_name, set_code, set_id, set_data) values (?, '1', ?, ?)"),
			value[0],
			thread_bbs_id+"-"+topic_num,
			value[1],
		); err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	if err = tx.Commit(); err != nil {
		return err
	}

	tool.Search_bbs_index_update(db, thread_bbs_id, topic_num)
	tool.Search_bbs_index_update_comment(db, thread_bbs_id, topic_num, "1")
	return nil
}
