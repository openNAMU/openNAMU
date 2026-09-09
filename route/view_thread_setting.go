package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_thread_setting(config tool.Config, topic_num string, values url.Values) string {
	if values == nil {
		return tool.Get_redirect("/bbs/tool/" + tool.Url_parser(thread_bbs_id) + "/" + tool.Url_parser(topic_num))
	}

	db := tool.DB_connect()
	defer tool.DB_close(db)

	api_data := Api_thread_setting_post(config, topic_num, values.Get("stop"), values.Get("agree"), values.Get("why"))
	if api_data["response"] == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	}
	if api_data["response"] != "ok" {
		return tool.Get_redirect("/bbs/main")
	}
	return tool.Get_redirect("/bbs/tool/" + tool.Url_parser(thread_bbs_id) + "/" + tool.Url_parser(topic_num))
}
