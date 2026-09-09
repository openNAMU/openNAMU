package route

import "opennamu/route/tool"

func View_topic_list(config tool.Config, doc_name string, do_type string, num string) string {
	return tool.Get_redirect("/bbs/in/" + tool.Url_parser(thread_bbs_id))
}
