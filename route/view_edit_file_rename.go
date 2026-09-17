package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_edit_file_rename(config tool.Config, doc_name string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if values == nil && !tool.Check_acl(db, doc_name, "", "document_move", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	old_base_name, extension, valid := File_rename_parts(doc_name)
	if !valid {
		return tool.Get_error_page(db, config, "invalid file")
	}

	if values != nil {
		result := Api_edit_file_rename_post(config, doc_name, values)
		if result["response"] == "require auth" {
			return tool.Get_error_page(db, config, "auth")
		}
		if result["response"] != "ok" {
			error_name, _ := result["data"].(string)
			return tool.Get_error_page(db, config, error_name)
		}
		new_doc_name, _ := result["data"].(string)
		return tool.Get_redirect("/w/" + tool.Url_parser(new_doc_name))
	}

	if _, exists := tool.Get_data_title(db, doc_name); !exists {
		return tool.Get_redirect("/list/file")
	}

	body := `<p>` + tool.HTML_escape(doc_name) + `</p><form method="post">`
	body += `<input name="name" value="` + tool.HTML_escape(old_base_name) + `" placeholder="` + tool.Get_language(db, "file_name", true) + `">.` + tool.HTML_escape(extension) + `<hr class="main_hr">`
	body += `<input name="send" placeholder="` + tool.Get_language(db, "why", true) + `"><hr class="main_hr">`
	body += tool.Get_captcha_ui(db, config) + tool.Get_IP_warning_ui(db, config) + tool.Get_edit_check_box_ui(db) + tool.Get_edit_bottom_text_ui(db, "move")
	body += `<button type="submit">` + tool.Get_language(db, "file_rename", true) + `</button></form>`

	return tool.Get_template(
		db,
		config,
		doc_name,
		body,
		[]any{"(" + tool.Get_language(db, "file_rename", true) + ")"},
		[][]any{{"w/" + tool.Url_parser(doc_name), tool.Get_language(db, "return", true)}},
		map[string]string{},
	)
}
