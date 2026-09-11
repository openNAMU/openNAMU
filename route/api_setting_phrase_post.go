package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func Api_setting_phrase_post(config tool.Config, form map[string]string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "setting_phrase", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		for _, field := range setting_phrase_fields() {
			setting_save_value(tx, field.name, "", setting_form_value(form, field.name, ""))
		}
		tool.Do_insert_auth_history(tx, config.IP, "edit_set (phrase)")
		return nil
	}); err != nil {
		panic(err)
	}
	return_data["response"] = "ok"
	return return_data
}
