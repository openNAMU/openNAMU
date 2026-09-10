package route

import "opennamu/route/tool"

func Api_bbs_w_blind_post(config tool.Config, set_id string, set_code string, blind bool) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "bbs_post_manage", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	_, exists := tool.Get_bbs_data_value(db, set_id, set_code, "title")
	if !exists {
		return_data["response"] = "not exist"
		return_data["data"] = "post"
		return return_data
	}

	if blind {
		_, blind_exists := tool.Get_bbs_data_value(db, set_id, set_code, "blind")
		if blind_exists {
			tool.Exec_DB(
				db,
				"update bbs_data set set_data = 'O' where set_name = 'blind' and set_id = ? and set_code = ?",
				set_id,
				set_code,
			)
		} else {
			tool.Exec_DB(
				db,
				"insert into bbs_data (set_name, set_id, set_code, set_data) values ('blind', ?, ?, 'O')",
				set_id,
				set_code,
			)
		}
	} else {
		tool.Exec_DB(
			db,
			"delete from bbs_data where set_name = 'blind' and set_id = ? and set_code = ?",
			set_id,
			set_code,
		)
	}

	return_data["response"] = "ok"
	if blind {
		return_data["data"] = "O"
	} else {
		return_data["data"] = ""
	}
	return return_data
}
