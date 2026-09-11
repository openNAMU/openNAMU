package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func Api_setting_sitemap_set_post(config tool.Config, form map[string]string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "setting_sitemap", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		setting_save_fields(tx, setting_sitemap_fields(), form)
		tool.Do_insert_auth_history(tx, config.IP, "edit_set (sitemap)")
		return nil
	}); err != nil {
		panic(err)
	}
	sync_indexnow_key(db)
	return_data["response"] = "ok"
	return return_data
}
