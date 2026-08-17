package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_thread_comment_blind(config tool.Config, topic_num string, comment_num string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "toron_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	block := ""
	if !tool.QueryRow_DB(db, "select block from topic where code = ? and id = ?", []any{&block}, topic_num, comment_num) {
		return tool.Get_redirect("/thread/" + tool.Url_parser(topic_num))
	}
	if values != nil {
		if block == "O" {
			block = ""
		} else {
			block = "O"
		}
		tool.Exec_DB(db, "update topic set block = ? where code = ? and id = ?", block, topic_num, comment_num)
		tool.Exec_DB(db, "update rd set date = ? where code = ?", tool.Get_time(), topic_num)
		tool.Do_insert_auth_history(db, config.IP, "blind (code "+topic_num+"#"+comment_num+")")
		return tool.Get_redirect("/thread/" + tool.Url_parser(topic_num) + "#" + tool.Url_parser(comment_num))
	}

	action := "hide"
	if block == "O" {
		action = "hide_release"
	}
	data := `<form method="post"><button type="submit">` + tool.Get_language(db, action, true) + `</button></form>`
	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, "discussion_tool", true),
		data,
		[]any{"(#" + tool.HTML_escape(comment_num) + ")"},
		[][]any{{"thread/" + tool.Url_parser(topic_num) + "/comment/" + tool.Url_parser(comment_num) + "/tool", tool.Get_language(db, "return", true)}},
		map[string]string{},
	)
}
