package route

import (
	"database/sql"
	"strconv"

	"opennamu/route/tool"
)

func change_bbs_tabom_count(db *sql.DB, set_name string, set_id string, set_code string, amount int) {
	count := "0"
	exists := tool.QueryRow_DB(
		db,
		"select set_data from bbs_data where set_name = ? and set_id = ? and set_code = ?",
		[]any{&count},
		set_name,
		set_id,
		set_code,
	)
	if !exists {
		tool.Exec_DB(
			db,
			"insert into bbs_data (set_name, set_data, set_id, set_code) values (?, ?, ?, ?)",
			set_name,
			"0",
			set_id,
			set_code,
		)
	}

	count_int := tool.Str_to_int(count) + amount
	if count_int < 0 {
		count_int = 0
	}
	tool.Exec_DB(
		db,
		"update bbs_data set set_data = ? where set_name = ? and set_id = ? and set_code = ?",
		strconv.Itoa(count_int),
		set_name,
		set_id,
		set_code,
	)
}

func bbs_tabom_user_exists(db *sql.DB, set_name string, user string, set_id string, set_code string) bool {
	data := ""
	return tool.QueryRow_DB(
		db,
		"select set_data from bbs_data where set_name = ? and set_data = ? and set_id = ? and set_code = ?",
		[]any{&data},
		set_name,
		user,
		set_id,
		set_code,
	)
}

func Api_bbs_w_tabom_post(config tool.Config, set_id string, set_code string, vote_type string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)

	if !tool.Check_permission(db, "bbs_comment", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	if !bbs_post_blind_allowed(db, set_id, set_code, config.IP, nil) {
		return_data["response"] = "require auth"
		return return_data
	}

	return api_bbs_tabom_post(db, config.IP, set_id, set_code, vote_type)
}
