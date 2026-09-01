package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_history_hidden_safe(config tool.Config, doc_name string, rev string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if values == nil && !tool.Check_permission(db, "hidel", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	hide := tool.Get_history_hide(db, doc_name, rev)
	if values != nil {
		api_data := Api_history_hidden_post(config, doc_name, rev)
		response, _ := api_data["response"].(string)
		if response == "require auth" {
			return tool.Get_error_page(db, config, "auth")
		}
		if response != "ok" {
			return tool.Get_error_page(db, config, "error")
		}
		return tool.Get_redirect("/history/" + tool.Url_parser(doc_name))
	}

	action := "hide"
	if hide == "O" {
		action = "hide_release"
	}
	body := `<form method="post"><button type="submit">` + tool.Get_language(db, action, true) + `</button></form>`
	return document_safe_page(db, config, tool.Get_language(db, action, true), body)
}
