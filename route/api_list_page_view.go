package route

import "opennamu/route/tool"

func Api_list_page_view(config tool.Config) map[string]any {
	return Api_list_title_index(config, "1")
}
