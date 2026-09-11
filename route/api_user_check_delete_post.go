package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func Api_user_check_delete_post(config tool.Config, user_name string, user_ip string, today string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "user_manage", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		if _, err := tx.Exec(tool.DB_change("delete from ua_d where name = ? and ip = ? and today = ?"), user_name, user_ip, today); err != nil {
			return err
		}
		tool.Do_insert_auth_history(tx, config.IP, "user_check_delete ("+user_name+")")
		return nil
	}); err != nil {
		panic(err)
	}
	return_data["response"] = "ok"
	return return_data
}
