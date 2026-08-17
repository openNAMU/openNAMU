package route

import "opennamu/route/tool"

func View_alarm_read(config tool.Config, user_name string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if user_name == "" {
		user_name = config.IP
	}
	if user_name != config.IP && !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	tool.Exec_DB(db, "update user_notice set readme = '1' where name = ?", user_name)
	return tool.Get_redirect("/alarm/" + tool.Url_parser(user_name))
}
