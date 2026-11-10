package route

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"opennamu/route/tool"
)

func Api_thread_bbs_post(config tool.Config, topic_num string, doc_name string, content string, topic string, title string) map[string]any {
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
		// Board-level slow mode applies to new discussion topics. The check and
		// the insert share one transaction so a denied write records no cooldown
		// and no content.
		board_interval := tool.Bbs_board_slow_interval(db, thread_bbs_id)
		retry_after := 0

		if err := Bbs_identifier_transaction(db, thread_bbs_id, "", true, func(tx *sql.Tx) error {
			remaining, err := tool.Bbs_slow_enforce(db, tx, config.IP, thread_bbs_id, "", board_interval, 0)
			if err != nil {
				return err
			}
			if remaining > 0 {
				retry_after = remaining
				return tool.BbsSlowModeError
			}
			last_code := ""
			tool.QueryRow_DB(
				tx,
				"select set_code from bbs_data where set_name = 'title' and set_id = ? order by set_code + 0 desc limit 1",
				[]any{&last_code},
				thread_bbs_id,
			)
			topic_num = strconv.Itoa(tool.Str_to_int(last_code) + 1)
			return Thread_bbs_insert_post_tx(tx, topic_num, name, sub, data, config.IP)
		}); err != nil {
			if errors.Is(err, tool.BbsSlowModeError) {
				return_data["response"] = "error"
				return_data["data"] = "slow_mode"
				return_data["retry_after_seconds"] = retry_after
				return return_data
			}
			return_data["response"] = "error"
			return_data["data"] = "thread"
			return return_data
		}
		tool.Search_bbs_index_update(db, thread_bbs_id, topic_num)
		tool.Search_bbs_index_update_comment(db, thread_bbs_id, topic_num, "1")
		Topic_reference_notify(db, config, data, "1", topic_num, thread_bbs_id, name, sub, "bbs")
		Bbs_watch_notify(db, config, thread_bbs_id, topic_num, "1", "thread", sub, config.IP, "")
		Bbs_post_watch_notify(db, config, thread_bbs_id, topic_num, "thread", sub, config.IP, Bbs_tag_list(name))
	} else {
		id := ""
		date := tool.Get_time()
		comment_set_id := thread_bbs_id + "-" + topic_num

		// Board-level and this topic's post-level slow mode both apply to a
		// comment/reply. The post scope keys on the root topic; nested replies
		// share that cooldown and never touch another topic's.
		board_interval := tool.Bbs_board_slow_interval(db, thread_bbs_id)
		post_interval := tool.Bbs_post_slow_interval(db, thread_bbs_id, topic_num)
		retry_after := 0

		if err := Bbs_identifier_transaction(db, thread_bbs_id, topic_num, board_interval > 0, func(tx *sql.Tx) error {
			remaining, err := tool.Bbs_slow_enforce(db, tx, config.IP, thread_bbs_id, topic_num, board_interval, post_interval)
			if err != nil {
				return err
			}
			if remaining > 0 {
				retry_after = remaining
				return tool.BbsSlowModeError
			}
			id = Thread_bbs_next_comment(tx, topic_num)
			for _, value := range [][]string{
				{"comment", data},
				{"comment_date", date},
				{"comment_user_id", config.IP},
			} {
				if _, err := tx.Exec(
					tool.DB_change("insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)"),
					value[0],
					id,
					comment_set_id,
					value[1],
				); err != nil {
					return err
				}
			}
			Bbs_post_comment_count_update(tx, thread_bbs_id, topic_num, 1)
			Bbs_post_last_activity_update(tx, thread_bbs_id, topic_num, date)
			_, execErr := tx.Exec(
				tool.DB_change("update bbs_data set set_data = ? where set_name = 'date' and set_id = ? and set_code = ?"),
				date,
				thread_bbs_id,
				topic_num,
			)
			return execErr
		}); err != nil {
			if errors.Is(err, tool.BbsSlowModeError) {
				return_data["response"] = "error"
				return_data["data"] = "slow_mode"
				return_data["retry_after_seconds"] = retry_after
				return return_data
			}
			panic(err)
		}
		tool.Search_bbs_index_update_comment(db, thread_bbs_id, topic_num, id)
		Topic_reference_notify(db, config, data, id, topic_num, thread_bbs_id, name, sub, "bbs")
		Bbs_watch_notify(db, config, thread_bbs_id, topic_num, id, "thread", sub, Thread_bbs_post_user(db, topic_num), "")
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

func Thread_bbs_post_user(db *sql.DB, topic_num string) string {
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

func Thread_bbs_next_comment(db tool.DB_runner, topic_num string) string {
	last_code := ""
	tool.QueryRow_DB(
		db,
		"select set_code from bbs_data where set_name = 'comment' and set_id = ? order by set_code + 0 desc limit 1",
		[]any{&last_code},
		thread_bbs_id+"-"+topic_num,
	)
	return strconv.Itoa(tool.Str_to_int(last_code) + 1)
}

// Thread_bbs_insert_post inserts a new topic post and its first comment in its
// own transaction, preserving the existing helper interface.
func Thread_bbs_insert_post(db *sql.DB, topic_num string, name string, sub string, data string, user_id string) error {
	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		return Thread_bbs_insert_post_tx(tx, topic_num, name, sub, data, user_id)
	}); err != nil {
		return err
	}
	tool.Search_bbs_index_update(db, thread_bbs_id, topic_num)
	tool.Search_bbs_index_update_comment(db, thread_bbs_id, topic_num, "1")
	return nil
}

// Thread_bbs_insert_post_tx inserts a new topic post and its first comment into
// a caller-owned transaction. The caller owns Begin/Commit/rollback.
func Thread_bbs_insert_post_tx(tx *sql.Tx, topic_num string, name string, sub string, data string, user_id string) error {
	date := tool.Get_time()
	for _, value := range [][]string{
		{"title", sub},
		{"data", ""},
		{"date", date},
		{"last_activity", date},
		{"user_id", user_id},
		{"comment_count", "1"},
		{"document", name},
		{"tag", name},
		{"topic_source", "bbs"},
		{"prefix", "열림"},
	} {
		if _, err := tx.Exec(
			tool.DB_change("insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)"),
			value[0],
			topic_num,
			thread_bbs_id,
			value[1],
		); err != nil {
			return err
		}
	}
	for _, value := range [][]string{
		{"comment", data},
		{"comment_date", date},
		{"comment_user_id", user_id},
	} {
		if _, err := tx.Exec(
			tool.DB_change("insert into bbs_data (set_name, set_code, set_id, set_data) values (?, '1', ?, ?)"),
			value[0],
			thread_bbs_id+"-"+topic_num,
			value[1],
		); err != nil {
			return err
		}
	}
	return nil
}
