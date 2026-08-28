package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_history_add_safe(config tool.Config, doc_name string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if values == nil && !tool.Check_permission(db, "history_manage", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	if values != nil {
		api_data := Api_history_add_post(config, doc_name, values.Get("content"), values.Get("get_ip"), values.Get("send"))
		response, _ := api_data["response"].(string)
		if response == "require auth" {
			return tool.Get_error_page(db, config, "auth")
		}
		if response != "ok" {
			return tool.Get_error_page(db, config, "error")
		}
		return tool.Get_redirect("/history/" + tool.Url_parser(doc_name))
	}
	body := `<form method="post"><input name="send"><input name="get_ip"><textarea name="content" class="opennamu_textarea_500"></textarea><button type="submit">` + tool.Get_language(db, "add", true) + `</button></form>`
	return document_safe_page(db, config, "history add", body)
}
