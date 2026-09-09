package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_thread_comment_notice(config tool.Config, topic_num string, comment_num string, values url.Values) string {
	if values == nil {
		return tool.Get_redirect("/bbs/tool/" + tool.Url_parser(thread_bbs_id) + "/" + tool.Url_parser(topic_num) + "/" + tool.Url_parser(comment_num))
	}

	db := tool.DB_connect()
	defer tool.DB_close(db)

	api_data := Api_thread_comment_notice_post(config, topic_num, comment_num)
	if api_data["response"] == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	}
	return tool.Get_redirect("/bbs/w/" + tool.Url_parser(thread_bbs_id) + "/" + tool.Url_parser(topic_num) + "#" + tool.Url_parser(comment_num))
}
