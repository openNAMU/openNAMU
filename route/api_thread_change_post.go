package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func Api_thread_change_post(config tool.Config, topic_num string, new_title string, new_sub string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "bbs_post_manage", config.IP) {
		return map[string]any{"response": "require auth"}
	}
	if !thread_bbs_root_exists(db, topic_num) {
		return map[string]any{"response": "not exist", "data": "thread"}
	}

	old_title := ""
	old_sub := ""
	tool.QueryRow_DB(db, "select set_data from bbs_data where set_name = 'document' and set_id = ? and set_code = ?", []any{&old_title}, thread_bbs_id, topic_num)
	tool.QueryRow_DB(db, "select set_data from bbs_data where set_name = 'title' and set_id = ? and set_code = ?", []any{&old_sub}, thread_bbs_id, topic_num)
	if new_title == "" {
		new_title = old_title
	}
	if new_sub == "" {
		new_sub = old_sub
	}
	if !tool.Do_title_length_check(db, new_title, "document") || !tool.Do_title_length_check(db, new_sub, "topic") {
		return map[string]any{"response": "error", "data": "title length"}
	}

	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		for _, update := range []struct {
			query  string
			values []any
		}{
			{"update bbs_data set set_data = ? where set_name = 'document' and set_id = ? and set_code = ?", []any{new_title, thread_bbs_id, topic_num}},
			{"update bbs_data set set_data = ? where set_name = 'tag' and set_id = ? and set_code = ? and set_data = ?", []any{new_title, thread_bbs_id, topic_num, old_title}},
			{"update bbs_data set set_data = ? where set_name = 'title' and set_id = ? and set_code = ?", []any{new_sub, thread_bbs_id, topic_num}},
			{"update bbs_data set set_data = ? where set_name = 'date' and set_id = ? and set_code = ?", []any{tool.Get_time(), thread_bbs_id, topic_num}},
		} {
			if _, err := tx.Exec(tool.DB_change(update.query), update.values...); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}

	return map[string]any{"response": "ok"}
}
