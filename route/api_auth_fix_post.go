package route

import (
	"database/sql"
	"net/url"

	"opennamu/route/tool"
)

func Api_auth_fix_post(config tool.Config, user_name string, values url.Values) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "auth_fix", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	password := ""
	if !tool.QueryRow_DB(db, "select data from user_set where id = ? and name = 'pw'", []any{&password}, user_name) {
		return_data["response"] = "error"
		return return_data
	}

	choice := values.Get("select")
	if (choice == "password_change" || choice == "2fa_password_change") && values.Get("new_password") != values.Get("password_check") {
		return_data["response"] = "password different"
		return return_data
	}

	encode := ""
	hash := ""
	set_name := "pw"
	if choice == "password_change" || choice == "2fa_password_change" {
		encode = tool.Get_user_encode(db, user_name)
		hash = tool.Password_encode(db, values.Get("new_password"), encode)
		if choice == "2fa_password_change" {
			set_name = "2fa_pw"
		}
	}

	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		if choice == "password_change" || choice == "2fa_password_change" {
			user_save(tx, user_name, set_name, hash)
			if choice == "2fa_password_change" {
				user_save(tx, user_name, "2fa_pw_encode", encode)
				user_save(tx, user_name, "2fa", "on")
			}
		} else if choice == "2fa_off" {
			user_delete(tx, user_name, "2fa")
			user_delete(tx, user_name, "2fa_pw")
			user_delete(tx, user_name, "2fa_pw_encode")
		}
		tool.Do_insert_auth_history(tx, config.IP, "user_fix ("+user_name+")")
		return nil
	}); err != nil {
		panic(err)
	}

	return_data["response"] = "ok"
	return return_data
}
