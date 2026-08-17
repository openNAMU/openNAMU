package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_history_delete(config tool.Config, doc_name string, rev string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	revision, ok := history_revision_value(rev)
	if !ok || doc_name == "" {
		return tool.Get_error_page(db, config, "error")
	}
	if values != nil {
		tool.Exec_DB(db, "delete from history where id = ? and title = ?", revision, doc_name)
		tool.Do_insert_auth_history(db, config.IP, "history_delete ("+doc_name+" r"+revision+")")
		return tool.Get_redirect("/history/" + tool.Url_parser(doc_name))
	}
	return history_destructive_page(db, config, tool.Get_language(db, "history_delete", true)+" (r"+revision+")", tool.Get_language(db, "delete", true), "history/"+tool.Url_parser(doc_name))
}
