package route

import (
	"opennamu/route/tool"

	"net/url"
)

func View_login_register_post_full(config tool.Config, id string, password string, password_check string, captcha string) string {
	return user_register_post(config, id, password, password_check, captcha)
}

func View_login_register_values(config tool.Config, values url.Values) string {
	return user_register_post(config, values.Get("id"), values.Get("password"), values.Get("password_check"), values.Get("g-recaptcha-response"))
}
