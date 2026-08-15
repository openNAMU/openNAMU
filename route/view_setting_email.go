package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func View_setting_email_test(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	data := `<form method="post">`
	data += setting_input("title", "", "text") + setting_hr()
	data += setting_input("email", "", "email") + setting_hr()
	data += setting_textarea("data", "", "opennamu_textarea_500") + setting_hr()
	data += `<button type="submit">` + tool.Get_language(db, "send", true) + `</button></form>`

	return setting_page(db, config, tool.Get_language(db, "email_test", true), data, "setting/external")
}

func View_setting_email_test_post(config tool.Config, form map[string]string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	api_data := Api_func_email_post(
		config,
		setting_form_value(form, "email", ""),
		setting_form_value(form, "title", ""),
		setting_form_value(form, "data", ""),
	)
	response, _ := api_data["response"].(string)
	message := tool.Get_language(db, "error", true)
	if response == "ok" {
		message = tool.Get_language(db, "ok", true)
	}

	return view_setting_email_test_result(db, config, message)
}

func view_setting_email_test_result(db *sql.DB, config tool.Config, message string) string {
	return setting_page(db, config, tool.Get_language(db, "email_test", true), message, "setting/external")
}
