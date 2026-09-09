package route

import "opennamu/route/tool"

func Api_thread_comment_delete_post(config tool.Config, topic_num string, comment_num string) map[string]any {
	return Api_bbs_w_comment_one_delete(config, thread_bbs_id, topic_num+"-"+comment_num)
}
