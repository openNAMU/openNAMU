package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_thread_change(config tool.Config, topic_num string, values url.Values) string {
	if values == nil {
		return tool.Get_redirect("/bbs/edit/" + tool.Url_parser(thread_bbs_id) + "/" + tool.Url_parser(topic_num))
	}

	db := tool.DB_connect()
	defer tool.DB_close(db)

	api_data := Api_thread_change_post(config, topic_num, values.Get("title"), values.Get("sub"))
	if api_data["response"] == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	}
	if api_data["response"] != "ok" {
		return tool.Get_error_page(db, config, "error")
	}
	return tool.Get_redirect("/bbs/w/" + tool.Url_parser(thread_bbs_id) + "/" + tool.Url_parser(topic_num))
}
