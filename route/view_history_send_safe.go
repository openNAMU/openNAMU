package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_history_send_safe(config tool.Config, doc_name string, rev string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	if values != nil {
		tool.Exec_DB(db, "update history set send = ? where title = ? and id = ?", values.Get("send"), doc_name, rev)
		return tool.Get_redirect("/history/" + tool.Url_parser(doc_name))
	}
	send := ""
	tool.QueryRow_DB(db, "select send from history where title = ? and id = ?", []any{&send}, doc_name, rev)
	body := `<form method="post"><input name="send" value="` + tool.HTML_escape(send) + `"><button type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`
	return document_safe_page(db, config, "history send", body)
}
