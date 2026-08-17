package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_user_password(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !user_auth(db, config) {
		return tool.Get_redirect("/user")
	}
	if values != nil {
		current_password := values.Get("password_now")
		password := values.Get("password_new")
		password_repeat := values.Get("password_new_repeat")
		if password == "" && values.Has("password") {
			password = values.Get("password")
		}
		if password_repeat == "" && values.Has("password_check") {
			password_repeat = values.Get("password_check")
		}
		if password == "" {
			return tool.Get_error_page(db, config, "password empty")
		}
		if password != password_repeat {
			return tool.Get_error_page(db, config, "password different")
		}
		minimum := user_other(db, "password_min_length")
		if minimum != "" && tool.Get_len(password) < tool.Str_to_int(minimum) {
			return tool.Get_error_page(db, config, "password too short")
		}
		if !tool.Password_check(db, config.IP, current_password) {
			return tool.Get_error_page(db, config, "password error")
		}
		user_save(db, config.IP, "pw", tool.Password_encode(db, password, tool.Get_user_encode(db, config.IP)))
		return tool.Get_redirect("/user")
	}
	minimum := user_other(db, "password_min_length")
	minimum_text := ""
	if minimum != "" {
		minimum_text = " (" + tool.Get_language(db, "password_min_length", true) + " : " + tool.HTML_escape(minimum) + ")"
	}
	body := `<form method="post"><input placeholder="` + tool.Get_language(db, "now_password", true) + `" name="password_now" type="password"><hr class="main_hr"><input placeholder="` + tool.Get_language(db, "new_password", true) + minimum_text + `" name="password_new" type="password"><hr class="main_hr"><input placeholder="` + tool.Get_language(db, "password_confirm", true) + `" name="password_new_repeat" type="password"><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "save", true) + `</button>` + tool.Get_http_warning(db) + `</form>`
	return user_form_page(db, config, tool.Get_language(db, "password_change", true), body)
}
