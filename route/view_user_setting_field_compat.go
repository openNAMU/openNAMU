package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_user_setting_field_compat(config tool.Config, field string, values url.Values) string {
	return View_user_setting_field(config, field, user_field_values(values, field))
}
