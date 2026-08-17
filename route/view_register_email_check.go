package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_register_email_check(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	id, _ := config.Session.Get("reg_id").(string)
	pw, _ := config.Session.Get("reg_pw").(string)
	email, _ := config.Session.Get("reg_email").(string)
	key, _ := config.Session.Get("reg_key").(string)
	if id == "" || pw == "" || email == "" || key == "" {
		return tool.Get_redirect("/register")
	}
	if values == nil || values.Get("key") != key {
		if values != nil {
			return tool.Get_error_page(db, config, "key error")
		}
		instruction := user_other(db, "check_key_text")
		body := ""
		if instruction != "" {
			body += tool.HTML_escape(instruction) + `<hr class="main_hr">`
		}
		body += `<form method="post"><input placeholder="` + tool.Get_language(db, "key", true) + `" name="key" type="text"><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`
		return user_form_page(db, config, tool.Get_language(db, "check_key", true), body)
	}

	if user_other(db, "requires_approval") != "" {
		config.Session.Set("submit_id", id)
		config.Session.Set("submit_pw", pw)
		config.Session.Set("submit_email", email)
		for _, name := range []string{"reg_id", "reg_pw", "reg_email", "reg_key"} {
			config.Session.Delete(name)
		}
		_ = config.Session.Save()
		return tool.Get_redirect("/register/submit")
	}

	result := Api_add_user(config, id, pw, email, "")
	if result["response"] != "ok" {
		return tool.Get_error_page(db, config, "register error")
	}
	for _, name := range []string{"reg_id", "reg_pw", "reg_email", "reg_key"} {
		config.Session.Delete(name)
	}
	_ = config.Session.Save()
	return tool.Get_redirect("/login")
}
