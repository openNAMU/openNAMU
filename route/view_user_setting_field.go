package route

import (
	"net/url"
	"opennamu/route/tool"
)

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
