package route

import "opennamu/route/tool"

func View_login_login_post_full(config tool.Config, id string, password string, captcha string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !tool.Get_auth_info(db, config.IP)["login_available"] {
		return tool.Get_error_page(db, config, "ban")
	}
	if !tool.Captcha_check(db, config.Session, config.IP, captcha) {
		return tool.Get_error_page(db, config, "recaptcha")
	}
	if !tool.Password_check(db, id, password) {
		return tool.Get_error_page(db, config, "password error")
	}
	if !tool.Get_auth_info(db, id)["login_available"] {
		return tool.Get_error_page(db, config, "ban")
	}
	if user_value(db, id, "2fa") != "" || user_value(db, id, "2fa_pw") != "" {
		config.Session.Set("login_id", id)
		_ = config.Session.Save()
		return tool.Get_redirect("/login/2fa")
	}
	config.Session.Set("id", id)
	_ = config.Session.Save()
	tool.Record_user_agent(db, id, config.IP, config.UserAgent, tool.Get_time())
	return tool.Get_redirect("/user")
}
