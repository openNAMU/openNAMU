package route

import (
	"database/sql"
	"net/url"

	"opennamu/route/tool"
)

type user_choice struct {
	value string
	label string
}

func user_title_list(db *sql.DB, user_name string) []user_choice {
	choice_list := []user_choice{{"", tool.Get_language(db, "default", true)}, {"🌳", "🌳 newbie"}}
	challenge_list := []struct {
		name  string
		value string
		label string
	}{
		{"challenge_first_contribute", "🔰", "🔰 first_contribute"},
		{"challenge_tenth_contribute", "📝", "📝 tenth_contribute"},
		{"challenge_hundredth_contribute", "🖊️", "🖊️ hundredth_contribute"},
		{"challenge_thousandth_contribute", "🏅", "🏅 thousandth_contribute"},
		{"challenge_first_discussion", "💬", "💬 first_discussion"},
		{"challenge_tenth_discussion", "💡", "💡 tenth_discussion"},
		{"challenge_hundredth_discussion", "📢", "📢 hundredth_discussion"},
		{"challenge_thousandth_discussion", "📜", "📜 thousandth_discussion"},
	}
	for _, challenge := range challenge_list {
		if user_value(db, user_name, challenge.name) != "" {
			choice_list = append(choice_list, user_choice{challenge.value, challenge.label})
		}
	}
	if user_value(db, user_name, "challenge_admin") != "" {
		choice_list = append(choice_list, user_choice{"☑️", "☑️ before_admin"})
	}
	if tool.Check_acl(db, "", "", "all_admin_auth", user_name) {
		choice_list = append(choice_list, user_choice{"✅", "✅ admin"})
	}
	var egg string
	if tool.QueryRow_DB(db, "select name from user_set where id = ? and name = 'get_🥚'", []any{&egg}, user_name) {
		choice_list = append(choice_list, user_choice{"🥚", "🥚 easter_egg"})
	}
	return choice_list
}

func user_language_list(db *sql.DB) []user_choice {
	choice_list := []user_choice{{"default", tool.Get_language(db, "default", true)}}
	set_list := tool.Get_init_set_list("language")
	if language_set, ok := set_list["language"]; ok {
		if language_list, ok := language_set["list"].([]string); ok {
			for _, language := range language_list {
				choice_list = append(choice_list, user_choice{language, language})
			}
		}
	}
	return choice_list
}

