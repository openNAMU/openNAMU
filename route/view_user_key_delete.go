package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_user_key_delete(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !user_auth(db, config) {
		return tool.Get_redirect("/login")
	}
	if values != nil {
		user_delete(db, config.IP, "random_key")
		return tool.Get_redirect("/change")
	}
	body := `<form method="post"><span>` + tool.Get_language(db, "delete_warning", true) + `</span><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "delete", true) + `</button></form>`
	return user_form_page(db, config, tool.Get_language(db, "key_delete", true), body)
}
