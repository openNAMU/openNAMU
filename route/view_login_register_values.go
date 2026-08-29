package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_login_register_values(config tool.Config, values url.Values) string {
	return user_register_post(config, values.Get("id"), values.Get("password"), values.Get("password_check"), tool.Captcha_response(values.Get("g-recaptcha"), values.Get("g-recaptcha-response"), values.Get("h-captcha-response"), values.Get("cf-turnstile-response"), values.Get("altcha")))
}
