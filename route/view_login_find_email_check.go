package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_login_find_email_check(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	key, _ := config.Session.Get("email_reset_key").(string)
	if values != nil {
		captcha := tool.Captcha_response(values.Get("g-recaptcha"), values.Get("g-recaptcha-response"), values.Get("h-captcha-response"), values.Get("cf-turnstile-response"))
		if !tool.Captcha_check(db, config.Session, config.IP, captcha) {
			return tool.Get_error_page(db, config, "recaptcha")
		}
	}
	if values != nil && key != "" && values.Get("key") == key {
		reset_id, _ := config.Session.Get("email_reset_id").(string)
		config.Session.Set("reset_id", reset_id)
		config.Session.Delete("email_reset_key")
		config.Session.Delete("email_reset_id")
		_ = config.Session.Save()
		return tool.Get_redirect("/login/find/key")
	}
	return user_form_page(db, config, "email check", "<form method='post'><input name='key'>"+tool.Get_captcha_ui(db, config)+"<button type='submit'>check</button></form>")
}
