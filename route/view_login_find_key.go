package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_login_find_key(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	user_id, _ := config.Session.Get("reset_id").(string)
	if user_id == "" {
		return tool.Get_redirect("/login/find")
	}
	if values != nil {
		captcha := tool.Captcha_response(values.Get("g-recaptcha"), values.Get("g-recaptcha-response"), values.Get("h-captcha-response"), values.Get("cf-turnstile-response"))
		if !tool.Captcha_check(db, config.Session, config.IP, captcha) {
			return tool.Get_error_page(db, config, "recaptcha")
		}
		password := values.Get("password")
		if password == "" || password != values.Get("password_check") {
			return tool.Get_error_page(db, config, "password different")
		}
		user_save(db, user_id, "pw", tool.Password_encode(db, password, tool.Get_user_encode(db, user_id)))
		user_delete(db, user_id, "2fa")
		user_delete(db, user_id, "2fa_pw")
		user_delete(db, user_id, "2fa_pw_encode")
		user_delete(db, user_id, "random_key")
		config.Session.Delete("reset_id")
		_ = config.Session.Save()
		return tool.Get_redirect("/login")
	}
	return user_form_page(db, config, "password change", `<form method="post"><input type="password" name="password"><input type="password" name="password_check">`+tool.Get_captcha_ui(db, config)+`<button type="submit">save</button></form>`)
}
