package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func Api_user_head_reset_post(config tool.Config, skin_name string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if tool.IP_or_user(config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		user_save(tx, config.IP, "custom_css", "")
		user_save(tx, config.IP, "custom_css_"+skin_name, "")
		user_delete(tx, config.IP, "head")
		return nil
	}); err != nil {
		panic(err)
	}
	return_data["response"] = "ok"
	return return_data
}
