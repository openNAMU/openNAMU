package route

import (
	"net/url"
	"opennamu/route/tool"
)

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
			api_data := Api_user_skin_post(config, skin)
			response, _ := api_data["response"].(string)
			if response != "ok" {
				return tool.Get_error_page(db, config, "error")
			}
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
