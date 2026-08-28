package route

import (
	"net/url"
	"opennamu/route/tool"
)

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
			api_data := Api_user_head_skin_post(config, storage_name, content)
			response, _ := api_data["response"].(string)
			if response != "ok" {
				return tool.Get_error_page(db, config, "error")
			}
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
	return user_form_page(db, config, tool.Get_language(db, "user_head", false)+sub_name, body)
}
