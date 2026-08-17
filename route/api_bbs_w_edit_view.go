package route

import (
	"opennamu/route/tool"
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
		"title":  "",
		"data":   "",
		"prefix": "",
		"tags":   "",
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
	prefix := ""
	tool.QueryRow_DB(
		db,
		"select set_data from bbs_data where set_name = 'prefix' and set_id = ? and set_code = ?",
		[]any{&prefix},
		set_id,
		set_code,
	)
	data["prefix"] = prefix
	rows := tool.Query_DB(
		db,
		"select set_data from bbs_data where set_name = 'tag' and set_id = ? and set_code = ?",
		set_id,
		set_code,
	)
	tag_list := []string{}
	for rows.Next() {
		var tag string
		if rows.Scan(&tag) == nil {
			tag_list = append(tag_list, tag)
		}
	}
	rows.Close()
	data["tags"] = strings.Join(tag_list, ", ")
	return_data["response"] = "ok"
	return_data["data"] = data

	return return_data
}
