package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func Api_auth_group_delete_post(config tool.Config, name string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if tool.Auth_group_name_default(name) {
		return_data["response"] = "default"
		return return_data
	}
	if !tool.Check_permission(db, "auth_group_manage", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if tool.Auth_group_in_use(db, name) {
		return_data["response"] = "error"
		return return_data
	}
	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		if _, err := tx.Exec(tool.DB_change("delete from alist where name = ?"), name); err != nil {
			return err
		}
		tool.Do_insert_auth_history(tx, config.IP, "auth_group_delete ("+name+")")
		return nil
	}); err != nil {
		panic(err)
	}

	return_data["response"] = "ok"
	return return_data
}
