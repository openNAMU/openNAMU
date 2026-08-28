package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_delete_multiple(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "document_bulk_delete", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	if values != nil {
		if !tool.Captcha_check(db, config.Session, config.IP, tool.Captcha_response(values.Get("g-recaptcha"), values.Get("g-recaptcha-response"), values.Get("h-captcha-response"), values.Get("cf-turnstile-response"))) {
			return tool.Get_error_page(db, config, "recaptcha")
		}
		api_data := Api_edit_delete_multiple_post(config, values.Get("content"), values.Get("send"), values.Get("copyright_agreement"))
		response, _ := api_data["response"].(string)
		if response == "require auth" {
			return tool.Get_error_page(db, config, "auth")
		}
		if response != "ok" {
			error_name, _ := api_data["data"].(string)
			if error_name == "" {
				error_name = "error"
			}
			return tool.Get_error_page(db, config, error_name)
		}
		return tool.Get_redirect("/recent_change")
	}

	data := `<form method="post"><textarea class="opennamu_textarea_500" name="content" placeholder="` + tool.Get_language(db, "many_delete_help", true) + `"></textarea><hr class="main_hr"><input name="send" placeholder="` + tool.Get_language(db, "why", true) + `"><hr class="main_hr">` + tool.Get_captcha_ui(db, config) + tool.Get_IP_warning_ui(db, config) + tool.Get_edit_check_box_ui(db) + tool.Get_edit_bottom_text_ui(db, "edit") + `<button type="submit">` + tool.Get_language(db, "delete", true) + `</button></form>`
	return tool.Get_template(db, config, tool.Get_language(db, "many_delete", true), data, []any{}, [][]any{{"manager", tool.Get_language(db, "return", true)}}, map[string]string{})
}
