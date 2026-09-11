package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func Api_user_field_delete_post(config tool.Config, field string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if tool.IP_or_user(config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	remove_email_2fa := field == "email" && user_value(db, config.IP, "2fa") == "email"
	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		user_delete(tx, config.IP, field)
		if remove_email_2fa {
			user_delete(tx, config.IP, "2fa")
		}
		return nil
	}); err != nil {
		panic(err)
	}
	return_data["response"] = "ok"
	return return_data
}
