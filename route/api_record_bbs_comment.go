package route

import "opennamu/route/tool"

func Api_record_bbs_comment(config tool.Config, user_name string, page string) map[string]any {
	if !Record_bbs_allowed(config, user_name) {
		return map[string]any{"response": "require auth", "data": [][]string{}}
	}

	db := tool.DB_connect()
	defer tool.DB_close(db)

	page_int := tool.Str_to_int(page)
	num := 0
	if page_int*50 > 0 {
		num = page_int*50 - 50
	}

	rows := tool.Query_DB(
		db,
		`select case when u.set_id like '-1-%' then '-1' else substr(u.set_id, 1, instr(u.set_id, '-') - 1) end as bbs_id, coalesce(max(d.set_data), '') as date
		 from bbs_data u
		 left join bbs_data d on d.set_name = 'comment_date' and d.set_id = u.set_id and d.set_code = u.set_code
		 where u.set_name = 'comment_user_id' and u.set_data = ?
		 group by bbs_id
		 order by bbs_id + 0 desc
		 limit ?, 50`,
		user_name,
		num,
	)
	defer rows.Close()

	data_list := [][]string{}

	for rows.Next() {
		var bbs_id string
		var date string

		err := rows.Scan(&bbs_id, &date)
		if err != nil {
			panic(err)
		}

		data_list = append(data_list, []string{bbs_id, date})
	}

	result_data := make(map[string]any)
	result_data["response"] = "ok"
	result_data["data"] = data_list

	return result_data
}
