package route

import (
	"opennamu/route/tool"
)

func View_user(config tool.Config, user_name string) string {
	return View_user_safe(config, user_name)
}
