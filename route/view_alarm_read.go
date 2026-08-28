package route

import "opennamu/route/tool"

func View_alarm_read(config tool.Config, user_name string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if user_name == "" {
		user_name = config.IP
	}
	result := Api_alarm_read_post(config, user_name)
	if result["response"] == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	}
	if result["response"] != "ok" {
		return tool.Get_error_page(db, config, "error")
	}
	return_path := "/alarm"
	if user_name != config.IP {
		return_path = "/alarm_user/" + tool.Url_parser(user_name)
	}
	return tool.Get_redirect(return_path)
}
