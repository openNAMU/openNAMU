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
	invite_hash, _ := config.Session.Get("submit_invite").(string)
	if id == "" || pw == "" {
		return tool.Get_redirect("/register")
	}
	question := User_other(db, "approval_question")
	if question == "" {
		for _, name := range []string{"submit_id", "submit_pw", "submit_email"} {
			config.Session.Delete(name)
		}
		_ = config.Session.Save()
		return tool.Get_redirect("/register")
	}
	if values != nil {
		result := Api_register_submit_post(config, id, pw, email, question, values.Get("answer"), invite_hash)
		if result["response"] != "ok" {
			if result["data"] == "invite error" {
				return tool.Get_error_page(db, config, "invite error")
			}
			return tool.Get_error_page(db, config, "error")
		}
		for _, name := range []string{"submit_id", "submit_pw", "submit_email", "submit_invite"} {
			config.Session.Delete(name)
		}
		_ = config.Session.Save()
		return User_form_page(db, config, tool.Get_language(db, "register", true), `<p>`+tool.Get_language(db, "register_submit_done", true)+`</p><div><a href="/user">`+tool.Get_language(db, "return", true)+`</a></div>`)
	}
	body := "<form method='post'><p><label for='register_answer'>" + tool.HTML_escape(question) + "</label></p><div><input id='register_answer' name='answer'></div><hr class='main_hr'><div><button type='submit'>" + tool.Get_language(db, "send", true) + "</button></div></form>"
	return User_form_page(db, config, tool.Get_language(db, "register", true), body)
}
