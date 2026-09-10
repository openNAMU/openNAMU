package route

import "opennamu/route/tool"

type bbs_comment_page_state struct {
	offset int
	limit  int
	data   []map[string]string
}

func api_bbs_w_comment_page_data(config tool.Config, sub_code string, already_auth_check bool, do_type string, state *bbs_comment_page_state) {
	if len(state.data) >= state.limit {
		return
	}

	return_data := Api_bbs_w_comment_one(config, already_auth_check, do_type, sub_code)
	return_data_in, _ := return_data["data"].([]map[string]string)
	for _, comment := range return_data_in {
		if state.offset > 0 {
			state.offset -= 1
		} else {
			state.data = append(state.data, comment)
		}

		if len(state.data) >= state.limit {
			return
		}
		api_bbs_w_comment_page_data(config, sub_code+"-"+comment["code"], already_auth_check, do_type, state)
		if len(state.data) >= state.limit {
			return
		}
	}
}

func Api_bbs_w_comment_page(config tool.Config, sub_code string, already_auth_check bool, do_type string, page int) map[string]any {
	if page < 1 {
		page = 1
	}

	state := &bbs_comment_page_state{
		offset: (page - 1) * 50,
		limit:  50,
		data:   []map[string]string{},
	}
	api_bbs_w_comment_page_data(config, sub_code, already_auth_check, do_type, state)

	return map[string]any{
		"response": "ok",
		"data":     state.data,
	}
}

func api_bbs_w_comment_pinned_data(config tool.Config, set_id string, set_code string) []map[string]string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	root := set_id + "-" + set_code
	rows := tool.Query_DB(
		db,
		"select set_id, set_code from bbs_data where set_name = 'pinned' and (set_id = ? or set_id like ?) order by set_id, set_code",
		root,
		root+"-%",
	)
	defer rows.Close()

	data := []map[string]string{}
	for rows.Next() {
		comment_set_id := ""
		comment_set_code := ""
		if rows.Scan(&comment_set_id, &comment_set_code) != nil {
			continue
		}
		post_set_id, post_set_code, _, valid := bbs_comment_storage_location(comment_set_id, comment_set_code)
		if !valid || post_set_id != set_id || post_set_code != set_code {
			continue
		}

		comment_data := Api_bbs_w_comment_one(config, true, "", comment_set_id+"-"+comment_set_code)
		comment_list, _ := comment_data["data"].([]map[string]string)
		if len(comment_list) > 0 {
			data = append(data, comment_list[0])
		}
	}

	return data
}

func Api_bbs_w_comment_all(config tool.Config, sub_code string, already_auth_check bool, do_type string) map[string]any {
	end_data := []map[string]string{}

	return_data := Api_bbs_w_comment_one(config, already_auth_check, do_type, sub_code)
	return_data_in := return_data["data"].([]map[string]string)

	for for_a := 0; for_a < len(return_data_in); for_a++ {
		end_data = append(end_data, return_data_in[for_a])

		temp_data := Api_bbs_w_comment_all(config, sub_code+"-"+return_data_in[for_a]["code"], already_auth_check, do_type)
		temp, _ := temp_data["data"].([]map[string]string)
		if len(temp) > 0 {
			for for_b := 0; for_b < len(temp); for_b++ {
				end_data = append(end_data, temp[for_b])
			}
		}
	}

	return map[string]any{
		"response": "ok",
		"data":     end_data,
	}
}
