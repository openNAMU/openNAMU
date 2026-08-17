package route

import "opennamu/route/tool"

func Api_bbs_w_comment(config tool.Config, do_type string, sub_code string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if do_type == "" || do_type == "normal" {
		do_type = "around"
	}

	if do_type == "length" {
		bbs_and_post_num := sub_code

		comment_length := "0"
		tool.QueryRow_DB(
			db,
			"select count(*) from bbs_data where set_name = 'comment_date' and set_id = ? order by set_code + 0 desc",
			[]any{&comment_length},
			bbs_and_post_num,
		)

		reply_length := "0"
		tool.QueryRow_DB(
			db,
			"select count(*) from bbs_data where set_name = 'comment_date' and set_id like ? order by set_code + 0 desc",
			[]any{&reply_length},
			bbs_and_post_num+"-%",
		)

		comment_length_int := tool.Str_to_int(comment_length)
		reply_length_int := tool.Str_to_int(reply_length)

		length_int := comment_length_int + reply_length_int

		data_list := map[string]any{
			"response": "ok",
			"comment":  comment_length,
			"reply":    reply_length,
			"data":     length_int,
		}

		return data_list
	} else {
		return_data := make(map[string]any)

		temp := []map[string]string{}
		if !tool.Check_acl(db, "", "", "bbs_comment", config.IP) {
			return_data["response"] = "require auth"
			return_data["data"] = temp

			return return_data
		}

		temp_data := Api_bbs_w_comment_all(config, sub_code, true, do_type)
		temp, _ = temp_data["data"].([]map[string]string)

		return_data["response"] = "ok"
		return_data["data"] = temp

		return return_data
	}
}