func user_option(choice user_choice, current string) string {
	selected := ""
	if choice.value == current {
		selected = ` selected="selected"`
	}
	return `<option value="` + tool.HTML_escape(choice.value) + `"` + selected + `>` + tool.HTML_escape(choice.label) + `</option>`
}
func View_user_setting(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	language_list := user_language_list(db)
	skin_options := func(current string) string {
		if current == "" {
			current = "default"
		}
		skin_list := tool.Get_skin_list(current, true)
		if !tool.Arr_in_str(skin_list, current) {
			skin_list = append([]string{current}, skin_list...)
		}
		options := ""
		for _, skin := range skin_list {
			options += user_option(user_choice{skin, skin}, current)
		}
		return options
	}
	language_options := func(current string) string {
		if current == "" {
			current = "default"
		}
		options := ""
		for _, language := range language_list {
			options += user_option(language, current)
		}
		return options
	}

	if tool.IP_or_user(config.IP) {
		if values != nil {
			if values.Has("skin") {
				skin := values.Get("skin")
				if tool.Arr_in_str(tool.Get_skin_list("", true), skin) {
					config.Session.Set("skin", skin)
				}
			}
			if values.Has("lang") {
				for _, language := range language_list {
					if language.value == values.Get("lang") {
						config.Session.Set("lang", language.value)
						break
					}
				}
			}
			_ = config.Session.Save()
			return tool.Get_redirect("/change")
		}

		current_skin, _ := config.Session.Get("skin").(string)
		current_language, _ := config.Session.Get("lang").(string)
		body := `<form method="post"><div id="opennamu_get_user_info">` + tool.HTML_escape(config.IP) + `</div><hr class="main_hr"><h2>` + tool.Get_language(db, "main", true) + `</h2>`
		body += `<span>` + tool.Get_language(db, "skin", true) + `</span><hr class="main_hr"><select name="skin">` + skin_options(current_skin) + `</select><hr class="main_hr">`
		body += `<a href="/change/skin_set">(` + tool.Get_language(db, "skin_set", true) + `)</a> <a href="/change/skin_set/main">(` + tool.Get_language(db, "main_skin_set", true) + `)</a><hr class="main_hr">`
		body += `<span>` + tool.Get_language(db, "language", true) + `</span><hr class="main_hr"><select name="lang">` + language_options(current_language) + `</select><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "save", true) + `</button>` + tool.Get_http_warning(db) + `</form>`
		return user_form_page(db, config, tool.Get_language(db, "user_setting", true), body)
	}

	if values != nil {
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
		return tool.Get_redirect("/change")
	}

	current_skin := user_value(db, config.IP, "skin")
	current_language := user_value(db, config.IP, "lang")
	current_title := user_value(db, config.IP, "user_title")
	user_name := user_value(db, config.IP, "user_name")
	if user_name == "" {
		user_name = config.IP
	}
	email := user_value(db, config.IP, "email")
	if email == "" {
		email = "-"
	}
	random_key := user_value(db, config.IP, "random_key")
	if random_key == "" {
		random_key = "-"
	}
	twofa := user_value(db, config.IP, "2fa")
	twofa_password := "2fa_password"
	if user_value(db, config.IP, "2fa_pw") != "" {
		twofa_password = "2fa_password_change"
	}

	title_options := ""
	for _, choice := range user_title_list(db, config.IP) {
		title_options += user_option(choice, current_title)
	}
	twofa_options := user_option(user_choice{"", tool.Get_language(db, "off", true)}, twofa)
	twofa_options += user_option(user_choice{"on", tool.Get_language(db, "password", true)}, twofa)
	body := `<form method="post"><div id="opennamu_get_user_info">` + tool.HTML_escape(config.IP) + `</div><hr class="main_hr">`
	body += `<a href="/change/pw">(` + tool.Get_language(db, "password_change", true) + `)</a><hr class="main_hr">`
	body += `<span>` + tool.Get_language(db, "email", true) + ` : ` + tool.HTML_escape(email) + `</span> <a href="/change/email">(` + tool.Get_language(db, "email_change", true) + `)</a> <a href="/change/email/delete">(` + tool.Get_language(db, "email_delete", true) + `)</a><hr class="main_hr">`
	body += `<span>` + tool.Get_language(db, "password_instead_key", true) + ` : ` + tool.HTML_escape(random_key) + `</span> <a href="/change/key">(` + tool.Get_language(db, "key_change", true) + `)</a> <a href="/change/key/delete">(` + tool.Get_language(db, "key_delete", true) + `)</a><h2>` + tool.Get_language(db, "main", true) + `</h2>`
	body += `<a href="/change/head">(` + tool.Get_language(db, "user_head", true) + `)</a> <a href="/change/top_menu">(` + tool.Get_language(db, "user_added_menu", true) + `)</a><hr class="main_hr"><span>` + tool.Get_language(db, "skin", true) + `</span><hr class="main_hr"><select name="skin">` + skin_options(current_skin) + `</select><hr class="main_hr">`
	body += `<a href="/change/skin_set">(` + tool.Get_language(db, "skin_set", true) + `)</a> <a href="/change/skin_set/main">(` + tool.Get_language(db, "main_skin_set", true) + `)</a><hr class="main_hr"><span>` + tool.Get_language(db, "language", true) + `</span><hr class="main_hr"><select name="lang">` + language_options(current_language) + `</select><hr class="main_hr"><span>` + tool.Get_language(db, "user_title", true) + `</span><hr class="main_hr"><select name="user_title">` + title_options + `</select><h2>` + tool.Get_language(db, "2fa", true) + `</h2><select name="2fa">` + twofa_options + `</select><hr class="main_hr"><input type="password" name="2fa_pw" placeholder="` + tool.Get_language(db, twofa_password, true) + `"><h2>` + tool.Get_language(db, "main_user_name", true) + `</h2><a href="/change/user_name">(` + tool.Get_language(db, "change_user_name", true) + `)</a><hr class="main_hr">`
	body += tool.Get_language(db, "user_name", true) + ` : ` + tool.HTML_escape(user_name) + `<h2>` + tool.Get_language(db, "sub_user_name", true) + `</h2><input name="sub_user_name" value="` + tool.HTML_escape(user_value(db, config.IP, "sub_user_name")) + `" placeholder="` + tool.Get_language(db, "sub_user_name", true) + `"><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "save", true) + `</button>` + tool.Get_http_warning(db) + `</form>`
	return user_form_page(db, config, tool.Get_language(db, "user_setting", true), body)
}

func View_user_setting_field(config tool.Config, field string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !user_auth(db, config) {
		return tool.Get_redirect("/user")
	}
	if values != nil {
		value := values.Get("data")
		if field == "user_name" {
			if !tool.Get_user_name_check(db, value) && value != user_value(db, config.IP, "user_name") {
				return tool.Get_error_page(db, config, "user name error")
			}
		}
		if field == "email" && value == "" {
			user_delete(db, config.IP, field)
		} else {
			user_save(db, config.IP, field, value)
		}
		return tool.Get_redirect("/change")
	}
	value := user_value(db, config.IP, field)
	input_type := "text"
	if field == "email" {
		input_type = "email"
	}
	body := `<form method="post"><input type="` + input_type + `" name="data" value="` + tool.HTML_escape(value) + `"><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`
	return user_form_page(db, config, field, body)
}

