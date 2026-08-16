package route

import (
	"database/sql"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"opennamu/route/tool"
)

func user_field_values(values url.Values, field string) url.Values {
	if values == nil {
		return nil
	}
	copy_values := url.Values{}
	for key, list := range values {
		copy_values[key] = append([]string{}, list...)
	}
	if copy_values.Get("data") == "" {
		if value := copy_values.Get(field); value != "" {
			copy_values.Set("data", value)
		}
	}
	if field == "user_name" && copy_values.Get("data") == "" {
		copy_values.Set("data", copy_values.Get("new_user_name"))
	}
	return copy_values
}

func View_user_setting_field_compat(config tool.Config, field string, values url.Values) string {
	return View_user_setting_field(config, field, user_field_values(values, field))
}

func View_user_key_delete(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !user_auth(db, config) {
		return tool.Get_redirect("/login")
	}
	if values != nil {
		user_delete(db, config.IP, "random_key")
		return tool.Get_redirect("/change")
	}
	body := `<form method="post"><span>` + tool.Get_language(db, "delete_warning", true) + `</span><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "delete", true) + `</button></form>`
	return user_form_page(db, config, tool.Get_language(db, "key_delete", true), body)
}

func View_user_field_delete(config tool.Config, field string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !user_auth(db, config) {
		return tool.Get_redirect("/login")
	}
	if values != nil {
		user_delete(db, config.IP, field)
		return tool.Get_redirect("/change")
	}
	title := tool.Get_language(db, field+"_delete", true)
	body := `<form method="post"><span>` + tool.Get_language(db, "delete_warning", true) + `</span><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "delete", true) + `</button></form>`
	return user_form_page(db, config, title, body)
}

