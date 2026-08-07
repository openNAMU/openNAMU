package route

import "opennamu/route/tool"

func Api_bbs_delete(config tool.Config, set_id string) map[string]any {
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

	if set_id == "0" {
		return_data["response"] = "error"
		return_data["data"] = "not allowed"

		return return_data
	}

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return_data["response"] = "require auth"

		return return_data
	}

	tool.Exec_DB(db, "delete from bbs_data where set_id = ?", set_id)
	tool.Exec_DB(db, "delete from bbs_set where set_id = ?", set_id)
	tool.Exec_DB(db, "delete from bbs_data where set_id like ?", set_id+"-%")

	return_data["response"] = "ok"

	return return_data
}