func View_user_password(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !user_auth(db, config) {
		return tool.Get_redirect("/user")
	}
	if values != nil {
		current_password := values.Get("password_now")
		password := values.Get("password_new")
		password_repeat := values.Get("password_new_repeat")
		if password == "" && values.Has("password") {
			password = values.Get("password")
		}
		if password_repeat == "" && values.Has("password_check") {
			password_repeat = values.Get("password_check")
		}
		if password == "" {
			return tool.Get_error_page(db, config, "password empty")
		}
		if password != password_repeat {
			return tool.Get_error_page(db, config, "password different")
		}
		minimum := user_other(db, "password_min_length")
		if minimum != "" && tool.Get_len(password) < tool.Str_to_int(minimum) {
			return tool.Get_error_page(db, config, "password too short")
		}
		if !tool.Password_check(db, config.IP, current_password) {
			return tool.Get_error_page(db, config, "password error")
		}
		user_save(db, config.IP, "pw", tool.Password_encode(db, password, tool.Get_user_encode(db, config.IP)))
		return tool.Get_redirect("/user")
	}
	minimum := user_other(db, "password_min_length")
	minimum_text := ""
	if minimum != "" {
		minimum_text = " (" + tool.Get_language(db, "password_min_length", true) + " : " + tool.HTML_escape(minimum) + ")"
	}
	body := `<form method="post"><input placeholder="` + tool.Get_language(db, "now_password", true) + `" name="password_now" type="password"><hr class="main_hr"><input placeholder="` + tool.Get_language(db, "new_password", true) + minimum_text + `" name="password_new" type="password"><hr class="main_hr"><input placeholder="` + tool.Get_language(db, "password_confirm", true) + `" name="password_new_repeat" type="password"><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "save", true) + `</button>` + tool.Get_http_warning(db) + `</form>`
	return user_form_page(db, config, tool.Get_language(db, "password_change", true), body)
}

func View_user_key(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !user_auth(db, config) {
		return tool.Get_redirect("/user")
	}
	value := ""
	for value == "" {
		value = tool.Get_random_key(128)
		existing := ""
		if tool.QueryRow_DB(db, "select data from user_set where name = ? and data = ?", []any{&existing}, "random_key", value) {
			value = ""
		}
	}
	user_save(db, config.IP, "random_key", value)
	return tool.Get_redirect("/change")
}

func View_user_head(config tool.Config, values url.Values) string {
	return View_user_head_skin(config, "", values)
}

func View_user_top_menu(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !user_auth(db, config) {
		return tool.Get_redirect("/user")
	}
	if values != nil {
		content := values.Get("content")
		if !values.Has("content") {
			content = values.Get("data")
		}
		user_save(db, config.IP, "top_menu", content)
		return tool.Get_redirect("/change/top_menu")
	}
	content := user_value(db, config.IP, "top_menu")
	body := `<span>EX)<br>ONTS<br>https://2du.pythonanywhere.com/<br>FrontPage<br>/w/FrontPage</span><hr class="main_hr">` + tool.Get_language(db, "not_support_skin_warning", true) + `<hr class="main_hr"><form method="post"><textarea class="opennamu_textarea_500" placeholder="` + tool.Get_language(db, "enter_top_menu_setting", true) + `" name="content" id="content">` + tool.HTML_escape(content) + `</textarea><hr class="main_hr"><button id="opennamu_save_button" type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`
	return user_form_page(db, config, tool.Get_language(db, "user_added_menu", true), body)
}

func View_user_name(config tool.Config, values url.Values) string {
	return View_user_name_for(config, "", values)
}

func View_user_skin(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	skin_list := tool.Get_skin_list("ringo", true)
	if values != nil {
		skin := values.Get("skin")
		if !tool.Arr_in_str(skin_list, skin) {
			return tool.Get_redirect("/change")
		}
		if tool.IP_or_user(config.IP) {
			config.Session.Set("skin", skin)
			_ = config.Session.Save()
		} else {
			user_save(db, config.IP, "skin", skin)
		}
		return tool.Get_redirect("/change")
	}

	current := ""
	if tool.IP_or_user(config.IP) {
		current, _ = config.Session.Get("skin").(string)
	} else {
		current = user_value(db, config.IP, "skin")
	}
	if current == "" {
		current = "default"
	}

	body := `<form method="post"><select name="skin">`
	for _, skin := range skin_list {
		selected := ""
		if skin == current {
			selected = ` selected="selected"`
		}
		body += `<option value="` + tool.HTML_escape(skin) + `"` + selected + `>` + tool.HTML_escape(skin) + `</option>`
	}
	body += `</select><button type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`
	return user_form_page(db, config, "skin", body)
}
