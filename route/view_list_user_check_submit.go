package route

import "opennamu/route/tool"

func View_list_user_check_submit(config tool.Config, user_name string) string {
	return View_user_check(config, user_name, "check", "1", "")
}
