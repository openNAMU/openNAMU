package route

import "opennamu/route/tool"

func Api_record_bbs(config tool.Config, user_name string, page string) map[string]any {
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
		`select u.set_id, coalesce(max(d.set_data), '') as date
		 from bbs_data u
		 left join bbs_data d on d.set_name = 'date' and d.set_id = u.set_id and d.set_code = u.set_code
		 where u.set_name = 'user_id' and u.set_data = ? and (u.set_id + 0 > 0 or u.set_id in ('0', '-1'))
		 group by u.set_id
		 order by u.set_id + 0 desc
		 limit ?, 50`,
		user_name,
		num,
	)
	defer rows.Close()

	data_list := [][]string{}

	for rows.Next() {
		var set_id string
		var date string

		err := rows.Scan(&set_id, &date)
		if err != nil {
			panic(err)
		}

		data_list = append(data_list, []string{set_id, date})
	}

	result_data := make(map[string]any)
	result_data["response"] = "ok"
	result_data["data"] = data_list

	return result_data
}
