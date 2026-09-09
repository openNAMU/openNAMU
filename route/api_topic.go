package route

import "opennamu/route/tool"

func Api_topic(config tool.Config, tool_name string, topic_num string, s_num string, e_num string) map[string]any {
	return api_topic_bbs(config, tool_name, topic_num, s_num, e_num)
}
