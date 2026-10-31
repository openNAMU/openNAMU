package route

import (
	"database/sql"
	"fmt"

	"opennamu/route/tool"
)

func Api_setting_main_post(config tool.Config, form map[string]string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "setting_main", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if !Acl_value_valid(db, form["user_document_view_acl_all"]) {
		return_data["response"] = "error"
		return_data["data"] = "error"
		return return_data
	}
	if enabled, exists := form[related_document_enabled_setting]; exists && enabled != "1" && enabled != "0" {
		return_data["response"] = "error"
		return_data["data"] = "error"
		return return_data
	}
	// Validate trusted_proxies before it reaches the database, so an unsafe or
	// malformed value is rejected and the previous valid setting is preserved.
	if trusted_raw, exists := form["trusted_proxies"]; exists {
		if err := tool.Validate_trusted_proxies(trusted_raw); err != nil {
			return_data["response"] = "error"
			return_data["data"] = fmt.Sprintf("trusted_proxies: %v", err)
			return return_data
		}
	}
	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		Setting_save_fields(tx, Setting_main_fields(), form)
		tool.Do_insert_auth_history(tx, config.IP, "edit_set (main)")
		return nil
	}); err != nil {
		panic(err)
	}
	return_data["response"] = "ok"
	return return_data
}
