package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_login_2fa(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	login_id, _ := config.Session.Get("login_id").(string)
	if login_id == "" {
		return tool.Get_redirect("/login")
	}
	if values != nil {
		if !tool.Captcha_check(db, config.Session, config.IP, tool.Captcha_response(values.Get("g-recaptcha"), values.Get("g-recaptcha-response"), values.Get("h-captcha-response"), values.Get("cf-turnstile-response"), values.Get("altcha"))) {
			return tool.Get_error_page(db, config, "recaptcha")
		}
		if !tool.Get_auth_info(db, config.IP)["login_available"] || !tool.Get_auth_info(db, login_id)["login_available"] {
			return tool.Get_error_page(db, config, "ban")
		}
		stored := user_value(db, login_id, "2fa_pw")
		encode := user_value(db, login_id, "2fa_pw_encode")
		if encode == "" {
			encode = tool.Get_user_encode(db, login_id)
		}
		if stored != "" && stored != tool.Password_encode(db, values.Get("pw"), encode) {
			return tool.Get_error_page(db, config, "password error")
		}
		config.Session.Delete("login_id")
		config.Session.Delete("b_id")
		Api_record_user_agent_post(config, login_id)
		config.Session.Set("id", login_id)
		_ = config.Session.Save()
		return tool.Get_redirect("/user")
	}
	body := "<form method='post'><input type='password' name='pw'>" + tool.Get_captcha_ui(db, config) + "<button type='submit'>" + tool.Get_language(db, "login", true) + "</button>" + tool.Get_http_warning(db) + "</form>"
	return user_form_page(db, config, tool.Get_language(db, "2fa", true), body)
}
