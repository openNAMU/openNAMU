package route

import (
	"net/url"
	"opennamu/route/tool"
	"strings"
)

func View_user_skin_main(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	auth_name := tool.Get_user_auth(db, config.IP)
	if tool.Auth_group_name_ban(auth_name) {
		return tool.Get_error_page(db, config, "ban")
	}

	anonymous := tool.IP_or_user(config.IP)
	set_list := Get_main_skin_set_list(db)
	if values != nil {
		user_set_list := map[string]string{}
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
				user_set_list[field] = value
			}
		}
		if !anonymous {
			api_data := Api_user_setting_skin_set_main_post(config, user_set_list)
			response, _ := api_data["response"].(string)
			if response != "ok" {
				return tool.Get_error_page(db, config, "error")
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
		server_default := tool.Get_setting_value(db, field, "", "default")

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
	body.WriteString(field_data("main_css_video_set", "video", "h3"))
	body.WriteString(field_data("main_css_toc_set", "toc", "h3"))
	body.WriteString(field_data("main_css_exter_link", "exter_link", "h3"))
	body.WriteString(field_data("main_css_link_delimiter", "link_delimiter", "h3"))
	body.WriteString(field_data("main_css_darkmode", "force_darkmode", "h3"))
	body.WriteString(`<h3>` + tool.Get_language(db, "table", true) + `</h3>`)
	body.WriteString(field_data("main_css_table_scroll", "table_scroll", "h4"))
	body.WriteString(field_data("main_css_table_transparent", "table_transparent", "h4"))
	body.WriteString(field_data("main_css_table_auto_color", "table_auto_color", "h4"))
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
