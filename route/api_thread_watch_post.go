package route

import "opennamu/route/tool"

func Api_thread_watch_post(config tool.Config, topic_num string) map[string]any {
	return Api_bbs_watch_post(config, thread_bbs_id, topic_num)
}
