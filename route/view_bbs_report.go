package route

import "opennamu/route/tool"

func View_bbs_report(config tool.Config, set_id string, set_code string, comment_code string, reason string, captcha string, submit bool) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if set_id == report_bbs_id {
		return tool.Get_redirect("/bbs/main")
	}
	if !tool.Check_permission(db, "bbs_comment", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	target_api := Api_bbs_w(config, set_id, set_code)
	if target_api["response"] == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	}
	target_data, _ := target_api["data"].(map[string]string)
	if target_api["response"] != "ok" || target_data["title"] == "" {
		return tool.Get_redirect("/bbs/main")
	}

	target_path := "/bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code)
	target_title := target_data["title"]
	if comment_code != "" {
		if !bbs_comment_code_regex.MatchString(comment_code) {
			return tool.Get_redirect("/bbs/main")
		}
		comment_api := Api_bbs_w_comment_one(config, false, "", set_id+"-"+set_code+"-"+comment_code)
		comment_list, _ := comment_api["data"].([]map[string]string)
		if comment_api["response"] == "require auth" {
			return tool.Get_error_page(db, config, "auth")
		}
		if comment_api["response"] != "ok" || len(comment_list) == 0 || comment_list[0]["comment"] == "" {
			return tool.Get_redirect("/bbs/main")
		}
		target_path += "/comment/" + tool.Url_parser(comment_code)
		target_title += " #" + comment_code
	}

	if submit {
		api_data := Api_bbs_report_post(config, set_id, set_code, comment_code, reason, captcha)
		response, _ := api_data["response"].(string)
		if response == "ok" {
			body := `<div>` + tool.Get_language(db, "report_submitted", true) + `</div><hr class="main_hr"><a href="` + tool.HTML_escape(target_path) + `">` + tool.HTML_escape(target_title) + `</a>`
			return tool.Get_template(
				db,
				config,
				tool.Get_language(db, "report", true),
				body,
				[]any{},
				[][]any{{"bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code), tool.Get_language(db, "return", true)}},
				map[string]string{},
			)
		}
		if response == "require auth" {
			return tool.Get_error_page(db, config, "auth")
		}
		if api_data["data"] == "recaptcha" {
			return tool.Get_error_page(db, config, "recaptcha")
		}
		if api_data["data"] == "daily limit" {
			return tool.Get_error_page(db, config, "daily limit")
		}
		if api_data["data"] != "report reason" {
			return tool.Get_error_page(db, config, "error")
		}
	}

	form_path := "/bbs/report/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code)
	if comment_code != "" {
		form_path += "/comment/" + tool.Url_parser(comment_code)
	}
	body := `<div>` + tool.Get_language(db, "report_target", true) + `: <a href="` + tool.HTML_escape(target_path) + `">` + tool.HTML_escape(target_title) + `</a></div><hr class="main_hr">`
	if submit {
		body += `<div>` + tool.Get_language(db, "report_reason_error", true) + `</div><hr class="main_hr">`
	}
	body += `<h3>` + tool.Get_language(db, "report_reason", true) + `</h3><form method="post" action="` + tool.HTML_escape(form_path) + `"><textarea class="opennamu_textarea_100" name="reason">` + tool.HTML_escape(reason) + `</textarea><hr class="main_hr">` + tool.Get_captcha_ui(db, config) + `<hr class="main_hr"><button type="submit">` + tool.Get_language(db, "send", true) + `</button></form>`

	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, "report", true),
		body,
		[]any{},
		[][]any{{"bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code), tool.Get_language(db, "return", true)}},
		map[string]string{},
	)
}
