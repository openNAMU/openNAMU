package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func Api_setting_external_post(config tool.Config, form map[string]string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "setting_external", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if form["recaptcha_ver"] == "" {
		form["recaptcha_ver"] = "altcha_high"
	}
	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		setting_save_fields(tx, setting_external_fields(), form)
		tool.Do_insert_auth_history(tx, config.IP, "edit_set (external)")
		return nil
	}); err != nil {
		panic(err)
	}
	return_data["response"] = "ok"
	return return_data
}
