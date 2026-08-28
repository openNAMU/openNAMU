package route

import "opennamu/route/tool"

func View_thread_tool(config tool.Config, topic_num string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	title := ""
	sub := ""
	stop := ""
	agree := ""
	acl := ""
	view_acl := ""
	if !tool.QueryRow_DB(db, "select title, sub, stop, agree, acl from rd where code = ?", []any{&title, &sub, &stop, &agree, &acl}, topic_num) {
		return tool.Get_redirect("/")
	}
	tool.QueryRow_DB(db, "select set_data from topic_set where thread_code = ? and set_name = 'thread_view_acl'", []any{&view_acl}, topic_num)

	state := tool.Get_language(db, "topic_normal", true)
	if stop == "S" {
		state = tool.Get_language(db, "topic_stop", true)
	} else if stop == "O" {
		state = tool.Get_language(db, "topic_close", true)
	}
	if agree == "O" {
		state += " (" + tool.Get_language(db, "topic_agree", true) + ")"
	}

	acl_view_text := view_acl
	if acl_view_text == "" {
		acl_view_text = tool.Get_language(db, "normal", true)
	}
	acl_text := acl
	if acl_text == "" {
		acl_text = tool.Get_language(db, "normal", true)
	}
	data := `<h2>` + tool.Get_language(db, "tool", true) + `</h2><ul><li>` + tool.Get_language(db, "topic_state", true) + ` : ` + state + `</li><li>` + tool.Get_language(db, "topic_acl", true) + ` : ` + tool.HTML_escape(acl_text) + `</li><li>` + tool.Get_language(db, "topic_view_acl", true) + ` : ` + tool.HTML_escape(acl_view_text) + `</li></ul>`
	if tool.Check_permission(db, "thread_manage", config.IP) {
		data += `<h2>` + tool.Get_language(db, "admin_tool", true) + `</h2><ul><li><a href="/thread/` + tool.Url_parser(topic_num) + `/setting">` + tool.Get_language(db, "topic_setting", true) + `</a></li><li><a href="/thread/` + tool.Url_parser(topic_num) + `/acl">` + tool.Get_language(db, "topic_acl_setting", true) + `</a></li></ul>`
	}
	if tool.Check_permission(db, "thread_change", config.IP) {
		data += `<h2>` + tool.Get_language(db, "owner", true) + `</h2><ul><li><a href="/thread/` + tool.Url_parser(topic_num) + `/change">` + tool.Get_language(db, "topic_name_change", true) + `</a></li></ul>`
	}
	if tool.Check_permission(db, "thread_delete", config.IP) {
		data += `<h2>` + tool.Get_language(db, "owner", true) + `</h2><ul><li><a href="/thread/` + tool.Url_parser(topic_num) + `/delete">` + tool.Get_language(db, "topic_delete", true) + `</a></li></ul>`
	}

	return tool.Get_template(db, config, tool.Get_language(db, "topic_tool", true), data, []any{"(" + tool.HTML_escape(sub) + ")"}, [][]any{{"thread/" + tool.Url_parser(topic_num), tool.Get_language(db, "return", true)}}, map[string]string{})
}
