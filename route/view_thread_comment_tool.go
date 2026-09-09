package route

import "opennamu/route/tool"

func View_thread_comment_tool(config tool.Config, topic_num string, comment_num string) string {
	return tool.Get_redirect("/bbs/tool/" + tool.Url_parser(thread_bbs_id) + "/" + tool.Url_parser(topic_num) + "/" + tool.Url_parser(comment_num))
}
