package route

import (
	"database/sql"
	"net/url"

	"opennamu/route/tool"
)

func Api_user_setting_post(config tool.Config, values url.Values) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if tool.IP_or_user(config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	language_list := user_language_list(db)
	title_choices := user_title_list(db, config.IP)
	profile_image := ""
	profile_image_set := values.Has("profile_image")
	if profile_image_set {
		var valid bool
		profile_image, valid = tool.Get_user_profile_image_name(db, values.Get("profile_image"))
		if !valid {
			return_data["response"] = "error"
			return_data["data"] = "invalid file"
			return return_data
		}
	}

	twofa_mode := ""
	twofa_password := ""
	twofa_password_hash := ""
	twofa_encode := ""
	if values.Has("2fa") {
		twofa_mode = values.Get("2fa")
		switch twofa_mode {
		case "":
		case "on":
			twofa_password = values.Get("2fa_pw")
			if twofa_password == "" && user_value(db, config.IP, "2fa_pw") == "" {
				return_data["response"] = "error"
				return_data["data"] = "password empty"
				return return_data
			}
			if twofa_password != "" {
				twofa_encode = tool.Get_user_encode(db, config.IP)
				twofa_password_hash = tool.Password_encode(db, twofa_password, twofa_encode)
			}
		case "email":
			if user_value(db, config.IP, "email") == "" {
				return_data["response"] = "error"
				return_data["data"] = "not found"
				return return_data
			}
		default:
			return_data["response"] = "error"
			return_data["data"] = "invalid data"
			return return_data
		}
	} else if values.Has("2fa_pw") {
		twofa_password = values.Get("2fa_pw")
		if twofa_password != "" {
			twofa_encode = tool.Get_user_encode(db, config.IP)
			twofa_password_hash = tool.Password_encode(db, twofa_password, twofa_encode)
		}
	}

	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		if values.Has("skin") {
			skin := values.Get("skin")
			if tool.Arr_in_str(tool.Get_skin_list("", true), skin) {
				user_save(tx, config.IP, "skin", skin)
			}
		}
		if values.Has("lang") {
			for _, language := range language_list {
				if language.value == values.Get("lang") {
					user_save(tx, config.IP, "lang", language.value)
					break
				}
			}
		}
		if values.Has("user_title") {
			title := ""
			for _, choice := range title_choices {
				if choice.value == values.Get("user_title") {
					title = choice.value
					break
				}
			}
			user_save(tx, config.IP, "user_title", title)
		}
		for _, name := range []string{"sub_user_name", "top_menu"} {
			if values.Has(name) {
				user_save(tx, config.IP, name, values.Get(name))
			}
		}

		if profile_image_set {
			if profile_image == "" {
				user_delete(tx, config.IP, "profile_image")
			} else {
				user_save(tx, config.IP, "profile_image", profile_image)
			}
		}

		if values.Has("2fa") {
			switch twofa_mode {
			case "":
				user_delete(tx, config.IP, "2fa")
				user_delete(tx, config.IP, "2fa_pw")
				user_delete(tx, config.IP, "2fa_pw_encode")
			case "on":
				if twofa_password != "" {
					user_save(tx, config.IP, "2fa_pw", twofa_password_hash)
					user_save(tx, config.IP, "2fa_pw_encode", twofa_encode)
				}
				user_save(tx, config.IP, "2fa", "on")
			case "email":
				user_save(tx, config.IP, "2fa", "email")
				user_delete(tx, config.IP, "2fa_pw")
				user_delete(tx, config.IP, "2fa_pw_encode")
			}
		} else if values.Has("2fa_pw") {
			if twofa_password == "" {
				user_delete(tx, config.IP, "2fa_pw")
				user_delete(tx, config.IP, "2fa_pw_encode")
				user_delete(tx, config.IP, "2fa")
			} else {
				user_save(tx, config.IP, "2fa_pw", twofa_password_hash)
				user_save(tx, config.IP, "2fa_pw_encode", twofa_encode)
				user_save(tx, config.IP, "2fa", "on")
			}
		}

		return nil
	}); err != nil {
		panic(err)
	}

	return_data["response"] = "ok"
	return return_data
}
