package route

import (
	"database/sql"

	"opennamu/route/tool"
	"strconv"
	"strings"
)

func Api_bbs_w_edit_post(config tool.Config, set_id string, set_code string, comment_code string, title string, data string, prefix string, tags string, document string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)

	bbs_name := ""
	if !tool.QueryRow_DB(
		db,
		"select set_data from bbs_set where set_id = ? and set_name = 'bbs_name'",
		[]any{&bbs_name},
		set_id,
	) {
		return_data["response"] = "not exist"
		return_data["data"] = "bbs"

		return return_data
	}

	if !bbs_post_blind_allowed(db, set_id, set_code, config.IP, nil) {
		return_data["response"] = "require auth"
		return return_data
	}

	edit_acl := "bbs_edit"
	if comment_code != "" {
		edit_acl = "bbs_comment"
	}
	if !tool.Check_acl(db, set_id, "", edit_acl, config.IP) {
		return_data["response"] = "require auth"

		return return_data
	}
	if set_id == "0" && set_code == "" && comment_code == "" {
		return_data["response"] = "not exist"
		return_data["data"] = "bbs"

		return return_data
	}
	if document != "" {
		if set_id != thread_bbs_id || set_code != "" || comment_code != "" {
			return_data["response"] = "not exist"
			return_data["data"] = "bbs"

			return return_data
		}
		if !tool.Check_acl(db, document, "", "render", config.IP) {
			return_data["response"] = "require auth"

			return return_data
		}
	}

	data = strings.ReplaceAll(data, "\r", "")
	if title == "" {
		title = "test"
	}
	if data == "" {
		return_data["response"] = "error"
		return_data["data"] = "empty data"

		return return_data
	}

	if !tool.Do_edit_filter(db, config, "", title) || !tool.Do_edit_filter(db, config, "", data) {
		return_data["response"] = "error"
		return_data["data"] = "edit filter (content)"

		return return_data
	}
	if !tool.Do_bbs_max_length_check(db, config, data) {
		return_data["response"] = "error"
		return_data["data"] = "bbs overflow max length"

		return return_data
	}

	if comment_code != "" {
		comment_code_split := strings.Split(comment_code, "-")
		comment_set_id := set_id + "-" + set_code
		comment_set_code := ""

		if len(comment_code_split) > 0 {
			comment_set_code = comment_code_split[len(comment_code_split)-1]
			if len(comment_code_split) > 1 {
				comment_set_id += "-" + strings.Join(comment_code_split[:len(comment_code_split)-1], "-")
			}
		}

		if comment_set_code == "" {
			return_data["response"] = "not exist"
			return_data["data"] = "comment"

			return return_data
		}

		comment := ""
		if !tool.QueryRow_DB(
			db,
			"select set_data from bbs_data where set_name = 'comment' and set_id = ? and set_code = ?",
			[]any{&comment},
			comment_set_id,
			comment_set_code,
		) {
			return_data["response"] = "not exist"
			return_data["data"] = "comment"

			return return_data
		}

		comment_user_id := ""
		if !tool.QueryRow_DB(
			db,
			"select set_data from bbs_data where set_name = 'comment_user_id' and set_id = ? and set_code = ?",
			[]any{&comment_user_id},
			comment_set_id,
			comment_set_code,
		) {
			return_data["response"] = "not exist"
			return_data["data"] = "comment"

			return return_data
		}

		comment_manage := tool.Check_permission(db, "bbs_comment_manage", config.IP)
		blind := ""
		tool.QueryRow_DB(
			db,
			"select set_data from bbs_data where set_name = 'blind' and set_id = ? and set_code = ?",
			[]any{&blind},
			comment_set_id,
			comment_set_code,
		)
		if blind == "O" && !comment_manage {
			return_data["response"] = "require auth"
			return return_data
		}
		if comment_user_id != config.IP && !comment_manage {
			return_data["response"] = "require auth"

			return return_data
		}
		if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
			if _, err := tx.Exec(
				tool.DB_change("update bbs_data set set_data = ? where set_name = 'comment' and set_code = ? and set_id = ?"),
				data,
				comment_set_code,
				comment_set_id,
			); err != nil {
				return err
			}
			bbs_post_last_activity_update(tx, set_id, set_code, tool.Get_time())
			return nil
		}); err != nil {
			panic(err)
		}
		tool.Search_bbs_index_update_comment(db, set_id, set_code, comment_code)

		return_data["response"] = "ok"
		return_data["data"] = set_code

		return return_data
	}

	prefix = bbs_prefix_check(db, set_id, prefix)
	if set_id == thread_bbs_id && prefix == "" {
		prefix = "열림"
	}
	if tool.Get_len(title) > bbs_title_max_length {
		return_data["response"] = "error"
		return_data["data"] = "bbs title length"

		return return_data
	}

	if tool.Get_len(tags) > bbs_tag_max_length {
		return_data["response"] = "error"
		return_data["data"] = "bbs tag length"

		return return_data
	}
	tag_list := bbs_tag_list(tags)
	if document != "" && !tool.Arr_in_str(tag_list, document) {
		if tool.Get_len(document) > bbs_tag_max_length {
			return_data["response"] = "error"
			return_data["data"] = "bbs tag length"

			return return_data
		}
		tag_list = append(tag_list, document)
	}

	if set_code == "" {
		if !tool.Check_daily_limit(db, config.IP, "bbs_edit") {
			return_data["response"] = "error"
			return_data["data"] = "daily limit"
			return return_data
		}

		if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
			last_code := ""
			tool.QueryRow_DB(
				tx,
				"select set_code from bbs_data where set_name = 'title' and set_id = ? order by set_code + 0 desc",
				[]any{&last_code},
				set_id,
			)
			set_code = strconv.Itoa(tool.Str_to_int(last_code) + 1)
			date := tool.Get_time()

			insert_db := [][]string{
				{"title", title},
				{"data", data},
				{"date", date},
				{"last_activity", date},
				{"user_id", config.IP},
				{"comment_count", "0"},
			}
			if document != "" {
				insert_db = append(insert_db, []string{"document", document})
			}
			if prefix != "" {
				insert_db = append(insert_db, []string{"prefix", prefix})
			}
			for _, tag := range tag_list {
				insert_db = append(insert_db, []string{"tag", tag})
			}
			for _, value := range insert_db {
				if _, err := tx.Exec(
					tool.DB_change("insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)"),
					value[0],
					set_code,
					set_id,
					value[1],
				); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			panic(err)
		}
		tool.Search_bbs_index_update(db, set_id, set_code)

		return_data["response"] = "ok"
		return_data["data"] = set_code

		return return_data
	}

	old_title := ""
	old_data := ""
	if !tool.QueryRow_DB(
		db,
		"select set_data from bbs_data where set_name = 'title' and set_id = ? and set_code = ?",
		[]any{&old_title},
		set_id,
		set_code,
	) || !tool.QueryRow_DB(
		db,
		"select set_data from bbs_data where set_name = 'data' and set_id = ? and set_code = ?",
		[]any{&old_data},
		set_id,
		set_code,
	) {
		return_data["response"] = "not exist"
		return_data["data"] = "post"

		return return_data
	}

	user_id := ""
	if !tool.QueryRow_DB(
		db,
		"select set_data from bbs_data where set_name = 'user_id' and set_id = ? and set_code = ?",
		[]any{&user_id},
		set_id,
		set_code,
	) {
		return_data["response"] = "not exist"
		return_data["data"] = "post"

		return return_data
	}

	if user_id != config.IP && !tool.Check_permission(db, "bbs_post_manage", config.IP) {
		return_data["response"] = "require auth"

		return return_data
	}

	date := tool.Get_time()
	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		for _, value := range []struct {
			query string
			args  []any
		}{
			{
				"update bbs_data set set_data = ? where set_name = 'title' and set_code = ? and set_id = ?",
				[]any{title, set_code, set_id},
			},
			{
				"update bbs_data set set_data = ? where set_name = 'data' and set_code = ? and set_id = ?",
				[]any{data, set_code, set_id},
			},
			{
				"update bbs_data set set_data = ? where set_name = 'date' and set_code = ? and set_id = ?",
				[]any{date, set_code, set_id},
			},
			{
				"delete from bbs_data where set_name = 'prefix' and set_code = ? and set_id = ?",
				[]any{set_code, set_id},
			},
			{
				"delete from bbs_data where set_name = 'tag' and set_code = ? and set_id = ?",
				[]any{set_code, set_id},
			},
		} {
			if _, err := tx.Exec(tool.DB_change(value.query), value.args...); err != nil {
				return err
			}
		}
		bbs_post_last_activity_update(tx, set_id, set_code, date)
		if prefix != "" {
			if _, err := tx.Exec(
				tool.DB_change("insert into bbs_data (set_name, set_code, set_id, set_data) values ('prefix', ?, ?, ?)"),
				set_code,
				set_id,
				prefix,
			); err != nil {
				return err
			}
		}
		for _, tag := range tag_list {
			if _, err := tx.Exec(
				tool.DB_change("insert into bbs_data (set_name, set_code, set_id, set_data) values ('tag', ?, ?, ?)"),
				set_code,
				set_id,
				tag,
			); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
	tool.Search_bbs_index_update(db, set_id, set_code)

	return_data["response"] = "ok"
	return_data["data"] = set_code

	return return_data
}
