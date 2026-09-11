package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func Api_user_email_check_post(config tool.Config, key string, email string, input_key string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if tool.IP_or_user(config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if key == "" || email == "" || key != input_key {
		return_data["response"] = "error"
		return_data["data"] = "key error"
		return return_data
	}
	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		user_delete(tx, config.IP, "email")
		user_save(tx, config.IP, "email", email)
		return nil
	}); err != nil {
		panic(err)
	}
	return_data["response"] = "ok"
	return return_data
}
