package route

import "opennamu/route/tool"

func Api_record_bbs_in(config tool.Config, user_name string, set_id string, page string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	page_int := tool.Str_to_int(page)
	num := 0
	if page_int*50 > 0 {
		num = page_int*50 - 50
	}

	query := `select set_code from bbs_data where set_name = "user_id" and set_id = ? and set_data = ?`
	if !tool.Check_permission(db, "bbs_post_manage", config.IP) {
		query += ` and not exists (select 1 from bbs_data blind_data where blind_data.set_name = "blind" and blind_data.set_data = "O" and blind_data.set_id = bbs_data.set_id and blind_data.set_code = bbs_data.set_code)`
	}
	query += ` order by set_code desc limit ?, 50`
	rows := tool.Query_DB(
		db,
		query,
		set_id,
		user_name,
		num,
	)
	defer rows.Close()

	data_list := []string{}

	for rows.Next() {
		var set_code string

		err := rows.Scan(&set_code)
		if err != nil {
			panic(err)
		}

		data_list = append(data_list, set_code)
	}

	result_data := make(map[string]any)
	result_data["response"] = "ok"
	result_data["data"] = data_list

	return result_data
}
