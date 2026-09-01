package route

import "opennamu/route/tool"

func View_thread_comment_tool(config tool.Config, topic_num string, comment_num string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", topic_num, "topic_view", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	topic_data, topic_exists := tool.Get_topic_data(db, topic_num, comment_num)
	if !topic_exists {
		return tool.Get_redirect("/thread/" + tool.Url_parser(topic_num))
	}
	ip := topic_data["ip"]
	date := topic_data["date"]

	data := `<h2>` + tool.Get_language(db, "state", true) + `</h2><ul><li>` + tool.Get_language(db, "writer", true) + ` : ` + tool.IP_parser(db, ip, config.IP) + `</li><li>` + tool.Get_language(db, "time", true) + ` : ` + tool.HTML_escape(date) + `</li></ul>`
	data += `<h2>` + tool.Get_language(db, "other_tool", true) + `</h2><ul><li><a href="/thread/` + tool.Url_parser(topic_num) + `/comment/` + tool.Url_parser(comment_num) + `/raw">` + tool.Get_language(db, "raw", true) + `</a></li></ul>`

	if tool.Check_permission(db, "give", config.IP) {
		data += `<h2>` + tool.Get_language(db, "admin_tool", true) + `</h2><ul><li><a href="/auth/give/` + tool.Url_parser(ip) + `">` + tool.Get_language(db, "ban", true) + ` | ` + tool.Get_language(db, "release", true) + `</a></li></ul>`
	}
	if tool.Check_permission(db, "thread_comment_manage", config.IP) {
		data += `<h2>` + tool.Get_language(db, "admin_tool", true) + `</h2><ul><li><a href="/thread/` + tool.Url_parser(topic_num) + `/comment/` + tool.Url_parser(comment_num) + `/blind">` + tool.Get_language(db, "hide", true) + ` | ` + tool.Get_language(db, "hide_release", true) + `</a></li><li><a href="/thread/` + tool.Url_parser(topic_num) + `/comment/` + tool.Url_parser(comment_num) + `/notice">` + tool.Get_language(db, "pinned", true) + ` | ` + tool.Get_language(db, "pinned_release", true) + `</a></li></ul>`
	}
	if tool.Check_permission(db, "thread_comment_delete", config.IP) {
		data += `<h2>` + tool.Get_language(db, "admin_tool", true) + `</h2><ul><li><a href="/thread/` + tool.Url_parser(topic_num) + `/comment/` + tool.Url_parser(comment_num) + `/delete">` + tool.Get_language(db, "delete", true) + `</a></li></ul>`
	}

	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, "discussion_tool", true),
		data,
		[]any{"(#" + tool.HTML_escape(comment_num) + ")"},
		[][]any{{"thread/" + tool.Url_parser(topic_num) + "#" + tool.Url_parser(comment_num), tool.Get_language(db, "return", true)}},
		map[string]string{},
	)
}
