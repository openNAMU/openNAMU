package route

import "opennamu/route/tool"

func Api_topic_list(config tool.Config, num string, doc_name string, do_type string) map[string]any {
	return api_topic_list_bbs(config, num, doc_name, do_type)
}
