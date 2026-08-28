package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_thread_setting(config tool.Config, topic_num string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if values == nil && !tool.Check_acl(db, "", "", "toron_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	name := ""
	stop := ""
	agree := ""
	if !tool.QueryRow_DB(db, "select title, stop, agree from rd where code = ?", []any{&name, &stop, &agree}, topic_num) {
		return tool.Get_redirect("/")
	}

	if values != nil {
		stop_value := values.Get("stop")
		if stop_value == "" {
			stop_value = values.Get("stop_d")
		}
		api_data := Api_thread_setting_post(config, topic_num, stop_value, values.Get("agree"), values.Get("why"))
		response, _ := api_data["response"].(string)
		if response == "require auth" {
			return tool.Get_error_page(db, config, "auth")
		}
		if response == "not exist" {
			return tool.Get_redirect("/")
		}
		if response != "ok" {
			return tool.Get_error_page(db, config, "error")
		}
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
