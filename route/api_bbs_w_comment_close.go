package route

import "opennamu/route/tool"

func Api_bbs_w_comment_close(config tool.Config, set_id string, set_code string, closed bool) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)

	title := ""
	if !tool.QueryRow_DB(
		db,
		"select set_data from bbs_data where set_name = 'title' and set_id = ? and set_code = ?",
		[]any{&title},
		set_id,
		set_code,
	) {
		return_data["response"] = "not exist"
		return_data["data"] = "post"
		return return_data
	}

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	tool.Exec_DB(
		db,
		"delete from bbs_data where set_name = 'comment_close' and set_id = ? and set_code = ?",
		set_id,
		set_code,
	)
	if closed {
		tool.Exec_DB(
			db,
			"insert into bbs_data (set_name, set_code, set_id, set_data) values ('comment_close', ?, ?, '1')",
			set_code,
			set_id,
		)
	}

	return_data["response"] = "ok"
	return return_data
}
