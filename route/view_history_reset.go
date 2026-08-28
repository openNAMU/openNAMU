package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_history_reset(config tool.Config, doc_name string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if values == nil && !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	if doc_name == "" {
		return tool.Get_error_page(db, config, "error")
	}
	if values != nil {
		api_data := Api_history_reset_post(config, doc_name)
		response, _ := api_data["response"].(string)
		if response == "require auth" {
			return tool.Get_error_page(db, config, "auth")
		}
		if response != "ok" {
			return tool.Get_error_page(db, config, "error")
		}
		return tool.Get_redirect("/history/" + tool.Url_parser(doc_name))
	}
	return history_destructive_page(db, config, tool.Get_language(db, "history_reset", true), tool.Get_language(db, "reset", true), "history/"+tool.Url_parser(doc_name))
}
