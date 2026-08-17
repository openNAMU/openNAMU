package route

import (
	"net/url"
	"opennamu/route/tool"
	"strings"
)

func View_history_add_safe(config tool.Config, doc_name string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	if values != nil {
		content := strings.ReplaceAll(values.Get("content"), "\r", "")
		Do_add_history(db, doc_name, content, tool.Get_time(), "Add:"+values.Get("get_ip"), values.Get("send"), "+"+tool.Get_edit_length_diff("", content), "", "add")
		return tool.Get_redirect("/history/" + tool.Url_parser(doc_name))
	}
	body := `<form method="post"><input name="send"><input name="get_ip"><textarea name="content" class="opennamu_textarea_500"></textarea><button type="submit">` + tool.Get_language(db, "add", true) + `</button></form>`
	return document_safe_page(db, config, "history add", body)
}
