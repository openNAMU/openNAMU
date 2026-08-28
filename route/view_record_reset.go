package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_record_reset(config tool.Config, user_name string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if values == nil && !tool.Check_permission(db, "record_manage", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	if user_name == "" {
		return tool.Get_error_page(db, config, "error")
	}
	if values != nil {
		result := Api_record_reset_post(config, user_name)
		if result["response"] != "ok" {
			return tool.Get_error_page(db, config, "error")
		}
		return tool.Get_redirect("/record/" + tool.Url_parser(user_name))
	}
	return history_destructive_page(db, config, tool.Get_language(db, "record_reset", true), tool.Get_language(db, "reset", true), "record/"+tool.Url_parser(user_name))
}
