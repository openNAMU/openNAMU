package route

import "opennamu/route/tool"

func View_bbs_in_w_post_secure(config tool.Config, set_id string, set_code string, comment_select string, data string, captcha string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Captcha_check(db, config.Session, config.IP, captcha) {
		return tool.Get_error_page(db, config, "recaptcha_error")
	}

	api_data := Api_bbs_w_comment_post(config, set_id, set_code, comment_select, data)
	response, _ := api_data["response"].(string)
	if response == "ok" {
		comment_code, _ := api_data["data"].(string)
		return tool.Get_redirect("/bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code) + "#" + tool.Url_parser(comment_code))
	}

	if response == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	}
	if response == "error" {
		error_name, _ := api_data["data"].(string)
		if error_name == "" {
			error_name = "error"
		}
		return tool.Get_error_page(db, config, error_name)
	}

	return tool.Get_redirect("/bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code))
}
