package route

import (
	"database/sql"
	"strings"

	"opennamu/route/tool"
)

func bbs_post_location(db *sql.DB, sub_code string) (string, string, bool) {
	best_id := ""
	best_code := ""
	for _, set_id := range bbs_list(db) {
		prefix := set_id + "-"
		if !strings.HasPrefix(sub_code, prefix) || len(set_id) <= len(best_id) {
			continue
		}

		post_code := strings.TrimPrefix(sub_code, prefix)
		post_code_parts := strings.SplitN(post_code, "-", 2)
		if post_code_parts[0] == "" {
			continue
		}
		best_id = set_id
		best_code = post_code_parts[0]
	}

	return best_id, best_code, best_id != ""
}

func Api_bbs_w_comment_one(config tool.Config, already_auth_check bool, do_type string, sub_code string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !already_auth_check {
		set_id, set_code, exists := bbs_post_location(db, sub_code)
		if !exists {
			return map[string]any{"response": "not exist", "data": []map[string]string{}}
		}
		if _, allowed := bbs_post_view_auth(db, set_id, set_code, config.IP); !allowed {
			return map[string]any{"response": "require auth", "data": []map[string]string{}}
		}
	}

	sub_code_parts := strings.Split(sub_code, "-")
	sub_code_last := ""
	new_sub_code := ""

	if do_type == "around" {
		new_sub_code = sub_code
	} else {
		if len(sub_code_parts) > 2 {
			sub_code_last = sub_code_parts[len(sub_code_parts)-1]
			sub_code_parts = sub_code_parts[:len(sub_code_parts)-1]

			new_sub_code = strings.Join(sub_code_parts, "-")
		}
	}

	var rows *sql.Rows
	if do_type == "around" {
		rows = tool.Query_DB(
			db,
			"select set_name, set_data, set_code, set_id from bbs_data where (set_name = 'comment' or set_name like 'comment%' or set_name in ('blind', 'pinned', 'tabom_count', 'tabom_down_count')) and set_id = ? order by set_code + 0 asc, set_name asc",
			new_sub_code,
		)
	} else {
		rows = tool.Query_DB(
			db,
			"select set_name, set_data, set_code, set_id from bbs_data where (set_name = 'comment' or set_name like 'comment%' or set_name in ('blind', 'pinned', 'tabom_count', 'tabom_down_count')) and set_id = ? and set_code = ? order by set_name asc",
			new_sub_code, sub_code_last,
		)
	}
	defer rows.Close()

	data_list := []map[string]string{}
	temp_dict := map[string]string{}
	ip_parser_temp := map[string][]string{}
	before_set_code := ""

	for rows.Next() {
		var set_name string
		var set_data string
		var set_code string
		var set_id string

		err := rows.Scan(&set_name, &set_data, &set_code, &set_id)
		if err != nil {
			panic(err)
		}

		if before_set_code != set_code {
			if before_set_code != "" {
				data_list = append(data_list, temp_dict)
			}

			temp_dict = map[string]string{}
			temp_dict["id"] = set_id
			temp_dict["code"] = set_code

			before_set_code = set_code
		}

		if set_name == "comment_user_id" {
			var ip_pre string
			var ip_render string

			if _, ok := ip_parser_temp[set_data]; ok {
				ip_pre = ip_parser_temp[set_data][0]
				ip_render = ip_parser_temp[set_data][1]
			} else {
				ip_pre = tool.IP_preprocess(db, set_data, config.IP)[0]
				ip_render = tool.Get_user_profile_image_ui(db, set_data) + tool.IP_parser(db, set_data, config.IP)

				ip_parser_temp[set_data] = []string{ip_pre, ip_render}
			}

			temp_dict["comment_user_id"] = ip_pre
			temp_dict["comment_user_id_render"] = ip_render
		} else {
			temp_dict[set_name] = set_data
		}
	}

	if before_set_code != "" {
		data_list = append(data_list, temp_dict)
	}

	if !tool.Check_permission(db, "bbs_comment_manage", config.IP) {
		for _, comment_data := range data_list {
			if comment_data["blind"] == "O" {
				comment_data["comment"] = ""
			}
		}
	}

	return_data := make(map[string]any)
	return_data["response"] = "ok"

	if do_type == "around" {
		return_data["data"] = data_list
	} else {
		if len(data_list) > 0 {
			return_data["data"] = []map[string]string{
				data_list[0],
			}
		} else {
			return_data["data"] = []map[string]string{}
		}
	}

	return return_data
}
