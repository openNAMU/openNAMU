package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_thread_comment_delete(config tool.Config, topic_num string, comment_num string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	exists_id := ""
	if !tool.QueryRow_DB(db, "select id from topic where code = ? and id = ?", []any{&exists_id}, topic_num, comment_num) {
		return tool.Get_redirect("/thread/" + tool.Url_parser(topic_num))
	}
	if values != nil {
		tool.Exec_DB(db, "delete from topic where code = ? and id = ?", topic_num, comment_num)
		tool.Exec_DB(db, "update rd set date = ? where code = ?", tool.Get_time(), topic_num)
		tool.Do_insert_auth_history(db, config.IP, "delete_topic_comment (code "+topic_num+"#"+comment_num+")")
		return tool.Get_redirect("/thread/" + tool.Url_parser(topic_num))
	}

	data := `<hr class="main_hr"><form method="post"><button type="submit">` + tool.Get_language(db, "delete", true) + `</button></form>`
	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, "topic_delete", true),
		data,
		[]any{"(#" + tool.HTML_escape(comment_num) + ")"},
		[][]any{{"thread/" + tool.Url_parser(topic_num) + "/comment/" + tool.Url_parser(comment_num) + "/tool", tool.Get_language(db, "return", true)}},
		map[string]string{},
	)
}
