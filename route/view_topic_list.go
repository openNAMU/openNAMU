package route

import (
	"strconv"

	"opennamu/route/tool"
)

func View_topic_list(config tool.Config, doc_name string, do_type string, num string) string {
	if doc_name == "" {
		return tool.Get_redirect("/bbs/in/" + tool.Url_parser(thread_bbs_id))
	}

	page := tool.Str_to_int(num)
	if page < 1 {
		page = 1
	}
	return tool.Get_redirect(
		"/bbs/in/" + tool.Url_parser(thread_bbs_id) +
			"/filter/tag/" + tool.Url_parser(doc_name) + "/" + strconv.Itoa(page),
	)
}
