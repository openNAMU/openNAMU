package route

import "opennamu/route/tool"

func View_list_recent_discuss(config tool.Config, limit string, num string, set_type string) string {
	return tool.Get_redirect("/bbs/in/" + tool.Url_parser(thread_bbs_id))
}
