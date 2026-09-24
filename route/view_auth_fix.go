package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_auth_fix(config tool.Config, user_name string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Get_user_set_exists(db, user_name, "pw") {
		return tool.Get_error_page(db, config, "error")
	}
	if values == nil && !tool.Check_permission(db, "auth_fix", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	if values != nil {
		result := Api_auth_fix_post(config, user_name, values)
		if result["response"] == "require auth" {
			return tool.Get_error_page(db, config, "auth")
		}
		if result["response"] == "password different" {
			return tool.Get_error_page(db, config, "password different")
		}
		if result["response"] != "ok" {
			return tool.Get_error_page(db, config, "error")
		}
		return tool.Get_redirect("/user/" + tool.Url_parser(user_name))
	}

	data := `<form method="post"><div id="opennamu_get_user_info">` + tool.HTML_escape(user_name) + `</div><hr class="main_hr">` +
		`<div><label for="auth_fix_action">` + tool.Get_language(db, "option", true) + `</label><select id="auth_fix_action" name="select"><option value="password_change">` + tool.Get_language(db, "password_change", true) + `</option><option value="2fa_password_change">` + tool.Get_language(db, "2fa_password_change", true) + `</option><option value="2fa_off">` + tool.Get_language(db, "2fa_off", true) + `</option></select></div><hr class="main_hr">` +
		`<div><label for="new_password">` + tool.Get_language(db, "new_password", true) + `</label><input id="new_password" name="new_password" type="password"></div><hr class="main_hr"><div><label for="password_check">` + tool.Get_language(db, "password_confirm", true) + `</label><input id="password_check" name="password_check" type="password"></div><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "go", true) + `</button></form>`
	return tool.Get_template(db, config, tool.Get_language(db, "user_fix", true), data, []any{}, [][]any{{"manager", tool.Get_language(db, "return", true)}}, map[string]string{})
}
