package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_thread_change(config tool.Config, topic_num string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if values == nil && !tool.Check_permission(db, "thread_change", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	title := ""
	sub := ""
	if !tool.QueryRow_DB(db, "select title, sub from rd where code = ?", []any{&title, &sub}, topic_num) {
		return tool.Get_redirect("/")
	}
	if values != nil {
		api_data := Api_thread_change_post(config, topic_num, values.Get("title"), values.Get("sub"))
		response, _ := api_data["response"].(string)
		if response == "require auth" {
			return tool.Get_error_page(db, config, "auth")
		}
		if response == "not exist" {
			return tool.Get_redirect("/")
		}
		if response != "ok" {
			error_name, _ := api_data["data"].(string)
			if error_name == "" {
				error_name = "error"
			}
			return tool.Get_error_page(db, config, error_name)
		}
		return tool.Get_redirect("/thread/" + tool.Url_parser(topic_num))
	}

	data := `<form method="post"><input name="title" value="` + tool.HTML_escape(title) + `"><hr class="main_hr"><input name="sub" value="` + tool.HTML_escape(sub) + `"><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`
	return tool.Get_template(db, config, tool.Get_language(db, "topic_name_change", true), data, []any{"(" + tool.HTML_escape(title) + ")"}, [][]any{{"thread/" + tool.Url_parser(topic_num) + "/tool", tool.Get_language(db, "return", true)}}, map[string]string{})
}
