package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_thread_acl(config tool.Config, topic_num string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "toron_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	title := ""
	thread_acl := ""
	if !tool.QueryRow_DB(db, "select sub, acl from rd where code = ?", []any{&title, &thread_acl}, topic_num) {
		return tool.Get_redirect("/")
	}
	if values != nil {
		acl := values.Get("acl")
		acl_view := values.Get("acl_view")
		if !acl_value_valid(db, acl) || !acl_value_valid(db, acl_view) {
			return tool.Get_error_page(db, config, "error")
		}
		tool.Exec_DB(db, "update rd set acl = ?, date = ? where code = ?", acl, tool.Get_time(), topic_num)
		tool.Exec_DB(db, "delete from topic_set where thread_code = ? and set_name = 'thread_view_acl'", topic_num)
		tool.Exec_DB(db, "insert into topic_set (thread_code, set_name, set_id, set_data) values (?, 'thread_view_acl', '1', ?)", topic_num, acl_view)
		thread_add(db, topic_num, thread_next_id(db, topic_num), tool.Get_language(db, "acl_thread_change", true)+" : "+acl, config.IP, "1")
		tool.Do_insert_auth_history(db, config.IP, "change_topic_acl (code "+topic_num+")")
		return tool.Get_redirect("/thread/" + tool.Url_parser(topic_num))
	}

	acl_view := ""
	tool.QueryRow_DB(db, "select set_data from topic_set where thread_code = ? and set_name = 'thread_view_acl'", []any{&acl_view}, topic_num)
	acl_values := acl_value_list(db, thread_acl)
	acl_view_values := acl_value_list(db, acl_view)
	data := `<form method="post"><a href="/acl/TEST#exp">(` + tool.Get_language(db, "reference", true) + `)</a><h2>` + tool.Get_language(db, "thread_acl", true) + `</h2>` + bbs_set_select(db, "acl", thread_acl, acl_values) + `<h2>` + tool.Get_language(db, "view_acl", true) + ` (` + tool.Get_language(db, "beta", true) + `)</h2>` + bbs_set_select(db, "acl_view", acl_view, acl_view_values) + `<hr class="main_hr"><button type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`
	return tool.Get_template(db, config, tool.Get_language(db, "topic_acl_setting", true), data, []any{"(" + tool.HTML_escape(title) + ")"}, [][]any{{"thread/" + tool.Url_parser(topic_num) + "/tool", tool.Get_language(db, "return", true)}}, map[string]string{})
}
