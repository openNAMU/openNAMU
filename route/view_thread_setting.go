package route

import (
	"net/url"
	"strings"

	"opennamu/route/tool"
)

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
