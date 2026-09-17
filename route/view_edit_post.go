package route

import (
	"opennamu/route/tool"
)

func View_edit_post(config tool.Config, doc_name string, data string, send string, agree string, captcha string, expected_revision string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !tool.Captcha_check(db, config.Session, config.IP, captcha) {
		return tool.Get_error_page(db, config, "recaptcha")
	}

	return_data := Api_edit_post(config, doc_name, data, send, agree, expected_revision)

	result_html := ""
	response, _ := return_data["response"].(string)
	if response == "ok" {
		result_html = tool.Get_redirect("/w/" + tool.Url_parser(doc_name))
	} else {
		error_name, _ := return_data["data"].(string)
		result_html = tool.Get_error_page(db, config, error_name)
	}

	return result_html
}
