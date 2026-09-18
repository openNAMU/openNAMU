package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_history_send_safe(config tool.Config, doc_name string, rev string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if values == nil && !tool.Check_permission(db, "history_manage", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	if values != nil {
		api_data := Api_history_send_post(config, doc_name, rev, values.Get("send"))
		return tool.Api_post_redirect(db, config, api_data, "/history/"+tool.Url_parser(doc_name))
	}
	send := tool.Get_history_send(db, doc_name, rev)
	body := `<form method="post"><div><input name="send" value="` + tool.HTML_escape(send) + `"></div><hr class="main_hr"><div><button type="submit">` + tool.Get_language(db, "save", true) + `</button></div></form>`
	return Document_safe_page(db, config, tool.Get_language(db, "history", true)+" "+tool.Get_language(db, "send", true), body)
}
