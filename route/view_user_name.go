package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_user_name(config tool.Config, values url.Values) string {
	return View_user_name_for(config, "", values)
}
