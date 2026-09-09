package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_thread_watch(config tool.Config, topic_num string, values url.Values) string {
	if values == nil {
		return tool.Get_redirect("/bbs_watch/" + tool.Url_parser(thread_bbs_id) + "/" + tool.Url_parser(topic_num))
	}
	return View_bbs_watch(config, thread_bbs_id, topic_num, values)
}
