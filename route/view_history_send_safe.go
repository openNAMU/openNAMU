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
		response, _ := api_data["response"].(string)
		if response == "require auth" {
			return tool.Get_error_page(db, config, "auth")
		}
		if response != "ok" {
			return tool.Get_error_page(db, config, "error")
		}
		return tool.Get_redirect("/history/" + tool.Url_parser(doc_name))
	}
	send := tool.Get_history_send(db, doc_name, rev)
	body := `<form method="post"><input name="send" value="` + tool.HTML_escape(send) + `"><button type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`
	return document_safe_page(db, config, "history send", body)
}
