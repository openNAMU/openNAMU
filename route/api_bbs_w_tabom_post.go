package route

import (
	"database/sql"
	"strconv"
	"strings"

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

func Api_bbs_w_tabom_post(config tool.Config, sub_code string, vote_type string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	sub_code_parts := strings.Split(sub_code, "-")

	bbs_num := ""
	post_num := ""

	if len(sub_code_parts) > 1 {
		bbs_num = sub_code_parts[0]
		post_num = sub_code_parts[1]
	}

	return_data := make(map[string]any)

	if !tool.Check_permission(db, "bbs_comment", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	if vote_type != "down" {
		vote_type = "up"
	}

	selected_list := "tabom_list"
	selected_count := "tabom_count"
	if vote_type == "down" {
		selected_list = "tabom_down_list"
		selected_count = "tabom_down_count"
	}

	if bbs_tabom_user_exists(db, selected_list, config.IP, bbs_num, post_num) {
		return_data["response"] = "same user exist"
		return return_data
	}

	change_bbs_tabom_count(db, selected_count, bbs_num, post_num, 1)
	tool.Exec_DB(
		db,
		"insert into bbs_data (set_name, set_data, set_id, set_code) values (?, ?, ?, ?)",
		selected_list,
		config.IP,
		bbs_num,
		post_num,
	)
	return_data["response"] = "ok"

	return return_data
}
