package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_thread_change(config tool.Config, topic_num string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	title := ""
	sub := ""
	if !tool.QueryRow_DB(db, "select title, sub from rd where code = ?", []any{&title, &sub}, topic_num) {
		return tool.Get_redirect("/")
	}
	if values != nil {
		new_title := values.Get("title")
		new_sub := values.Get("sub")
		if new_title == "" {
			new_title = title
		}
		if new_sub == "" {
			new_sub = sub
		}
		if !tool.Do_title_length_check(db, new_title, "document") {
			return tool.Get_error_page(db, config, "title length")
		}
		if !tool.Do_title_length_check(db, new_sub, "topic") {
			return tool.Get_error_page(db, config, "topic title length")
		}
		tool.Exec_DB(db, "update rd set title = ?, sub = ?, date = ? where code = ?", new_title, new_sub, tool.Get_time(), topic_num)
		thread_add(db, topic_num, thread_next_id(db, topic_num), tool.Get_language(db, "topic_name_change", true)+" : "+sub+" ("+title+") → "+new_sub+" ("+new_title+")", config.IP, "1")
		tool.Do_insert_auth_history(db, config.IP, "change_topic_name (code "+topic_num+")")
		return tool.Get_redirect("/thread/" + tool.Url_parser(topic_num))
	}
	data := `<form method="post"><input name="title" value="` + tool.HTML_escape(title) + `"><hr class="main_hr"><input name="sub" value="` + tool.HTML_escape(sub) + `"><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`
	return tool.Get_template(db, config, tool.Get_language(db, "topic_name_change", true), data, []any{"(" + tool.HTML_escape(title) + ")"}, [][]any{{"thread/" + tool.Url_parser(topic_num) + "/tool", tool.Get_language(db, "return", true)}}, map[string]string{})
}
