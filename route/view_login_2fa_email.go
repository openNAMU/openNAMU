package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_login_2fa_email(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	login_id, _ := config.Session.Get("login_id").(string)
	if login_id == "" {
		if email_id, ok := config.Session.Get("b_id").(string); ok && email_id != "" {
			login_id = email_id
			config.Session.Set("login_id", login_id)
			_ = config.Session.Save()
		}
	}
	if login_id == "" {
		return tool.Get_redirect("/login")
	}
	if user_value(db, login_id, "2fa") != "email" {
		return tool.Get_redirect("/login/2fa")
	}
	email := user_value(db, login_id, "email")
	if email == "" {
		return tool.Get_error_page(db, config, "not found")
	}

	key, _ := config.Session.Get("login_2fa_key").(string)
	if values != nil {
		captcha := tool.Captcha_response(values.Get("g-recaptcha"), values.Get("g-recaptcha-response"), values.Get("h-captcha-response"), values.Get("cf-turnstile-response"), values.Get("altcha"))
		if !tool.Captcha_check(db, config.Session, config.IP, captcha) {
			return tool.Get_error_page(db, config, "recaptcha")
		}
		if !tool.Get_auth_info(db, config.IP)["login_available"] || !tool.Get_auth_info(db, login_id)["login_available"] {
			return tool.Get_error_page(db, config, "ban")
		}
		if key == "" || values.Get("key") != key {
			return tool.Get_error_page(db, config, "key error")
		}
		config.Session.Delete("login_id")
		config.Session.Delete("b_id")
		config.Session.Delete("login_2fa_key")
		Api_record_user_agent_post(config, login_id)
		config.Session.Set("id", login_id)
		_ = config.Session.Save()
		return tool.Get_redirect("/user")
	}

	if key == "" {
		key = tool.Get_random_key(32)
		body := tool.Get_language(db, "key", true) + " : " + key
		if err := tool.Send_email(db, config.IP, email, tool.Get_language(db, "2fa", true), body); err != nil {
			config.Session.Delete("login_2fa_key")
			_ = config.Session.Save()
			return tool.Get_error_page(db, config, "email error")
		}
		config.Session.Set("login_2fa_key", key)
		if err := config.Session.Save(); err != nil {
			config.Session.Delete("login_2fa_key")
			_ = config.Session.Save()
			return tool.Get_error_page(db, config, "error")
		}
	}

	body := `<form method="post"><input name="key" type="text">` + tool.Get_captcha_ui(db, config) + `<button type="submit">` + tool.Get_language(db, "login", true) + `</button>` + tool.Get_http_warning(db) + `</form>`
	return user_form_page(db, config, tool.Get_language(db, "2fa", true), body)
}