func View_user_head_reset(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	skin_name := tool.Get_use_skin_name_session(db, config.IP, config.Session)
	if values != nil {
		if user_auth(db, config) {
			user_save(db, config.IP, "custom_css", "")
			user_save(db, config.IP, "custom_css_"+skin_name, "")
			user_delete(db, config.IP, "head")
		}
		config.Session.Set("head", "")
		config.Session.Set("head_"+skin_name, "")
		_ = config.Session.Save()
		return tool.Get_redirect("/change/head")
	}
	data := ""
	data_skin := ""
	if user_auth(db, config) {
		data = user_value(db, config.IP, "custom_css")
		data_skin = user_value(db, config.IP, "custom_css_"+skin_name)
	} else {
		data, _ = config.Session.Get("head").(string)
		data_skin, _ = config.Session.Get("head_" + skin_name).(string)
	}
	body := `<form method="post"><style>.main_hr { border: none; }</style>` + tool.Get_language(db, "all", true) + `<hr class="main_hr"><pre>` + tool.HTML_escape(data) + `</pre><hr class="main_hr">` + tool.HTML_escape(skin_name) + `<hr class="main_hr"><pre>` + tool.HTML_escape(data_skin) + `</pre><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "reset", true) + `</button></form>`
	return body
}

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

func View_user_email_check(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !user_auth(db, config) {
		return tool.Get_redirect("/login")
	}
	key, _ := config.Session.Get("c_key").(string)
	email, _ := config.Session.Get("c_email").(string)
	if key == "" || email == "" {
		return tool.Get_redirect("/change/email")
	}
	if values != nil {
		if values.Get("key") != key {
			return tool.Get_error_page(db, config, "key error")
		}
		user_delete(db, config.IP, "email")
		user_save(db, config.IP, "email", email)
		config.Session.Delete("c_key")
		config.Session.Delete("c_email")
		_ = config.Session.Save()
		return tool.Get_redirect("/change")
	}
	instruction := user_other(db, "check_key_text")
	body := ""
	if instruction != "" {
		body += tool.HTML_escape(instruction) + `<hr class="main_hr">`
	}
	body += `<form method="post"><input placeholder="` + tool.Get_language(db, "key", true) + `" name="key" type="text"><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`
	return user_form_page(db, config, tool.Get_language(db, "check_key", true), body)
}

func View_user_head_skin(config tool.Config, skin_name string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	storage_name := "custom_css"
	session_name := "head"
	if skin_name != "" {
		storage_name += "_" + skin_name
		session_name += "_" + skin_name
	}
	redirect_path := "/change/head"
	if skin_name != "" {
		redirect_path += "/" + tool.Url_parser(skin_name)
	}
	if values != nil {
		content := values.Get("content")
		if !values.Has("content") {
			content = values.Get("data")
		}
		if user_auth(db, config) {
			user_save(db, config.IP, storage_name, content)
		}
		config.Session.Set(session_name, content)
		_ = config.Session.Save()
		return tool.Get_redirect(redirect_path)
	}
	content, session_exists := config.Session.Get(session_name).(string)
	if !session_exists && user_auth(db, config) {
		content = user_value(db, config.IP, storage_name)
		if content == "" && skin_name == "" {
			content = user_value(db, config.IP, "head")
		}
	}
	body := ""
	if !user_auth(db, config) {
		body += `<span>` + tool.Get_language(db, "user_head_warning", true) + `</span><hr class="main_hr">`
	}
	body += `<a href="/change/head">(` + tool.Get_language(db, "all", true) + `)</a> `
	for _, skin := range tool.Get_skin_list("", true) {
		body += `<a href="/change/head/` + tool.Url_parser(skin) + `">(` + tool.HTML_escape(skin) + `)</a> `
	}
	sub_name := ""
	if skin_name != "" {
		sub_name = " (" + skin_name + ")"
	}
	body += `<hr class="main_hr"><span>&lt;style&gt;CSS&lt;/style&gt;<br>&lt;script&gt;JS&lt;/script&gt;</span><hr class="main_hr">`
	body += `<form method="post"><textarea class="opennamu_textarea_500" cols="100" name="content">` + tool.HTML_escape(content) + `</textarea><hr class="main_hr">` + tool.Get_language(db, "user_css_warning", true) + ` : <a href="/change/head_reset">/change/head_reset</a><hr class="main_hr"><button id="opennamu_save_button" type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`
	return user_form_page(db, config, tool.Get_language(db, "user_head", true)+sub_name, body)
}

func user_skin_choice(set_list map[string][][]string, field string, value string) bool {
	for _, choice := range set_list[field] {
		if len(choice) > 0 && choice[0] == value {
			return true
		}
	}
	return false
}

func user_skin_main_render_simple_set(db *sql.DB, data string) string {
	heading_regex := regexp.MustCompile("(?s)<h([1-6])>([^<>]+)</h[1-6]>")
	matches := heading_regex.FindAllStringSubmatch(data, -1)
	if len(matches) == 0 {
		return data
	}

	heading_stack := make([]int, 6)
	toc_data := strings.Builder{}
	toc_data.WriteString(`<div class="opennamu_TOC" id="toc"><span class="opennamu_TOC_title">` + tool.Get_language(db, "toc", true) + `</span><br>`)
	for _, match := range matches {
		heading_level, err := strconv.Atoi(match[1])
		if err != nil || heading_level < 1 || heading_level > 6 {
			continue
		}
		heading_stack[heading_level-1]++
		for i := heading_level; i < len(heading_stack); i++ {
			heading_stack[i] = 0
		}

		heading_number := strings.Builder{}
		for _, count := range heading_stack {
			if count == 0 {
				continue
			}
			if heading_number.Len() > 0 {
				heading_number.WriteString(".")
			}
			heading_number.WriteString(strconv.Itoa(count))
		}
		number := heading_number.String()
		indent := strings.Count(number, ".")
		toc_data.WriteString(`<br><span class="opennamu_TOC_list">` + strings.Repeat(`<span style="margin-left: 10px;"></span>`, indent) + `<a href="#s-` + number + `">` + number + `.</a> ` + match[2] + `</span>`)
		heading := `<h` + match[1] + ` id="s-` + number + `"><a href="#toc">` + number + `.</a> ` + match[2] + `</h` + match[1] + `>`
		data = strings.Replace(data, match[0], heading, 1)
	}
	toc_data.WriteString(`</div>`)
	return toc_data.String() + data
}

func View_user_skin_main(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	ban_data := tool.Get_user_ban(db, config.IP, "")
	if len(ban_data) > 0 && ban_data[0] == "true" {
		return tool.Get_error_page(db, config, "ban")
	}

	anonymous := tool.IP_or_user(config.IP)
	set_list := Get_main_skin_set_list(db)
	if values != nil {
		if config.Session != nil {
			config.Session.Delete("main_css_darkmode")
		}
		for field, choices := range set_list {
			value := values.Get(field)
			if value == "" && len(choices) > 0 {
				value = choices[0][0]
			}
			if !user_skin_choice(set_list, field, value) {
				continue
			}
			if field == "main_css_darkmode" {
				continue
			}
			if anonymous {
				if config.Session != nil {
					config.Session.Set(field, value)
				}
			} else {
				user_save(db, config.IP, field, value)
			}
		}
		if config.Session != nil {
			_ = config.Session.Save()
		}
		return tool.Get_redirect("/change/skin_set/main")
	}

	cookies := tool.Get_cookie_header(config.Cookies)
	current_value := func(field string) string {
		if field == "main_css_darkmode" {
			return cookies[field]
		}
		if anonymous {
			if config.Session == nil {
				return ""
			}
			value, _ := config.Session.Get(field).(string)
			return value
		}
		return user_value(db, config.IP, field)
	}
	field_data := func(field string, label string, heading string) string {
		choices := set_list[field]
		current := current_value(field)
		server_default := "default"
		tool.QueryRow_DB(db, "select data from other where name = ?", []any{&server_default}, field)

		server_label := ""
		options := strings.Builder{}
		for _, choice := range choices {
			if len(choice) < 2 {
				continue
			}
			if choice[0] == server_default {
				server_label = choice[1]
			}
			selected := ""
			if choice[0] == current {
				selected = ` selected="selected"`
			}
			options.WriteString(`<option value="` + tool.HTML_escape(choice[0]) + `"` + selected + `>` + tool.HTML_escape(choice[1]) + `</option>`)
		}

		data := strings.Builder{}
		if label != "" {
			data.WriteString(`<` + heading + `>` + tool.Get_language(db, label, true) + `</` + heading + `>`)
		}
		data.WriteString(tool.Get_language(db, "default", true) + " : " + tool.HTML_escape(server_label) + `<hr class="main_hr"><select name="` + tool.HTML_escape(field) + `">` + options.String() + `</select>`)
		return data.String()
	}

	body := strings.Builder{}
	body.WriteString(`<form method="post"><h2>` + tool.Get_language(db, "render", true) + `</h2>`)
	body.WriteString(field_data("main_css_strike", "strike", "h3"))
	body.WriteString(field_data("main_css_bold", "bold", "h3"))
	body.WriteString(`<h3>` + tool.Get_language(db, "category", true) + `</h3>`)
	body.WriteString(field_data("main_css_category_set", "position", "h4"))
	body.WriteString(field_data("main_css_category_change_title", "category_change_title", "h4"))
	body.WriteString(`<h3>` + tool.Get_language(db, "footnote", true) + ` (` + tool.Get_language(db, "beta", true) + `)</h3>`)
	body.WriteString(field_data("main_css_footnote_set", "footnote_render", "h4"))
	body.WriteString(field_data("main_css_footnote_number", "footnote_number", "h4"))
	body.WriteString(field_data("main_css_view_real_footnote_num", "footnote_real_num_view", "h4"))
	body.WriteString(field_data("main_css_include_link", "include_link", "h3"))
	body.WriteString(`<h3>` + tool.Get_language(db, "image", true) + ` (` + tool.Get_language(db, "beta", true) + `)</h3>`)
	body.WriteString(field_data("main_css_image_set", "", ""))
	body.WriteString(field_data("main_css_toc_set", "toc", "h3"))
	body.WriteString(field_data("main_css_exter_link", "exter_link", "h3"))
	body.WriteString(field_data("main_css_link_delimiter", "link_delimiter", "h3"))
	body.WriteString(field_data("main_css_darkmode", "force_darkmode", "h3"))
	body.WriteString(`<h3>` + tool.Get_language(db, "table", true) + `</h3>`)
	body.WriteString(field_data("main_css_table_scroll", "table_scroll", "h4"))
	body.WriteString(field_data("main_css_table_transparent", "table_transparent", "h4"))
	body.WriteString(field_data("main_css_list_view_change", "list_view_change", "h3"))
	body.WriteString(field_data("main_css_view_joke", "view_joke", "h3"))
	body.WriteString(field_data("main_css_math_scroll", "math_scroll", "h3"))
	body.WriteString(field_data("main_css_view_history", "view_history", "h3"))
	body.WriteString(field_data("main_css_font_size", "font_size", "h3"))
	body.WriteString(`<h2>` + tool.Get_language(db, "edit", true) + `</h2>`)
	body.WriteString(field_data("main_css_monaco", "monaco_editor", "h3"))
	body.WriteString(`<hr class="main_hr"><button type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`)

	menu := [][]any{
		{"change", tool.Get_language(db, "user_setting", true)},
		{"change/skin_set", tool.Get_language(db, "skin_set", true)},
		{"setting/skin_set", tool.Get_language(db, "main_skin_set_default", true)},
	}
	body_data := user_skin_main_render_simple_set(db, body.String())
	return tool.Get_template(db, config, tool.Get_language(db, "main_skin_set", true), body_data, []any{}, menu, map[string]string{})
}
