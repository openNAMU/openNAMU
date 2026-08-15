package route

import (
	"database/sql"
	"net/url"
	"strings"

	"opennamu/route/tool"
)

func View_login_login_post_full(config tool.Config, id string, password string, captcha string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	ban_data := tool.Get_user_ban(db, config.IP, "login")
	if ban_data[0] == "true" {
		return tool.Get_error_page(db, config, "ban")
	}
	if !tool.Captcha_check(db, config.Session, config.IP, captcha) {
		return tool.Get_error_page(db, config, "recaptcha")
	}
	if !tool.Password_check(db, id, password) {
		return tool.Get_error_page(db, config, "password error")
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

func View_login_2fa(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	login_id, _ := config.Session.Get("login_id").(string)
	if login_id == "" {
		return tool.Get_redirect("/login")
	}
	if values != nil {
		if !tool.Captcha_check(db, config.Session, config.IP, tool.Captcha_response(values.Get("g-recaptcha"), values.Get("g-recaptcha-response"), values.Get("h-captcha-response"), values.Get("cf-turnstile-response"))) {
			return tool.Get_error_page(db, config, "recaptcha")
		}
		ban_data := tool.Get_user_ban(db, config.IP, "login")
		if len(ban_data) > 0 && ban_data[0] == "true" {
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
		tool.Record_user_agent(db, login_id, config.IP, config.UserAgent, tool.Get_time())
		config.Session.Set("id", login_id)
		_ = config.Session.Save()
		return tool.Get_redirect("/user")
	}
	body := "<form method='post'><input type='password' name='pw'>" + tool.Get_captcha_ui(db, config) + "<button type='submit'>login</button>" + tool.Get_http_warning(db) + "</form>"
	return user_form_page(db, config, "2FA", body)
}

func View_login_2fa_email(config tool.Config, values url.Values) string {
	login_id, _ := config.Session.Get("login_id").(string)
	if login_id == "" {
		if email_id, ok := config.Session.Get("b_id").(string); ok && email_id != "" {
			config.Session.Set("login_id", email_id)
			_ = config.Session.Save()
		}
	}
	return View_login_2fa(config, values)
}

func View_login_find(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if values != nil {
		captcha := tool.Captcha_response(values.Get("g-recaptcha"), values.Get("g-recaptcha-response"), values.Get("h-captcha-response"), values.Get("cf-turnstile-response"))
		if !tool.Captcha_check(db, config.Session, config.IP, captcha) {
			return tool.Get_error_page(db, config, "recaptcha")
		}
		key := values.Get("key")
		if key != "" {
			user_id := ""
			if tool.QueryRow_DB(db, "select id from user_set where name = 'random_key' and data = ?", []any{&user_id}, key) {
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

func user_email_allowed(db *sql.DB, email string) bool {
	at_index := strings.LastIndex(email, "@")
	if at_index <= 0 || at_index == len(email)-1 {
		return false
	}
	domain := strings.TrimSpace(email[at_index+1:])
	rows := tool.Query_DB(db, "select html from html_filter where kind = 'email'")
	defer rows.Close()
	for rows.Next() {
		allowed_domain := ""
		if rows.Scan(&allowed_domain) != nil {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(allowed_domain), domain) {
			return true
		}
	}
	return false
}

func View_login_find_email(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if values != nil {
		captcha := tool.Captcha_response(values.Get("g-recaptcha"), values.Get("g-recaptcha-response"), values.Get("h-captcha-response"), values.Get("cf-turnstile-response"))
		if !tool.Captcha_check(db, config.Session, config.IP, captcha) {
			return tool.Get_error_page(db, config, "recaptcha")
		}
		email := strings.TrimSpace(values.Get("email"))
		user_id := ""
		if !tool.QueryRow_DB(db, "select id from user_set where name = 'email' and data = ?", []any{&user_id}, email) {
			return tool.Get_error_page(db, config, "not found")
		}
		key := tool.Get_random_key(32)
		config.Session.Set("email_reset_key", key)
		config.Session.Set("email_reset_id", user_id)
		_ = config.Session.Save()
		if err := tool.Send_email(db, config.IP, email, "password reset", key); err != nil {
			config.Session.Delete("email_reset_key")
			config.Session.Delete("email_reset_id")
			_ = config.Session.Save()
			return tool.Get_error_page(db, config, "email error")
		}
		return tool.Get_redirect("/login/find/email/check")
	}
	return user_form_page(db, config, "password search", "<form method='post'><input type='email' name='email'>"+tool.Get_captcha_ui(db, config)+"<button type='submit'>send</button></form>")
}

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
