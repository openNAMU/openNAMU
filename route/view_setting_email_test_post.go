package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func View_setting_email_test_post(config tool.Config, form map[string]string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	api_data := Api_func_email_post(
		config,
		setting_form_value(form, "email", ""),
		setting_form_value(form, "title", ""),
		setting_form_value(form, "data", ""),
	)
	response, _ := api_data["response"].(string)
	if response == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	}
	message := tool.Get_language(db, "error", true)
	if response == "ok" {
		message = tool.Get_language(db, "ok", true)
	}

	return view_setting_email_test_result(db, config, message)
}

func view_setting_email_test_result(db *sql.DB, config tool.Config, message string) string {
	return setting_page(db, config, tool.Get_language(db, "email_test", true), message, "setting/external")
}
