package route

import (
	"database/sql"
	"opennamu/route/tool"
)

func user_register_access(db *sql.DB, config tool.Config) string {
	owner_auth := tool.Check_permission(db, "user_manage", config.IP)
	if !owner_auth && !tool.IP_or_user(config.IP) {
		return "login user"
	}
	if !owner_auth && user_other(db, "reg") == "on" {
		return "register disabled"
	}
	if !tool.Get_auth_info(db, config.IP)["register_available"] {
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
	password_length_limit := tool.Get_setting_value(db, "password_min_length", "", "0")
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

	owner_auth := tool.Check_permission(db, "user_manage", config.IP)
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
	return tool.Get_setting_value(db, name, "", "")
}
