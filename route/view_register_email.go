package route

import (
	"net/url"
	"opennamu/route/tool"
	"strings"
)

func View_register_email(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	id, _ := config.Session.Get("reg_id").(string)
	pw, _ := config.Session.Get("reg_pw").(string)
	if id == "" || pw == "" {
		return tool.Get_redirect("/register")
	}
	if values != nil {
		email := strings.TrimSpace(values.Get("email"))
		if !user_email_allowed(db, email) {
			return tool.Get_error_page(db, config, "email domain")
		}
		_, email_exists := tool.Get_user_set_id(db, "email", email)
		if email_exists {
			return tool.Get_error_page(db, config, "email already exist")
		}
		key := tool.Get_random_key(32)
		title := user_other(db, "email_title")
		if title == "" {
			title = tool.Get_language(db, "register", true)
		}
		body := user_other(db, "email_text")
		if strings.Contains(body, "{}") {
			body = strings.ReplaceAll(body, "{}", key)
		} else {
			if body != "" {
				body += "\n\n"
			}
			body += tool.Get_language(db, "key", true) + " : " + key
		}
		if err := tool.Send_email(db, config.IP, email, title, body); err != nil {
			return tool.Get_error_page(db, config, "email error")
		}
		config.Session.Set("reg_email", email)
		config.Session.Set("reg_key", key)
		_ = config.Session.Save()
		return tool.Get_redirect("/register/email/check")
	}
	instruction := user_other(db, "email_insert_text")
	body := ""
	if instruction != "" {
		body += tool.HTML_escape(instruction) + `<hr class="main_hr">`
	}
	body += `<a href="/filter/email_filter">(` + tool.Get_language(db, "email_filter_list", true) + `)</a><hr class="main_hr"><form method="post"><input placeholder="` + tool.Get_language(db, "email", true) + `" name="email" type="email"><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`
	return user_form_page(db, config, tool.Get_language(db, "email", true), body)
}
