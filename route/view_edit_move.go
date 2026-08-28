package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_edit_move(config tool.Config, doc_name string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if values != nil {
		api_data := Api_edit_move_post(config, doc_name, values)
		response, _ := api_data["response"].(string)
		if response == "require auth" {
			return tool.Get_error_page(db, config, "auth")
		}
		if response == "not exist" {
			return tool.Get_redirect("/w/" + tool.Url_parser(doc_name))
		}
		if response != "ok" {
			error_name, _ := api_data["data"].(string)
			if error_name == "" {
				error_name = "error"
			}
			return tool.Get_error_page(db, config, error_name)
		}
		new_name, ok := api_data["data"].(string)
		if !ok || new_name == "" {
			return tool.Get_error_page(db, config, "error")
		}
		return tool.Get_redirect("/w/" + tool.Url_parser(new_name))
	}

	if !tool.Check_acl(db, doc_name, "", "document_move", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	var source_title string
	if !tool.QueryRow_DB(db, "select title from data where title = ?", []any{&source_title}, doc_name) {
		return tool.Get_redirect("/w/" + tool.Url_parser(doc_name))
	}
	owner_auth := tool.Check_acl(db, "", "", "owner_auth", config.IP)
	body := "<form method=\"post\">"
	body += "<input name=\"title\" value=\"" + tool.HTML_escape(doc_name) + "\" placeholder=\"" + tool.Get_language(db, "document_name", true) + "\"><hr class=\"main_hr\">"
	body += "<input name=\"send\" placeholder=\"" + tool.Get_language(db, "why", true) + "\"><hr class=\"main_hr\">"

	body += "<h2>" + tool.Get_language(db, "document", true) + "</h2>"
	body += "<select name=\"move_option\">"
	body += "<option value=\"normal\" selected>" + tool.Get_language(db, "normal", true) + "</option>"
	body += "<option value=\"none\">" + tool.Get_language(db, "dont_move", true) + "</option>"
	body += "<option value=\"reverse\">" + tool.Get_language(db, "replace_move", true) + "</option>"
	if owner_auth {
		body += "<option value=\"merge\">" + tool.Get_language(db, "merge_move", true) + "</option>"
	}
	body += "</select><hr class=\"main_hr\">"

	body += "<h2>" + tool.Get_language(db, "discussion", true) + "</h2>"
	body += "<select name=\"move_topic_option\">"
	body += "<option value=\"none\" selected>" + tool.Get_language(db, "dont_move", true) + "</option>"
	body += "<option value=\"normal\">" + tool.Get_language(db, "normal", true) + "</option>"
	body += "<option value=\"reverse\">" + tool.Get_language(db, "replace_move", true) + "</option>"
	if owner_auth {
		body += "<option value=\"merge\">" + tool.Get_language(db, "merge_move", true) + "</option>"
	}
	body += "</select><hr class=\"main_hr\">"

	if owner_auth {
		body += "<h2>" + tool.Get_language(db, "document_set", true) + "</h2>"
		body += "<select name=\"document_set_option\">"
		body += "<option value=\"none\" selected>" + tool.Get_language(db, "dont_move", true) + "</option>"
		body += "<option value=\"normal\">" + tool.Get_language(db, "normal", true) + "</option>"
		body += "<option value=\"reverse\">" + tool.Get_language(db, "replace_move", true) + "</option>"
		body += "</select><hr class=\"main_hr\">"
	}

	body += tool.Get_captcha_ui(db, config) + tool.Get_IP_warning_ui(db, config) + tool.Get_edit_check_box_ui(db) + tool.Get_edit_bottom_text_ui(db, "move")
	body += "<button type=\"submit\">" + tool.Get_language(db, "move", true) + "</button></form>"
	return tool.Get_template(db, config, doc_name, body, []any{"(" + tool.Get_language(db, "move", true) + ")"}, [][]any{{"w/" + tool.Url_parser(doc_name), tool.Get_language(db, "return", true)}, {"move_all", tool.Get_language(db, "multiple_move", true)}}, map[string]string{})
}
