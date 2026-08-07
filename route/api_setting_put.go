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
