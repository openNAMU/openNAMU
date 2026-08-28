package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_user_head_reset(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	skin_name := tool.Get_use_skin_name_session(db, config.IP, config.Session)
	if values != nil {
		if user_auth(db, config) {
			api_data := Api_user_head_reset_post(config, skin_name)
			response, _ := api_data["response"].(string)
			if response != "ok" {
				return tool.Get_error_page(db, config, "error")
			}
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
