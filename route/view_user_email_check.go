package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_user_email_check(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !user_auth(db, config) {
		return tool.Get_redirect("/login")
	}
	key, _ := config.Session.Get("c_key").(string)
	email, _ := config.Session.Get("c_email").(string)
	if key == "" || email == "" {
		return tool.Get_redirect("/change/email")
	}
	if values != nil {
		if values.Get("key") != key {
			return tool.Get_error_page(db, config, "key error")
		}
		user_delete(db, config.IP, "email")
		user_save(db, config.IP, "email", email)
		config.Session.Delete("c_key")
		config.Session.Delete("c_email")
		_ = config.Session.Save()
		return tool.Get_redirect("/change")
	}
	instruction := user_other(db, "check_key_text")
	body := ""
	if instruction != "" {
		body += tool.HTML_escape(instruction) + `<hr class="main_hr">`
	}
	body += `<form method="post"><input placeholder="` + tool.Get_language(db, "key", true) + `" name="key" type="text"><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`
	return user_form_page(db, config, tool.Get_language(db, "check_key", true), body)
}
