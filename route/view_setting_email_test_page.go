package route

import "opennamu/route/tool"

func View_setting_email_test(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "setting_email_test", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	data := `<form method="post">`
	data += Setting_input("title", "", "text") + Main_hr()
	data += Setting_input("email", "", "email") + Main_hr()
	data += Setting_textarea("data", "", "opennamu_textarea_500") + Main_hr()
	data += `<button type="submit">` + tool.Get_language(db, "send", true) + `</button></form>`

	return Setting_page(db, config, tool.Get_language(db, "email_test", true), data, "setting/external")
}
