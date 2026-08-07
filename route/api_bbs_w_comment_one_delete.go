package route

import (
	"opennamu/route/tool"
	"strings"
)

func Api_bbs_w_comment_one_delete(config tool.Config, set_id string, set_code string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)

	set_code_split := strings.Split(set_code, "-")
	if len(set_code_split) < 2 {
		return_data["response"] = "not exist"
		return_data["data"] = "comment"

		return return_data
	}

	post_code := set_code_split[0]
	comment_set_id := set_id + "-" + post_code
	comment_set_code := set_code_split[len(set_code_split)-1]
	if len(set_code_split) > 2 {
		comment_set_id += "-" + strings.Join(set_code_split[1:len(set_code_split)-1], "-")
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

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return_data["response"] = "require auth"

		return return_data
	}

	tool.Exec_DB(
		db,
		"update bbs_data set set_data = '' where set_id = ? and set_code = ?",
		comment_set_id,
		comment_set_code,
	)

	return_data["response"] = "ok"

	return return_data
}
