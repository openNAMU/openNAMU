package route

import "opennamu/route/tool"

func View_history_report(config tool.Config, doc_name string, rev string, reason string, captcha string, submit bool) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	target_api := Api_history_report_target(config, doc_name, rev)
	response, _ := target_api["response"].(string)
	if response == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	}
	if response == "not exist" {
		return tool.Get_redirect("/history/" + tool.Url_parser(doc_name))
	}
	if response != "ok" {
		return tool.Get_error_page(db, config, "error")
	}
	revision, _ := target_api["data"].(string)
	if revision == "" {
		return tool.Get_error_page(db, config, "error")
	}

	target_path := "/history_tool/" + revision + "/" + tool.Url_parser(doc_name)
	target_title := doc_name + " (r" + revision + ")"
	if submit {
		api_data := Api_history_report_post(config, doc_name, revision, reason, captcha)
		response, _ := api_data["response"].(string)
		if response == "ok" {
			body := `<div>` + tool.Get_language(db, "report_submitted", true) + `</div><hr class="main_hr"><a href="` + tool.HTML_escape(target_path) + `">` + tool.HTML_escape(target_title) + `</a>`
			return tool.Get_template(
				db,
				config,
				tool.Get_language(db, "report", true),
				body,
				[]any{},
				[][]any{{"history_tool/" + revision + "/" + tool.Url_parser(doc_name), tool.Get_language(db, "return", true)}},
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

	form_path := "/history_report/" + revision + "/" + tool.Url_parser(doc_name)
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
		[][]any{{"history_tool/" + revision + "/" + tool.Url_parser(doc_name), tool.Get_language(db, "return", true)}},
		map[string]string{},
	)
}
