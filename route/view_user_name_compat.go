package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_user_name_for(config tool.Config, target string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if target == "" {
		if values == nil && !user_auth(db, config) {
			return tool.Get_redirect("/login")
		}
		target = config.IP
	} else if values == nil && !tool.Check_permission(db, "user_manage", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	if values != nil {
		value := values.Get("new_user_name")
		if value == "" {
			value = values.Get("data")
		}
		api_data := Api_user_name_post(config, target, value)
		response, _ := api_data["response"].(string)
		if response == "require auth" {
			if target == config.IP {
				return tool.Get_redirect("/login")
			}
			return tool.Get_error_page(db, config, "auth")
		}
		if response != "ok" {
			return tool.Get_error_page(db, config, "user name error")
		}
		if target != config.IP {
			return tool.Get_redirect("/change/user_name/" + tool.Url_parser(target))
		}
		return tool.Get_redirect("/change/user_name")
	}

	current := user_value(db, target, "user_name")
	if current == "" {
		current = target
	}
	body := `<form method="post"><input name="new_user_name" value="` + tool.HTML_escape(current) + `"><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`
	return user_form_page(db, config, tool.Get_language(db, "change_user_name", true), body)
}
