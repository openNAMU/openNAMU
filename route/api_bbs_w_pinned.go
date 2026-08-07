package route

import "opennamu/route/tool"

func Api_bbs_w_pinned(config tool.Config, set_id string, set_code string, toggle bool) map[string]any {
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

	if !tool.Check_acl(db, "", "", "bbs_auth", config.IP) {
		return_data["response"] = "require auth"

		return return_data
	}

	if toggle {
		pinned := ""
		pinned_exist := tool.QueryRow_DB(
			db,
			"select set_data from bbs_data where set_name = 'pinned' and set_id = ? and set_code = ?",
			[]any{&pinned},
			set_id,
			set_code,
		)

		if pinned_exist {
			tool.Exec_DB(
				db,
				"delete from bbs_data where set_name = 'pinned' and set_id = ? and set_code = ?",
				set_id,
				set_code,
			)
		} else {
			tool.Exec_DB(
				db,
				"insert into bbs_data (set_name, set_code, set_id, set_data) values ('pinned', ?, ?, ?)",
				set_code,
				set_id,
				tool.Get_time(),
			)
		}
	}

	pinned := ""
	pinned_exist := tool.QueryRow_DB(
		db,
		"select set_data from bbs_data where set_name = 'pinned' and set_id = ? and set_code = ?",
		[]any{&pinned},
		set_id,
		set_code,
	)

	return_data["response"] = "ok"
	if pinned_exist {
		return_data["data"] = "pinned_release"
	} else {
		return_data["data"] = "pinned"
	}

	return return_data
}
