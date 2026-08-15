package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_auth_fix(config tool.Config, user_name string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	password := ""
	if !tool.QueryRow_DB(db, "select data from user_set where id = ? and name = 'pw'", []any{&password}, user_name) {
		return tool.Get_error_page(db, config, "error")
	}
	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	if values != nil {
		choice := values.Get("select")
		if choice == "password_change" || choice == "2fa_password_change" {
			if values.Get("new_password") != values.Get("password_check") {
				return tool.Get_error_page(db, config, "password different")
			}
			encode := tool.Get_user_encode(db, user_name)
			hash := tool.Password_encode(db, values.Get("new_password"), encode)
			set_name := "pw"
			if choice == "2fa_password_change" {
				set_name = "2fa_pw"
			}
			old := ""
			if tool.QueryRow_DB(db, "select data from user_set where id = ? and name = ?", []any{&old}, user_name, set_name) {
				tool.Exec_DB(db, "update user_set set data = ? where id = ? and name = ?", hash, user_name, set_name)
			} else {
				tool.Exec_DB(db, "insert into user_set (id, name, data) values (?, ?, ?)", user_name, set_name, hash)
			}
			if choice == "2fa_password_change" {
				user_save(db, user_name, "2fa_pw_encode", encode)
				user_save(db, user_name, "2fa", "on")
			}
		} else if choice == "2fa_off" {
			user_delete(db, user_name, "2fa")
			user_delete(db, user_name, "2fa_pw")
			user_delete(db, user_name, "2fa_pw_encode")
		}
		tool.Do_insert_auth_history(db, config.IP, "user_fix ("+user_name+")")
		return tool.Get_redirect("/user/" + tool.Url_parser(user_name))
	}

	data := `<form method="post"><div id="opennamu_get_user_info">` + tool.HTML_escape(user_name) + `</div><hr class="main_hr">` +
		`<select name="select"><option value="password_change">` + tool.Get_language(db, "password_change", true) + `</option><option value="2fa_password_change">` + tool.Get_language(db, "2fa_password_change", true) + `</option><option value="2fa_off">` + tool.Get_language(db, "2fa_off", true) + `</option></select><hr class="main_hr">` +
		`<input name="new_password" type="password" placeholder="` + tool.Get_language(db, "new_password", true) + `"><hr class="main_hr"><input name="password_check" type="password" placeholder="` + tool.Get_language(db, "password_confirm", true) + `"><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "go", true) + `</button></form>`
	return tool.Get_template(db, config, tool.Get_language(db, "user_fix", true), data, []any{}, [][]any{{"manager", tool.Get_language(db, "return", true)}}, map[string]string{})
}
