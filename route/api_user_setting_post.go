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
	if values.Has("2fa") {
		if values.Get("2fa") == "" {
			user_delete(db, config.IP, "2fa")
			user_delete(db, config.IP, "2fa_pw")
			user_delete(db, config.IP, "2fa_pw_encode")
		} else {
			user_save(db, config.IP, "2fa", "on")
			if password := values.Get("2fa_pw"); password != "" {
				encode := tool.Get_user_encode(db, config.IP)
				user_save(db, config.IP, "2fa_pw", tool.Password_encode(db, password, encode))
				user_save(db, config.IP, "2fa_pw_encode", encode)
			}
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
