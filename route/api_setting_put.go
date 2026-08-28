package route

import (
	"opennamu/route/tool"
)

func Api_setting_put(config tool.Config, set_name string, data string, coverage string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	auth_info := tool.Check_acl(db, "", "", "owner_auth", config.IP)

	setting_acl := Setting_list()
	return_data := make(map[string]any)

	if _, ok := setting_acl[set_name]; ok {
		if set_name == "rankup_condition" {
			if !auth_info {
				return_data["response"] = "require auth"
				return return_data
			}
			if !tool.Auth_group_name_rankup(coverage) {
				return_data["response"] = "error"
				return_data["data"] = "invalid name"
				return return_data
			}

			condition_list, ok := tool.Rankup_condition_list(data)
			if !ok {
				return_data["response"] = "error"
				return_data["data"] = "invalid condition"
				return return_data
			}

			tool.Exec_DB(db, "delete from other where name = 'rankup_condition' and coverage = ?", coverage)
			for _, condition_data := range condition_list {
				tool.Exec_DB(db, "insert into other (name, data, coverage) values ('rankup_condition', ?, ?)", condition_data, coverage)
			}
			tool.Do_insert_auth_history(db, config.IP, "rankup_condition ("+coverage+")")
			return_data["response"] = "ok"
			return return_data
		}

		if tool.Arr_in_str(bbs_global_acl_fields, set_name) && !acl_value_valid(db, data) {
			return_data["response"] = "error"
			return_data["data"] = "invalid acl"
			return return_data
		}
		if auth_info {
			if coverage == "" {
				tool.Exec_DB(
					db,
					"delete from other where name = ?",
					set_name,
				)
			}

			tool.Exec_DB(
				db,
				"insert into other (name, data, coverage) values (?, ?, ?)",
				set_name, data, coverage,
			)

			return_data["response"] = "ok"
		} else {
			return_data["response"] = "require auth"
		}
	} else {
		return_data["response"] = "not exist"
	}

	return return_data
}
