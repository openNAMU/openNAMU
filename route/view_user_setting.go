package route

import (
	"database/sql"
	"net/url"

	"opennamu/route/tool"
)

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
		api_data := Api_user_setting_post(config, values)
		response, _ := api_data["response"].(string)
		if response == "require auth" {
			return tool.Get_error_page(db, config, "auth")
		}
		if response != "ok" {
			return tool.Get_error_page(db, config, "error")
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
	body += `<a href="/change/head">(` + tool.Get_language(db, "user_head", false) + `)</a> <a href="/change/top_menu">(` + tool.Get_language(db, "user_added_menu", true) + `)</a><hr class="main_hr"><span>` + tool.Get_language(db, "skin", true) + `</span><hr class="main_hr"><select name="skin">` + skin_options(current_skin) + `</select><hr class="main_hr">`
	body += `<a href="/change/skin_set">(` + tool.Get_language(db, "skin_set", true) + `)</a> <a href="/change/skin_set/main">(` + tool.Get_language(db, "main_skin_set", true) + `)</a><hr class="main_hr"><span>` + tool.Get_language(db, "language", true) + `</span><hr class="main_hr"><select name="lang">` + language_options(current_language) + `</select><hr class="main_hr"><span>` + tool.Get_language(db, "user_title", true) + `</span><hr class="main_hr"><select name="user_title">` + title_options + `</select><h2>` + tool.Get_language(db, "2fa", true) + `</h2><select name="2fa">` + twofa_options + `</select><hr class="main_hr"><input type="password" name="2fa_pw" placeholder="` + tool.Get_language(db, twofa_password, true) + `"><h2>` + tool.Get_language(db, "main_user_name", true) + `</h2><a href="/change/user_name">(` + tool.Get_language(db, "change_user_name", true) + `)</a><hr class="main_hr">`
	body += tool.Get_language(db, "user_name", true) + ` : ` + tool.HTML_escape(user_name) + `<h2>` + tool.Get_language(db, "sub_user_name", true) + `</h2><input name="sub_user_name" value="` + tool.HTML_escape(user_value(db, config.IP, "sub_user_name")) + `" placeholder="` + tool.Get_language(db, "sub_user_name", true) + `"><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "save", true) + `</button>` + tool.Get_http_warning(db) + `</form>`
	return user_form_page(db, config, tool.Get_language(db, "user_setting", true), body)
}

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
