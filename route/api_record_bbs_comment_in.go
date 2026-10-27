package route

import "opennamu/route/tool"

func Api_record_bbs_comment_in(config tool.Config, user_name string, bbs_id string, page string) map[string]any {
	if !Record_bbs_allowed(config, user_name) {
		return map[string]any{"response": "require auth", "data": [][]string{}}
	}

	db := tool.DB_connect()
	defer tool.DB_close(db)

	page_int := tool.Str_to_int(page)
	skip := 0
	if page_int*50 > 0 {
		skip = page_int*50 - 50
	}

	visible_page_size := 50
	fetch_batch := 500

	var auth_info map[string]bool
	seen := 0
	data_list := [][]string{}
	offset := 0

	for len(data_list) < visible_page_size {
		rows := tool.Query_DB(
			db,
			`select c.set_id, c.set_code, coalesce(d.set_data, '') as date
			 from bbs_data c
			 left join bbs_data d on d.set_name = "comment_date" and d.set_id = c.set_id and d.set_code = c.set_code
			 where c.set_name = "comment_user_id" and c.set_id like ? and c.set_data = ?
			 order by coalesce(d.set_data, '') desc, c.set_id desc limit ?, ?`,
			bbs_id+"-%",
			user_name,
			offset,
			fetch_batch,
		)
		defer rows.Close()

		fetched := false
		for rows.Next() {
			var set_id string
			var set_code string
			var date string

			if err := rows.Scan(&set_id, &set_code, &date); err != nil {
				panic(err)
			}
			fetched = true

			post_set_id, post_set_code, _, exists := Bbs_comment_storage_location(set_id, set_code)
			if exists {
				if auth_info == nil {
					auth_info = tool.Get_auth_info(db, config.IP)
				}
				if !Bbs_post_blind_allowed(db, post_set_id, post_set_code, config.IP, auth_info) {
					continue
				}
			}

			seen++
			if seen <= skip {
				continue
			}
			if seen > skip+visible_page_size {
				break
			}

			data_list = append(data_list, []string{
				bbs_id,
				post_set_code,
				set_id,
				set_code,
				date,
			})
		}

		offset += fetch_batch
		if !fetched {
			break
		}
	}

	result_data := map[string]any{}
	result_data["response"] = "ok"
	result_data["data"] = data_list

	return result_data
}
