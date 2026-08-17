package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_user_head(config tool.Config, values url.Values) string {
	return View_user_head_skin(config, "", values)
}
