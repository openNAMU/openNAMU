package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_user_top_menu(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if values == nil && !user_auth(db, config) {
		return tool.Get_redirect("/user")
	}
	if values != nil {
		content := values.Get("content")
		if !values.Has("content") {
			content = values.Get("data")
		}
		api_data := Api_user_top_menu_post(config, content)
		response, _ := api_data["response"].(string)
		if response == "require auth" {
			return tool.Get_redirect("/user")
		}
		if response != "ok" {
			return tool.Get_error_page(db, config, "error")
		}
		return tool.Get_redirect("/change/top_menu")
	}
	content := user_value(db, config.IP, "top_menu")
	body := `<span>EX)<br>ONTS<br>https://2du.pythonanywhere.com/<br>FrontPage<br>/w/FrontPage</span><hr class="main_hr">` + tool.Get_language(db, "not_support_skin_warning", true) + `<hr class="main_hr"><form method="post"><textarea class="opennamu_textarea_500" placeholder="` + tool.Get_language(db, "enter_top_menu_setting", true) + `" name="content" id="content">` + tool.HTML_escape(content) + `</textarea><hr class="main_hr"><button id="opennamu_save_button" type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`
	return user_form_page(db, config, tool.Get_language(db, "user_added_menu", true), body)
}
