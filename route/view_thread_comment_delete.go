package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_thread_comment_delete(config tool.Config, topic_num string, comment_num string, values url.Values) string {
	if values == nil {
		return tool.Get_redirect("/bbs/delete/" + tool.Url_parser(thread_bbs_id) + "/" + tool.Url_parser(topic_num) + "/" + tool.Url_parser(comment_num))
	}
	return View_bbs_delete_post(config, thread_bbs_id, topic_num, comment_num)
}
