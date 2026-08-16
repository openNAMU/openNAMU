package route

import (
	"net/url"
	"strings"

	"opennamu/route/tool"
)

func View_thread_comment_tool(config tool.Config, topic_num string, comment_num string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", topic_num, "topic_view", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	block := ""
	ip := ""
	date := ""
	if !tool.QueryRow_DB(
		db,
		"select block, ip, date from topic where code = ? and id = ?",
		[]any{&block, &ip, &date},
		topic_num,
		comment_num,
	) {
		return tool.Get_redirect("/thread/" + tool.Url_parser(topic_num))
	}

	data := `<h2>` + tool.Get_language(db, "state", true) + `</h2><ul><li>` + tool.Get_language(db, "writer", true) + ` : ` + tool.IP_parser(db, ip, config.IP) + `</li><li>` + tool.Get_language(db, "time", true) + ` : ` + tool.HTML_escape(date) + `</li></ul>`
	data += `<h2>` + tool.Get_language(db, "other_tool", true) + `</h2><ul><li><a href="/thread/` + tool.Url_parser(topic_num) + `/comment/` + tool.Url_parser(comment_num) + `/raw">` + tool.Get_language(db, "raw", true) + `</a></li></ul>`

	if tool.Check_acl(db, "", "", "toron_auth", config.IP) {
		data += `<h2>` + tool.Get_language(db, "admin_tool", true) + `</h2><ul><li><a href="/auth/ban/` + tool.Url_parser(ip) + `">` + tool.Get_language(db, "ban", true) + ` | ` + tool.Get_language(db, "release", true) + `</a></li><li><a href="/thread/` + tool.Url_parser(topic_num) + `/comment/` + tool.Url_parser(comment_num) + `/blind">` + tool.Get_language(db, "hide", true) + ` | ` + tool.Get_language(db, "hide_release", true) + `</a></li><li><a href="/thread/` + tool.Url_parser(topic_num) + `/comment/` + tool.Url_parser(comment_num) + `/notice">` + tool.Get_language(db, "pinned", true) + ` | ` + tool.Get_language(db, "pinned_release", true) + `</a></li><li><a href="/thread/` + tool.Url_parser(topic_num) + `/comment/` + tool.Url_parser(comment_num) + `/delete">` + tool.Get_language(db, "delete", true) + `</a></li></ul>`
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

func View_thread_comment_notice(config tool.Config, topic_num string, comment_num string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "toron_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	top := ""
	if !tool.QueryRow_DB(db, "select top from topic where code = ? and id = ?", []any{&top}, topic_num, comment_num) {
		return tool.Get_redirect("/thread/" + tool.Url_parser(topic_num))
	}
	if values != nil {
		if top == "O" {
			top = ""
		} else {
			top = "O"
		}
		tool.Exec_DB(db, "update topic set top = ? where code = ? and id = ?", top, topic_num, comment_num)
		tool.Exec_DB(db, "update rd set date = ? where code = ?", tool.Get_time(), topic_num)
		tool.Do_insert_auth_history(db, config.IP, "notice (code "+topic_num+"#"+comment_num+")")
		return tool.Get_redirect("/thread/" + tool.Url_parser(topic_num) + "#" + tool.Url_parser(comment_num))
	}

	action := "pinned"
	if top == "O" {
		action = "pinned_release"
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
	if tool.Check_acl(db, "", "", "toron_auth", config.IP) {
		data += `<h2>` + tool.Get_language(db, "admin_tool", true) + `</h2><ul><li><a href="/thread/` + tool.Url_parser(topic_num) + `/setting">` + tool.Get_language(db, "topic_setting", true) + `</a></li><li><a href="/thread/` + tool.Url_parser(topic_num) + `/acl">` + tool.Get_language(db, "topic_acl_setting", true) + `</a></li></ul>`
	}
	if tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		data += `<h2>` + tool.Get_language(db, "owner", true) + `</h2><ul><li><a href="/thread/` + tool.Url_parser(topic_num) + `/change">` + tool.Get_language(db, "topic_name_change", true) + `</a></li><li><a href="/thread/` + tool.Url_parser(topic_num) + `/delete">` + tool.Get_language(db, "topic_delete", true) + `</a></li></ul>`
	}

	return tool.Get_template(db, config, tool.Get_language(db, "topic_tool", true), data, []any{"(" + tool.HTML_escape(sub) + ")"}, [][]any{{"thread/" + tool.Url_parser(topic_num), tool.Get_language(db, "return", true)}}, map[string]string{})
}

func View_thread_setting(config tool.Config, topic_num string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "toron_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	name := ""
	stop := ""
	agree := ""
	if !tool.QueryRow_DB(db, "select title, stop, agree from rd where code = ?", []any{&name, &stop, &agree}, topic_num) {
		return tool.Get_redirect("/")
	}

	if values != nil {
		old_stop := stop
		old_agree := agree
		new_stop := values.Get("stop")
		if new_stop == "" {
			new_stop = values.Get("stop_d")
		}
		if new_stop != "" && new_stop != "S" && new_stop != "O" {
			new_stop = ""
		}
		new_agree := ""
		if values.Get("agree") != "" {
			new_agree = "O"
		}
		tool.Exec_DB(db, "update rd set stop = ?, agree = ? where code = ?", new_stop, new_agree, topic_num)
		if old_stop != new_stop {
			state_key := "topic_state_change_normal"
			if new_stop == "S" {
				state_key = "topic_state_change_stop"
			} else if new_stop == "O" {
				state_key = "topic_state_change_close"
			}
			thread_add(db, topic_num, thread_next_id(db, topic_num), tool.Get_language(db, state_key, true), config.IP, "1")
		}
		if old_agree != new_agree {
			state_key := "topic_state_change_disagree"
			if new_agree == "O" {
				state_key = "topic_state_change_agree"
			}
			thread_add(db, topic_num, thread_next_id(db, topic_num), tool.Get_language(db, state_key, true), config.IP, "1")
		}
		if why := strings.TrimSpace(values.Get("why")); why != "" {
			thread_add(db, topic_num, thread_next_id(db, topic_num), tool.Get_language(db, "why", true)+" : "+why, config.IP, "1")
		}
		tool.Exec_DB(db, "update rd set date = ? where code = ?", tool.Get_time(), topic_num)
		tool.Do_insert_auth_history(db, config.IP, "change_topic_set (code "+topic_num+")")
		return tool.Get_redirect("/thread/" + tool.Url_parser(topic_num))
	}

	options := `<option value="">` + tool.Get_language(db, "topic_normal", true) + `</option><option value="S"` + thread_selected(stop, "S") + `>` + tool.Get_language(db, "topic_stop", true) + `</option><option value="O"` + thread_selected(stop, "O") + `>` + tool.Get_language(db, "topic_close", true) + `</option>`
	data := `<form method="post"><h2>` + tool.Get_language(db, "topic_progress", true) + `</h2><select name="stop">` + options + `</select><hr class="main_hr"><label><input type="checkbox" name="agree" value="O"` + thread_checked(agree == "O") + `> ` + tool.Get_language(db, "topic_change_agree", true) + `</label><h2>` + tool.Get_language(db, "why", true) + `</h2><input name="why"><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`
	return tool.Get_template(db, config, tool.Get_language(db, "topic_setting", true), data, []any{"(" + tool.HTML_escape(name) + ")"}, [][]any{{"thread/" + tool.Url_parser(topic_num) + "/tool", tool.Get_language(db, "return", true)}}, map[string]string{})
}

func thread_selected(value string, expected string) string {
	if value == expected {
		return ` selected`
	}
	return ""
}

func thread_checked(checked bool) string {
	if checked {
		return ` checked`
	}
	return ""
}

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
		if !tool.Arr_in_str(tool.List_acl("normal"), acl) || !tool.Arr_in_str(tool.List_acl("normal"), acl_view) {
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
	data := `<form method="post"><a href="/acl/TEST#exp">(` + tool.Get_language(db, "reference", true) + `)</a><h2>` + tool.Get_language(db, "thread_acl", true) + `</h2>` + bbs_set_select(db, "acl", thread_acl, tool.List_acl("normal")) + `<h2>` + tool.Get_language(db, "view_acl", true) + ` (` + tool.Get_language(db, "beta", true) + `)</h2>` + bbs_set_select(db, "acl_view", acl_view, tool.List_acl("normal")) + `<hr class="main_hr"><button type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`
	return tool.Get_template(db, config, tool.Get_language(db, "topic_acl_setting", true), data, []any{"(" + tool.HTML_escape(title) + ")"}, [][]any{{"thread/" + tool.Url_parser(topic_num) + "/tool", tool.Get_language(db, "return", true)}}, map[string]string{})
}

func View_thread_delete(config tool.Config, topic_num string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	title := ""
	if !tool.QueryRow_DB(db, "select title from rd where code = ?", []any{&title}, topic_num) {
		return tool.Get_redirect("/")
	}
	if values != nil {
		tool.Do_insert_auth_history(db, config.IP, "delete_topic (code "+topic_num+")")
		tool.Exec_DB(db, "delete from topic where code = ?", topic_num)
		tool.Exec_DB(db, "delete from topic_set where thread_code = ?", topic_num)
		tool.Exec_DB(db, "delete from rd where code = ?", topic_num)
		return tool.Get_redirect("/topic/" + tool.Url_parser(title))
	}
	data := `<form method="post"><p>` + tool.Get_language(db, "delete_warning", true) + `</p><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "delete", true) + `</button></form>`
	return tool.Get_template(db, config, tool.Get_language(db, "topic_delete", true), data, []any{"(" + tool.HTML_escape(title) + ")"}, [][]any{{"thread/" + tool.Url_parser(topic_num), tool.Get_language(db, "return", true)}}, map[string]string{})
}

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
		tool.Exec_DB(db, "update rd set title = ?, sub = ?, date = ? where code = ?", new_title, new_sub, tool.Get_time(), topic_num)
		thread_add(db, topic_num, thread_next_id(db, topic_num), tool.Get_language(db, "topic_name_change", true)+" : "+sub+" ("+title+") → "+new_sub+" ("+new_title+")", config.IP, "1")
		tool.Do_insert_auth_history(db, config.IP, "change_topic_name (code "+topic_num+")")
		return tool.Get_redirect("/thread/" + tool.Url_parser(topic_num))
	}
	data := `<form method="post"><input name="title" value="` + tool.HTML_escape(title) + `"><hr class="main_hr"><input name="sub" value="` + tool.HTML_escape(sub) + `"><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`
	return tool.Get_template(db, config, tool.Get_language(db, "topic_name_change", true), data, []any{"(" + tool.HTML_escape(title) + ")"}, [][]any{{"thread/" + tool.Url_parser(topic_num) + "/tool", tool.Get_language(db, "return", true)}}, map[string]string{})
}
