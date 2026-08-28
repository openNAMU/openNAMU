package route

import "opennamu/route/tool"

func View_thread_raw(config tool.Config, topic_num string, comment_num string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", topic_num, "topic_view", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	title := ""
	if !tool.QueryRow_DB(db, "select sub from rd where code = ?", []any{&title}, topic_num) {
		return tool.Get_redirect("/")
	}

	data := ""
	block := ""
	if !tool.QueryRow_DB(
		db,
		"select data, block from topic where code = ? and id = ?",
		[]any{&data, &block},
		topic_num,
		comment_num,
	) {
		return tool.Get_redirect("/thread/" + tool.Url_parser(topic_num))
	}
	if block == "O" && !tool.Check_permission(db, "hidel", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	data_html := `<div id="opennamu_preview_area"><textarea readonly id="opennamu_edit_textarea" class="opennamu_textarea_500 __ON_TEXTAREA__">` + tool.HTML_escape(data) + `</textarea></div>`
	return tool.Get_template(
		db,
		config,
		title,
		data_html,
		[]any{"(" + tool.Get_language(db, "discussion_raw", true) + ") (#" + tool.HTML_escape(comment_num) + ")"},
		[][]any{
			{"thread/" + tool.Url_parser(topic_num) + "#" + tool.Url_parser(comment_num), tool.Get_language(db, "discussion", true)},
			{"thread/" + tool.Url_parser(topic_num) + "/comment/" + tool.Url_parser(comment_num) + "/tool", tool.Get_language(db, "return", true)},
		},
		map[string]string{},
	)
}
