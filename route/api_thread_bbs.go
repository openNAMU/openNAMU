package route

import (
	"database/sql"

	"opennamu/route/tool"
)

const thread_bbs_id = "-1"

// Api_thread_bbs reads discussion comments from the thread BBS.
func Api_thread_bbs(config tool.Config, tool_name string, topic_num string, s_num string, e_num string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, thread_bbs_id, "", "bbs_view", config.IP) {
		return map[string]any{
			"response": "require auth",
			"data":     []map[string]string{},
		}
	}

	title := ""
	if !tool.QueryRow_DB(
		db,
		"select set_data from bbs_data where set_name = 'title' and set_id = ? and set_code = ?",
		[]any{&title},
		thread_bbs_id,
		topic_num,
	) {
		return map[string]any{
			"response": "not exist",
			"data":     "thread",
		}
	}

	if tool_name == "length" {
		length := "0"
		tool.QueryRow_DB(
			db,
			"select count(*) from bbs_data where set_name = 'comment' and set_id = ?",
			[]any{&length},
			thread_bbs_id+"-"+topic_num,
		)

		return map[string]any{
			"response": "ok",
			"comment":  length,
			"reply":    "0",
			"data":     tool.Str_to_int(length),
		}
	}

	comment_api := Api_bbs_w_comment_all(config, thread_bbs_id+"-"+topic_num, true, "around")
	comments, _ := comment_api["data"].([]map[string]string)
	admin_auth := tool.Check_permission(db, "bbs_comment_manage", config.IP)
	data_list := []map[string]string{}

	for _, comment := range comments {
		comment_code := comment["code"]
		if tool_name == "top" && comment["top"] != "O" {
			continue
		}
		if s_num != "" && e_num != "" {
			code := tool.Str_to_int(comment_code)
			if code < tool.Str_to_int(s_num) || code > tool.Str_to_int(e_num) {
				continue
			}
		}

		data := comment["comment"]
		if comment["blind"] == "O" && !admin_auth {
			data = ""
		}

		data_list = append(data_list, map[string]string{
			"id":                     topic_num,
			"code":                   comment_code,
			"comment":                data,
			"comment_date":           comment["comment_date"],
			"comment_user_id":        comment["comment_user_id"],
			"comment_user_id_render": comment["comment_user_id_render"],
			"blind":                  comment["blind"],
			"top":                    comment["top"],
		})
	}

	return map[string]any{
		"response": "ok",
		"data":     data_list,
	}
}

func thread_bbs_root_exists(db *sql.DB, topic_num string) bool {
	title := ""
	return tool.QueryRow_DB(
		db,
		"select set_data from bbs_data where set_name = 'title' and set_id = ? and set_code = ?",
		[]any{&title},
		thread_bbs_id,
		topic_num,
	)
}
