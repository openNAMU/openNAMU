package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_thread_delete(config tool.Config, topic_num string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if values == nil && !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	title := ""
	if !tool.QueryRow_DB(db, "select title from rd where code = ?", []any{&title}, topic_num) {
		return tool.Get_redirect("/")
	}
	if values != nil {
		api_data := Api_thread_delete_post(config, topic_num)
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
		delete_title, _ := api_data["data"].(string)
		if delete_title == "" {
			delete_title = title
		}
		return tool.Get_redirect("/topic/" + tool.Url_parser(delete_title))
	}

	data := `<form method="post"><p>` + tool.Get_language(db, "delete_warning", true) + `</p><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "delete", true) + `</button></form>`
	return tool.Get_template(db, config, tool.Get_language(db, "topic_delete", true), data, []any{"(" + tool.HTML_escape(title) + ")"}, [][]any{{"thread/" + tool.Url_parser(topic_num), tool.Get_language(db, "return", true)}}, map[string]string{})
}
