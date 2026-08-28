package route

import "opennamu/route/tool"

func View_list_recent_block(config tool.Config, num string, set_type string, why string, user_name string) string {
	return View_auth_give_list(config, num, false)
}
