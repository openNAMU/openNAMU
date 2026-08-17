package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_history_hidden_safe(config tool.Config, doc_name string, rev string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !tool.Check_acl(db, "", "", "hidel_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	hide := ""
	tool.QueryRow_DB(db, "select hide from history where title = ? and id = ?", []any{&hide}, doc_name, rev)
	if values != nil {
		if hide == "" {
			hide = "O"
		} else {
			hide = ""
		}
		tool.Exec_DB(db, "update history set hide = ? where title = ? and id = ?", hide, doc_name, rev)
		return tool.Get_redirect("/history/" + tool.Url_parser(doc_name))
	}

	action := "hide"
	if hide == "O" {
		action = "hide_release"
	}
	body := `<form method="post"><button type="submit">` + tool.Get_language(db, action, true) + `</button></form>`
	return document_safe_page(db, config, tool.Get_language(db, action, true), body)
}
