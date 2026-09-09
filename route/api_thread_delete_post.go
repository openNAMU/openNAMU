package route

import "opennamu/route/tool"

func Api_thread_delete_post(config tool.Config, topic_num string) map[string]any {
	return Api_bbs_w_delete(config, thread_bbs_id, topic_num)
}
