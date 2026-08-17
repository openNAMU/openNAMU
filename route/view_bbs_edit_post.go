package route

import "opennamu/route/tool"

func View_bbs_edit_post(config tool.Config, set_id string, set_code string, comment_code string, title string, data string, captcha string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !tool.Captcha_check(db, config.Session, config.IP, captcha) {
		return tool.Get_error_page(db, config, "recaptcha")
	}

	api_data := Api_bbs_w_edit_post(config, set_id, set_code, comment_code, title, data)
	response, _ := api_data["response"].(string)

	if response == "ok" {
		new_code, _ := api_data["data"].(string)
		target := "/bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(new_code)
		if comment_code != "" {
			target += "#" + tool.Url_parser(comment_code)
		}

		return tool.Get_redirect(target)
	}

	if response == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	}

	if response == "not exist" {
		return tool.Get_redirect("/bbs/main")
	}

	error_name, _ := api_data["data"].(string)
	if error_name == "empty data" {
		return tool.Get_redirect("/bbs/in/" + tool.Url_parser(set_id))
	}
	if error_name == "" {
		error_name = "error"
	}

	return tool.Get_error_page(db, config, error_name)
}
