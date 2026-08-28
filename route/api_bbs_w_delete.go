package route

import "opennamu/route/tool"

func Api_bbs_w_delete(config tool.Config, set_id string, set_code string) map[string]any {
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

	if !tool.Check_permission(db, "bbs_delete", config.IP) {
		return_data["response"] = "require auth"

		return return_data
	}

	tool.Exec_DB(
		db,
		"delete from bbs_data where set_id = ? and set_code = ?",
		set_id,
		set_code,
	)
	tool.Exec_DB(
		db,
		"delete from bbs_set where set_id = ? and set_code = ?",
		set_id,
		set_code,
	)
	tool.Exec_DB(
		db,
		"delete from bbs_data where set_id = ? or set_id like ?",
		set_id+"-"+set_code,
		set_id+"-"+set_code+"-%",
	)

	return_data["response"] = "ok"

	return return_data
}
