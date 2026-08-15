package route

import (
	"database/sql"
	"net/url"
	"strings"

	"opennamu/route/tool"
)

func document_safe_page(db *sql.DB, config tool.Config, title string, body string) string {
	return tool.Get_template(db, config, title, body, []any{}, [][]any{{"other", tool.Get_language(db, "return", true)}}, map[string]string{})
}

func View_history_hidden_safe(config tool.Config, doc_name string, rev string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !tool.Check_acl(db, "", "", "hidel_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	hide := ""
	tool.QueryRow_DB(db, "select hide from history where title = ? and id = ?", []any{&hide}, doc_name, rev)
	if hide == "" {
		hide = "O"
	} else {
		hide = ""
	}
	tool.Exec_DB(db, "update history set hide = ? where title = ? and id = ?", hide, doc_name, rev)
	return tool.Get_redirect("/history/" + tool.Url_parser(doc_name))
}

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
