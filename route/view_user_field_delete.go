package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_user_field_delete(config tool.Config, field string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if values == nil && !user_auth(db, config) {
		return tool.Get_redirect("/login")
	}
	if values != nil {
		api_data := Api_user_field_delete_post(config, field)
		response, _ := api_data["response"].(string)
		if response == "require auth" {
			return tool.Get_redirect("/login")
		}
		if response != "ok" {
			return tool.Get_error_page(db, config, "error")
		}
		return tool.Get_redirect("/change")
	}
	title := tool.Get_language(db, field+"_delete", true)
	body := `<form method="post"><span>` + tool.Get_language(db, "delete_warning", true) + `</span><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "delete", true) + `</button></form>`
	return user_form_page(db, config, title, body)
}
