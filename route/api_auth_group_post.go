package route

import (
	"net/url"

	"opennamu/route/tool"
)

func Api_auth_group_post(config tool.Config, name string, values url.Values) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "auth_group_manage", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if name == "" {
		return_data["response"] = "error"
		return return_data
	}

	tool.Exec_DB(db, "delete from alist where name = ?", name)
	for _, choice := range tool.Auth_choices() {
		if values.Get(choice.Key) != "" {
			tool.Exec_DB(db, "insert into alist (name, acl) values (?, ?)", name, choice.Key)
		}
	}
	tool.Exec_DB(db, "insert into alist (name, acl) values (?, 'nothing')", name)
	tool.Do_insert_auth_history(db, config.IP, "auth_group_save ("+name+")")

	return_data["response"] = "ok"
	return return_data
}
