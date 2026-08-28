package route

import "opennamu/route/tool"

func View_alarm_delete(config tool.Config, user_name string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if user_name == "" {
		user_name = config.IP
	}
	if user_name != config.IP && !tool.Check_permission(db, "owner", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	return_path := "/alarm"
	if user_name != config.IP {
		return_path = "/alarm_user/" + tool.Url_parser(user_name)
	}
	body := `<form method="post"><span>` + tool.Get_language(db, "delete_warning", true) + `</span><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "delete", true) + `</button></form>`
	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, "delete", true),
		body,
		[]any{},
		[][]any{{return_path, tool.Get_language(db, "return", true)}},
		map[string]string{},
	)
}
