package route

import (
	"opennamu/route/tool"
)

func View_edit_delete_post(config tool.Config, doc_name string, send string, agree string, captcha string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !tool.Check_acl(db, doc_name, "", "document_delete", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	if !tool.Captcha_check(db, config.Session, config.IP, captcha) {
		return tool.Get_error_page(db, config, "recaptcha")
	}

	api_data := Api_edit_delete_post(config, doc_name, send, agree)
	response := api_data["response"].(string)
	switch response {
	case "ok", "not exist":
		return tool.Get_redirect("/w/" + tool.Url_parser(doc_name))
	case "require auth":
		return tool.Get_error_page(db, config, "auth")
	}

	error_name, ok := api_data["data"].(string)
	if !ok {
		error_name = "error"
	}

	return tool.Get_error_page(db, config, error_name)
}
