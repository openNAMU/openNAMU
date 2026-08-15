package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_thread_route(config tool.Config, topic_num string, doc_name string, values url.Values) string {
	db := tool.DB_connect()
	if values != nil {
		captcha := tool.Captcha_response(values.Get("g-recaptcha"), values.Get("g-recaptcha-response"), values.Get("h-captcha-response"), values.Get("cf-turnstile-response"))
		if !tool.Captcha_check(db, config.Session, config.IP, captcha) {
			data := tool.Get_error_page(db, config, "recaptcha_error")
			tool.DB_close(db)
			return data
		}
	}
	tool.DB_close(db)

	return View_thread(config, topic_num, doc_name, values)
}
