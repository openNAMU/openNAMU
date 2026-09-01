package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_login_find(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if values != nil {
		captcha := tool.Captcha_response(values.Get("g-recaptcha"), values.Get("g-recaptcha-response"), values.Get("h-captcha-response"), values.Get("cf-turnstile-response"), values.Get("altcha"))
		if !tool.Captcha_check(db, config.Session, config.IP, captcha) {
			return tool.Get_error_page(db, config, "recaptcha")
		}
		key := values.Get("key")
		if key != "" {
			user_id, user_exists := tool.Get_user_set_id(db, "random_key", key)
			if user_exists {
				config.Session.Set("reset_id", user_id)
				_ = config.Session.Save()
				return tool.Get_redirect("/login/find/key")
			}
		}
		return tool.Get_error_page(db, config, "not found")
	}
	body := `<ul><li><a href="/login/find/email">` + tool.Get_language(db, "email", true) + `</a></li></ul><hr class="main_hr"><form method="post"><input name="key">` + tool.Get_captcha_ui(db, config) + `<button type="submit">next</button></form>`
	return user_form_page(db, config, "password search", body)
}
