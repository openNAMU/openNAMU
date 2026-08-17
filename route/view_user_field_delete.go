package route

import (
	"net/url"
	"opennamu/route/tool"
)

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
