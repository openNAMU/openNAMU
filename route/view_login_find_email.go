package route

import (
	"net/url"
	"opennamu/route/tool"
	"strings"
)

func View_login_find_email(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if values != nil {
		captcha := tool.Captcha_response(values.Get("g-recaptcha"), values.Get("g-recaptcha-response"), values.Get("h-captcha-response"), values.Get("cf-turnstile-response"), values.Get("altcha"))
		if !tool.Captcha_check(db, config.Session, config.IP, captcha) {
			return tool.Get_error_page(db, config, "recaptcha")
		}
		email := strings.TrimSpace(values.Get("email"))
		user_id, user_exists := tool.Get_user_set_id(db, "email", email)
		if !user_exists {
			return tool.Get_error_page(db, config, "not found")
		}
		key := tool.Get_random_key(32)
		config.Session.Set("email_reset_key", key)
		config.Session.Set("email_reset_id", user_id)
		_ = config.Session.Save()
		body := tool.Get_language(db, "id", true) + " : " + user_id + "\n\n" + tool.Get_language(db, "key", true) + " : " + key
		if err := tool.Send_email(db, config.IP, email, "password reset", body); err != nil {
			config.Session.Delete("email_reset_key")
			config.Session.Delete("email_reset_id")
			_ = config.Session.Save()
			return tool.Get_error_page(db, config, "email error")
		}
		return tool.Get_redirect("/login/find/email/check")
	}
	return user_form_page(db, config, "password search", "<form method='post'><input type='email' name='email'>"+tool.Get_captcha_ui(db, config)+"<button type='submit'>send</button></form>")
}
