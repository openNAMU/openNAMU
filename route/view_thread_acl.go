package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_thread_acl(config tool.Config, topic_num string, values url.Values) string {
	return tool.Get_redirect("/bbs/set/" + tool.Url_parser(thread_bbs_id))
}
