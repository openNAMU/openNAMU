package route

import "opennamu/route/tool"

func Api_list_recent_discuss(config tool.Config, limit string, num string, set_type string) map[string]any {
	return api_list_recent_discuss_bbs(config, limit, num, set_type)
}
