package route

import (
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
	if values.Has("skin") {
		skin := values.Get("skin")
		if tool.Arr_in_str(tool.Get_skin_list("", true), skin) {
			user_save(db, config.IP, "skin", skin)
		}
	}
	if values.Has("lang") {
		for _, language := range language_list {
			if language.value == values.Get("lang") {
				user_save(db, config.IP, "lang", language.value)
				break
			}
		}
	}
	if values.Has("user_title") {
		title := ""
		for _, choice := range user_title_list(db, config.IP) {
			if choice.value == values.Get("user_title") {
				title = choice.value
				break
			}
		}
		user_save(db, config.IP, "user_title", title)
	}
	for _, name := range []string{"sub_user_name", "top_menu"} {
		if values.Has(name) {
			user_save(db, config.IP, name, values.Get(name))
		}
	}
	if values.Has("profile_image") {
		profile_image, valid := tool.Get_user_profile_image_name(db, values.Get("profile_image"))
		if !valid {
			return_data["response"] = "error"
			return_data["data"] = "invalid file"
			return return_data
		}
		if profile_image == "" {
			user_delete(db, config.IP, "profile_image")
		} else {
			user_save(db, config.IP, "profile_image", profile_image)
		}
	}
	if values.Has("2fa") {
		switch values.Get("2fa") {
		case "":
			user_delete(db, config.IP, "2fa")
			user_delete(db, config.IP, "2fa_pw")
			user_delete(db, config.IP, "2fa_pw_encode")
		case "on":
			if password := values.Get("2fa_pw"); password != "" {
				encode := tool.Get_user_encode(db, config.IP)
				user_save(db, config.IP, "2fa_pw", tool.Password_encode(db, password, encode))
				user_save(db, config.IP, "2fa_pw_encode", encode)
			}
			if user_value(db, config.IP, "2fa_pw") == "" {
				return_data["response"] = "error"
				return_data["data"] = "password empty"
				return return_data
			}
			user_save(db, config.IP, "2fa", "on")
		case "email":
			if user_value(db, config.IP, "email") == "" {
				return_data["response"] = "error"
				return_data["data"] = "not found"
				return return_data
			}
			user_save(db, config.IP, "2fa", "email")
			user_delete(db, config.IP, "2fa_pw")
			user_delete(db, config.IP, "2fa_pw_encode")
		default:
			return_data["response"] = "error"
			return_data["data"] = "invalid data"
			return return_data
		}
	} else if values.Has("2fa_pw") {
		password := values.Get("2fa_pw")
		if password == "" {
			user_delete(db, config.IP, "2fa_pw")
			user_delete(db, config.IP, "2fa_pw_encode")
			user_delete(db, config.IP, "2fa")
		} else {
			encode := tool.Get_user_encode(db, config.IP)
			user_save(db, config.IP, "2fa_pw", tool.Password_encode(db, password, encode))
			user_save(db, config.IP, "2fa_pw_encode", encode)
			user_save(db, config.IP, "2fa", "on")
		}
	}

	return_data["response"] = "ok"
	return return_data
}
