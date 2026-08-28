package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_thread_acl(config tool.Config, topic_num string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if values == nil && !tool.Check_permission(db, "thread_acl", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	title := ""
	thread_acl := ""
	if !tool.QueryRow_DB(db, "select sub, acl from rd where code = ?", []any{&title, &thread_acl}, topic_num) {
		return tool.Get_redirect("/")
	}
	if values != nil {
		api_data := Api_thread_acl_post(config, topic_num, values.Get("acl"), values.Get("acl_view"))
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

	acl_view := ""
	tool.QueryRow_DB(db, "select set_data from topic_set where thread_code = ? and set_name = 'thread_view_acl'", []any{&acl_view}, topic_num)
	acl_values := acl_value_list(db, thread_acl)
	acl_view_values := acl_value_list(db, acl_view)
	data := `<form method="post"><a href="/acl/TEST#exp">(` + tool.Get_language(db, "reference", true) + `)</a><h2>` + tool.Get_language(db, "thread_acl", true) + `</h2>` + bbs_set_select(db, "acl", thread_acl, acl_values) + `<h2>` + tool.Get_language(db, "view_acl", true) + ` (` + tool.Get_language(db, "beta", true) + `)</h2>` + bbs_set_select(db, "acl_view", acl_view, acl_view_values) + `<hr class="main_hr"><button type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`
	return tool.Get_template(db, config, tool.Get_language(db, "topic_acl_setting", true), data, []any{"(" + tool.HTML_escape(title) + ")"}, [][]any{{"thread/" + tool.Url_parser(topic_num) + "/tool", tool.Get_language(db, "return", true)}}, map[string]string{})
}
