package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func Api_history_reset_post(config tool.Config, doc_name string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "history_manage", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if doc_name == "" {
		return_data["response"] = "error"
		return_data["data"] = "invalid document"
		return return_data
	}

	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		if _, err := tx.Exec(tool.DB_change("delete from history where title = ?"), doc_name); err != nil {
			return err
		}
		tool.Do_insert_auth_history(tx, config.IP, "history_reset ("+doc_name+")")
		return nil
	}); err != nil {
		panic(err)
	}

	return_data["response"] = "ok"
	return return_data
}
