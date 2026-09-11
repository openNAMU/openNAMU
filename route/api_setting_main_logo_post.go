package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func Api_setting_main_logo_post(config tool.Config, form map[string]string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "setting_main_logo", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		for _, skin := range setting_logo_skins() {
			coverage := ""
			field_name := "main_css"
			if skin != "default" {
				coverage = skin
				field_name = skin
			}
			setting_save_value(tx, "logo", coverage, setting_form_value(form, field_name, ""))
		}
		tool.Do_insert_auth_history(tx, config.IP, "edit_set (logo)")
		return nil
	}); err != nil {
		panic(err)
	}
	return_data["response"] = "ok"
	return return_data
}
