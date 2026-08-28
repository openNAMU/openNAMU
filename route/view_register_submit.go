package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_register_submit(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	id, _ := config.Session.Get("submit_id").(string)
	pw, _ := config.Session.Get("submit_pw").(string)
	email, _ := config.Session.Get("submit_email").(string)
	if id == "" || pw == "" {
		return tool.Get_redirect("/register")
	}
	question := user_other(db, "approval_question")
	if question == "" {
		for _, name := range []string{"submit_id", "submit_pw", "submit_email"} {
			config.Session.Delete(name)
		}
		_ = config.Session.Save()
		return tool.Get_redirect("/register")
	}
	if values != nil {
		result := Api_register_submit_post(config, id, pw, email, question, values.Get("answer"))
		if result["response"] != "ok" {
			return tool.Get_error_page(db, config, "error")
		}
		for _, name := range []string{"submit_id", "submit_pw", "submit_email"} {
			config.Session.Delete(name)
		}
		_ = config.Session.Save()
		return user_form_page(db, config, "register submit", "<p>submitted</p><a href='/user'>return</a>")
	}
	body := "<form method='post'><p>" + tool.HTML_escape(question) + "</p><input name='answer'><button type='submit'>send</button></form>"
	return user_form_page(db, config, "register submit", body)
}
