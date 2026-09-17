package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func Api_user_setting_field_post(config tool.Config, field string, value string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if tool.IP_or_user(config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if field == "user_name" && !tool.Get_user_name_check(db, value) && value != User_value(db, config.IP, "user_name") {
		return_data["response"] = "error"
		return_data["data"] = "user name error"
		return return_data
	}
	remove_email_2fa := field == "email" && value == "" && User_value(db, config.IP, "2fa") == "email"
	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		if field == "email" && value == "" {
			User_delete(tx, config.IP, field)
			if remove_email_2fa {
				User_delete(tx, config.IP, "2fa")
			}
		} else {
			User_save(tx, config.IP, field, value)
		}
		return nil
	}); err != nil {
		panic(err)
	}
	return_data["response"] = "ok"
	return return_data
}
