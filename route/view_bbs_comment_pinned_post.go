package route

import "opennamu/route/tool"

func View_bbs_comment_pinned_post(config tool.Config, set_id string, set_code string, comment_code string) string {
	Api_bbs_w_comment_pinned(config, set_id, set_code, comment_code, true)
	return tool.Get_redirect("/bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code) + "#" + tool.Url_parser(comment_code))
}
