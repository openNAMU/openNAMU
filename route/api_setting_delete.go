package route

import (
	"opennamu/route/tool"
)

func Api_setting_delete(config tool.Config, set_name string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	auth_info := tool.Check_acl(db, "", "", "owner_auth", config.IP)

	setting_acl := Setting_list()
	return_data := make(map[string]any)

	if _, ok := setting_acl[set_name]; ok {
		if auth_info {
			tool.Exec_DB(
				db,
				"delete from other where name = ?",
				set_name,
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
