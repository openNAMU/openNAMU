package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func api_bbs_tabom_post(db *sql.Tx, user string, set_id string, set_code string, vote_type string) map[string]any {
	return_data := make(map[string]any)

	if vote_type != "down" {
		vote_type = "up"
	}

	selected_list := "tabom_list"
	selected_count := "tabom_count"
	if vote_type == "down" {
		selected_list = "tabom_down_list"
		selected_count = "tabom_down_count"
	}

	if bbs_tabom_user_exists(db, selected_list, user, set_id, set_code) {
		return_data["response"] = "same user exist"
		return return_data
	}

	change_bbs_tabom_count(db, selected_count, set_id, set_code, 1)
	tool.Exec_DB(
		db,
		"insert into bbs_data (set_name, set_data, set_id, set_code) values (?, ?, ?, ?)",
		selected_list,
		user,
		set_id,
		set_code,
	)
	return_data["response"] = "ok"

	return return_data
}
