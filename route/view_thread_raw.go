package route

import "opennamu/route/tool"

func View_thread_raw(config tool.Config, topic_num string, comment_num string) string {
	path := "/bbs/raw/" + tool.Url_parser(thread_bbs_id) + "/" + tool.Url_parser(topic_num)
	if comment_num != "" {
		path += "/" + tool.Url_parser(comment_num)
	}
	return tool.Get_redirect(path)
}
