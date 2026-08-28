package route

import (
	"opennamu/route/tool"
)

func BBS_w_set_list() map[string]string {
	setting_acl := map[string]string{}

	setting_acl["bbs_view_acl"] = ""
	setting_acl["bbs_only_my_data_view_acl"] = ""
	setting_acl["bbs_acl"] = ""
	setting_acl["bbs_edit_acl"] = ""
	setting_acl["bbs_comment_acl"] = ""

	setting_acl["bbs_markup"] = ""
	setting_acl["bbs_name"] = ""
	setting_acl["bbs_prefix"] = ""

	return setting_acl
}

func Api_bbs_w_set(config tool.Config, set_id string, set_name string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	setting_acl := BBS_w_set_list()

	if val, ok := setting_acl[set_name]; ok {
		if val != "" {
			if !tool.Check_permission(db, "bbs_setting", config.IP) {
				return_data := make(map[string]any)
				return_data["response"] = "require auth"

				return return_data
			}
		}

		rows := tool.Query_DB(
			db,
			"select set_data, set_code from bbs_set where set_name = ? and set_id = ?",
			set_name, set_id,
		)
		defer rows.Close()

		data_list := [][]string{}

		for rows.Next() {
			var set_data string
			var set_coverage string

			err := rows.Scan(&set_data, &set_coverage)
			if err != nil {
				panic(err)
			}

			data_list = append(data_list, []string{set_data, set_coverage})
		}

		return_data := make(map[string]any)
		return_data["response"] = "ok"
		return_data["data"] = data_list

		return return_data
	} else {
		return_data := make(map[string]any)
		return_data["response"] = "not exist"

		return return_data
	}
}
