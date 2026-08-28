package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_history_delete(config tool.Config, doc_name string, rev string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if values == nil && !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	revision, ok := history_revision_value(rev)
	if !ok || doc_name == "" {
		return tool.Get_error_page(db, config, "error")
	}
	if values != nil {
		api_data := Api_history_delete_post(config, doc_name, revision)
		response, _ := api_data["response"].(string)
		if response == "require auth" {
			return tool.Get_error_page(db, config, "auth")
		}
		if response != "ok" {
			return tool.Get_error_page(db, config, "error")
		}
		return tool.Get_redirect("/history/" + tool.Url_parser(doc_name))
	}
	return history_destructive_page(db, config, tool.Get_language(db, "history_delete", true)+" (r"+revision+")", tool.Get_language(db, "delete", true), "history/"+tool.Url_parser(doc_name))
}
