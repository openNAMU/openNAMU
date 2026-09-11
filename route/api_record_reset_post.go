package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func Api_record_reset_post(config tool.Config, user_name string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "record_manage", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if user_name == "" {
		return_data["response"] = "error"
		return return_data
	}
	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		if _, err := tx.Exec(tool.DB_change("delete from history where ip = ?"), user_name); err != nil {
			return err
		}
		tool.Do_insert_auth_history(tx, config.IP, "record_reset ("+user_name+")")
		return nil
	}); err != nil {
		panic(err)
	}
	return_data["response"] = "ok"
	return return_data
}
