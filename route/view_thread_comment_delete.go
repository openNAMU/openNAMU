package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_thread_comment_delete(config tool.Config, topic_num string, comment_num string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if values == nil && !tool.Check_permission(db, "thread_comment_delete", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	_, topic_exists := tool.Get_topic_data(db, topic_num, comment_num)
	if !topic_exists {
		return tool.Get_redirect("/thread/" + tool.Url_parser(topic_num))
	}
	if values != nil {
		api_data := Api_thread_comment_delete_post(config, topic_num, comment_num)
		response, _ := api_data["response"].(string)
		if response == "require auth" {
			return tool.Get_error_page(db, config, "auth")
		}
		if response == "not exist" {
			return tool.Get_redirect("/thread/" + tool.Url_parser(topic_num))
		}
		if response != "ok" {
			return tool.Get_error_page(db, config, "error")
		}
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
