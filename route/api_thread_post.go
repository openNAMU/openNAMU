package route

import "opennamu/route/tool"

func Api_thread_post(config tool.Config, topic_num string, doc_name string, content string, topic string, title string) map[string]any {
	return api_thread_bbs_post(config, topic_num, doc_name, content, topic, title)
}
