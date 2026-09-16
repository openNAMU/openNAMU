package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_history_delete(config tool.Config, doc_name string, rev string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if values == nil && !tool.Check_permission(db, "history_manage", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	revision, ok := history_revision_value(rev)
	if !ok || doc_name == "" {
		return tool.Get_error_page(db, config, "error")
	}
	if values != nil {
		api_data := Api_history_delete_post(config, doc_name, revision)
		return tool.Api_post_redirect(db, config, api_data, "/history/"+tool.Url_parser(doc_name))
	}
	return history_destructive_page(db, config, tool.Get_language(db, "history_delete", true)+" (r"+revision+")", tool.Get_language(db, "delete", true), "history/"+tool.Url_parser(doc_name))
}
