package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_user_password(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if values == nil && !user_auth(db, config) {
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
		api_data := Api_user_password_post(config, current_password, password, password_repeat)
		response, _ := api_data["response"].(string)
		if response == "require auth" {
			return tool.Get_error_page(db, config, "auth")
		}
		if response != "ok" {
			error_name, _ := api_data["data"].(string)
			if error_name == "" {
				error_name = "error"
			}
			return tool.Get_error_page(db, config, error_name)
		}
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
