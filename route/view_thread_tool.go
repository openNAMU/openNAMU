package route

import "opennamu/route/tool"

func View_thread_tool(config tool.Config, topic_num string) string {
	return tool.Get_redirect("/bbs/tool/" + tool.Url_parser(thread_bbs_id) + "/" + tool.Url_parser(topic_num))
}
