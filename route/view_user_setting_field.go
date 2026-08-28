package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_user_setting_field(config tool.Config, field string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if values == nil && !user_auth(db, config) {
		return tool.Get_redirect("/user")
	}
	if values != nil {
		api_data := Api_user_setting_field_post(config, field, values.Get("data"))
		response, _ := api_data["response"].(string)
		if response == "require auth" {
			return tool.Get_redirect("/user")
		}
		if response != "ok" {
			return tool.Get_error_page(db, config, "user name error")
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
