package route

import "opennamu/route/tool"

func View_login_register_post_full(config tool.Config, id string, password string, password_check string, captcha string) string {
	return user_register_post(config, id, password, password_check, captcha)
}
