package route

import (
	"database/sql"
	stdjson "encoding/json"
	"net/url"
	"strings"

	"opennamu/route/tool"
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
		var email_id string
		if tool.QueryRow_DB(db, "select id from user_set where name = 'email' and data = ?", []any{&email_id}, email) {
			return tool.Get_error_page(db, config, "email already exist")
		}
		key := tool.Get_random_key(32)
		title := user_other(db, "email_title")
		if title == "" {
			title = tool.Get_language(db, "register", true)
		}
		body := user_other(db, "email_text")
		if body != "" {
			body += "\n\n"
		}
		body += "Key : " + key
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

func View_register_submit(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	id, _ := config.Session.Get("submit_id").(string)
	pw, _ := config.Session.Get("submit_pw").(string)
	email, _ := config.Session.Get("submit_email").(string)
	if id == "" || pw == "" {
		return tool.Get_redirect("/register")
	}
	question := user_other(db, "approval_question")
	if question == "" {
		for _, name := range []string{"submit_id", "submit_pw", "submit_email"} {
			config.Session.Delete(name)
		}
		_ = config.Session.Save()
		return tool.Get_redirect("/register")
	}
	if values != nil {
		encode := tool.Get_main_encode(db)
		application, _ := stdjson.Marshal(map[string]string{
			"id":       id,
			"pw_hash":  tool.Password_encode(db, pw, encode),
			"email":    email,
			"encode":   encode,
			"question": question,
			"answer":   values.Get("answer"),
		})
		tool.Exec_DB(db, "delete from user_set where id = ? and name = 'application'", id)
		tool.Exec_DB(db, "insert into user_set (id, name, data) values (?, 'application', ?)", id, string(application))
		for _, name := range []string{"submit_id", "submit_pw", "submit_email"} {
			config.Session.Delete(name)
		}
		_ = config.Session.Save()
		return user_form_page(db, config, "register submit", "<p>submitted</p><a href='/user'>return</a>")
	}
	body := "<form method='post'><p>" + tool.HTML_escape(question) + "</p><input name='answer'><button type='submit'>send</button></form>"
	return user_form_page(db, config, "register submit", body)
}

func user_register_access(db *sql.DB, config tool.Config) string {
	owner_auth := tool.Check_acl(db, "", "", "owner_auth", config.IP)
	if !owner_auth && !tool.IP_or_user(config.IP) {
		return "login user"
	}
	if !owner_auth && user_other(db, "reg") == "on" {
		return "register disabled"
	}
	ban_data := tool.Get_user_ban(db, config.IP, "register")
	if len(ban_data) > 0 && ban_data[0] == "true" {
		return "ban"
	}
	return ""
}

func user_register_validate(db *sql.DB, config tool.Config, id string, password string, password_check string) string {
	if error_name := user_register_access(db, config); error_name != "" {
		return error_name
	}
	if password != password_check {
		return "password error"
	}
	if password == "" {
		return "empty password"
	}
	if id == password {
		return "password same as id"
	}
	password_length_limit := "0"
	tool.QueryRow_DB(db, "select data from other where name = 'password_min_length'", []any{&password_length_limit})
	if tool.Get_len(password) < tool.Str_to_int(password_length_limit) {
		return "password too short"
	}
	if !tool.Get_user_name_check(db, id) {
		return "user name error"
	}
	return ""
}

func user_register_post(config tool.Config, id string, password string, password_check string, captcha string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !tool.Captcha_check(db, config.Session, config.IP, captcha) {
		return tool.Get_error_page(db, config, "recaptcha")
	}
	error_name := user_register_validate(db, config, id, password, password_check)
	if error_name != "" {
		return tool.Get_error_page(db, config, error_name)
	}

	owner_auth := tool.Check_acl(db, "", "", "owner_auth", config.IP)
	email_required := !owner_auth && user_other(db, "email_have") != ""
	approval_required := !owner_auth && user_other(db, "requires_approval") != ""
	if email_required {
		config.Session.Set("reg_id", id)
		config.Session.Set("reg_pw", password)
		_ = config.Session.Save()
		return tool.Get_redirect("/register/email")
	}
	if approval_required {
		config.Session.Set("submit_id", id)
		config.Session.Set("submit_pw", password)
		config.Session.Delete("submit_email")
		_ = config.Session.Save()
		return tool.Get_redirect("/register/submit")
	}
	result := Api_add_user(config, id, password, "", "")
	if result["response"] != "ok" {
		return tool.Get_error_page(db, config, "register error")
	}
	return tool.Get_redirect("/login")
}

func user_other(db *sql.DB, name string) string {
	value := ""
	tool.QueryRow_DB(db, "select data from other where name = ?", []any{&value}, name)
	return value
}
