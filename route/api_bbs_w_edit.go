package route

import (
	"opennamu/route/tool"
	"strconv"
	"strings"
)

func Api_bbs_w_edit_view(config tool.Config, set_id string, set_code string, comment_code string) map[string]any {
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

	if !tool.Check_acl(db, set_id, "", "bbs_edit", config.IP) {
		return_data["response"] = "require auth"

		return return_data
	}

	data := map[string]string{
		"title": "",
		"data":  "",
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

		if comment_user_id != config.IP && !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
			return_data["response"] = "require auth"

			return return_data
		}

		data["data"] = comment
		return_data["response"] = "ok"
		return_data["data"] = data

		return return_data
	}

	if set_code == "" {
		return_data["response"] = "ok"
		return_data["data"] = data

		return return_data
	}

	title := ""
	content := ""
	if !tool.QueryRow_DB(
		db,
		"select set_data from bbs_data where set_name = 'title' and set_id = ? and set_code = ?",
		[]any{&title},
		set_id,
		set_code,
	) || !tool.QueryRow_DB(
		db,
		"select set_data from bbs_data where set_name = 'data' and set_id = ? and set_code = ?",
		[]any{&content},
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

	if user_id != config.IP && !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return_data["response"] = "require auth"

		return return_data
	}

	data["title"] = title
	data["data"] = content
	return_data["response"] = "ok"
	return_data["data"] = data

	return return_data
}

func Api_bbs_w_edit_post(config tool.Config, set_id string, set_code string, comment_code string, title string, data string) map[string]any {
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

	if !tool.Check_acl(db, set_id, "", "bbs_edit", config.IP) {
		return_data["response"] = "require auth"

		return return_data
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

		if comment_user_id != config.IP && !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
			return_data["response"] = "require auth"

			return return_data
		}

		tool.Exec_DB(
			db,
			"update bbs_data set set_data = ? where set_name = 'comment' and set_code = ? and set_id = ?",
			data,
			comment_set_code,
			comment_set_id,
		)

		return_data["response"] = "ok"
		return_data["data"] = set_code

		return return_data
	}

	if set_code == "" {
		last_code := ""
		tool.QueryRow_DB(
			db,
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
			{"user_id", config.IP},
		}
		for _, v := range insert_db {
			tool.Exec_DB(
				db,
				"insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)",
				v[0],
				set_code,
				set_id,
				v[1],
			)
		}

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

	if user_id != config.IP && !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return_data["response"] = "require auth"

		return return_data
	}

	date := tool.Get_time()
	tool.Exec_DB(db, "update bbs_data set set_data = ? where set_name = 'title' and set_code = ? and set_id = ?", title, set_code, set_id)
	tool.Exec_DB(db, "update bbs_data set set_data = ? where set_name = 'data' and set_code = ? and set_id = ?", data, set_code, set_id)
	tool.Exec_DB(db, "update bbs_data set set_data = ? where set_name = 'date' and set_code = ? and set_id = ?", date, set_code, set_id)

	return_data["response"] = "ok"
	return_data["data"] = set_code

	return return_data
}
