package route

import (
	"net/url"
	"opennamu/route/tool"
	"strings"
)

func View_user_email(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !user_auth(db, config) {
		return tool.Get_redirect("/login")
	}
	if values != nil {
		email := strings.TrimSpace(values.Get("email"))
		if !user_email_allowed(db, email) {
			return tool.Get_error_page(db, config, "email domain")
		}
		var email_id string
		if tool.QueryRow_DB(db, "select id from user_set where name = ? and data = ?", []any{&email_id}, "email", email) {
			return tool.Get_error_page(db, config, "email already exist")
		}
		key := tool.Get_random_key(32)
		title := user_other(db, "email_title")
		if title == "" {
			title = tool.Get_language(db, "email", true) + " key"
		}
		body := user_other(db, "email_text")
		if body != "" {
			body += "\n\n"
		}
		body += "Key : " + key
		if err := tool.Send_email(db, config.IP, email, title, body); err != nil {
			return tool.Get_error_page(db, config, "email error")
		}
		config.Session.Set("c_key", key)
		config.Session.Set("c_email", email)
		_ = config.Session.Save()
		return tool.Get_redirect("/change/email/check")
	}
	instruction := user_other(db, "email_insert_text")
	body := ""
	if instruction != "" {
		body += tool.HTML_escape(instruction) + `<hr class="main_hr">`
	}
	body += `<a href="/filter/email_filter">(` + tool.Get_language(db, "email_filter_list", true) + `)</a><hr class="main_hr"><form method="post"><input placeholder="` + tool.Get_language(db, "email", true) + `" name="email" type="email"><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`
	return user_form_page(db, config, tool.Get_language(db, "email", true), body)
}
