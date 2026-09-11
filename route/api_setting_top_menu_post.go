package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func Api_setting_top_menu_post(config tool.Config, form map[string]string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "setting_top_menu", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		setting_save_value(tx, "top_menu", "", setting_form_value(form, "content", ""))
		tool.Do_insert_auth_history(tx, config.IP, "edit_set (top_menu)")
		return nil
	}); err != nil {
		panic(err)
	}
	return_data["response"] = "ok"
	return return_data
}
