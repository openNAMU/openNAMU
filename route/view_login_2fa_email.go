package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_login_2fa_email(config tool.Config, values url.Values) string {
	login_id, _ := config.Session.Get("login_id").(string)
	if login_id == "" {
		if email_id, ok := config.Session.Get("b_id").(string); ok && email_id != "" {
			config.Session.Set("login_id", email_id)
			_ = config.Session.Save()
		}
	}
	return View_login_2fa(config, values)
}
