package route

import "opennamu/route/tool"

func View_setting_email_test(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "setting_email_test", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	data := `<form method="post">`
	data += setting_input("title", "", "text") + setting_hr()
	data += setting_input("email", "", "email") + setting_hr()
	data += setting_textarea("data", "", "opennamu_textarea_500") + setting_hr()
	data += `<button type="submit">` + tool.Get_language(db, "send", true) + `</button></form>`

	return setting_page(db, config, tool.Get_language(db, "email_test", true), data, "setting/external")
}
