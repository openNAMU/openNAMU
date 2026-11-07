package route

import (
	"net/url"
	"path/filepath"

	"opennamu/route/tool"
)

func View_edit_file_delete(config tool.Config, doc_name string, values url.Values) string {
	if values != nil {
		result := Api_edit_file_delete_post(config, doc_name, values)
		response, _ := result["response"].(string)
		if response == "require auth" {
			db := tool.DB_connect()
			defer tool.DB_close(db)
			if _, valid := File_storage_path(tool.Get_file_main_dir(db), doc_name); !valid {
				return tool.Get_error_page(db, config, "invalid file")
			}
			return tool.Get_error_page(db, config, "auth")
		}
		if response != "ok" {
			error_name, _ := result["data"].(string)
			if error_name != "delete error" && error_name != "invalid file" {
				error_name = "file delete error"
			}
			db := tool.DB_connect()
			defer tool.DB_close(db)
			return tool.Get_error_page(db, config, error_name)
		}
		if values.Get("with_doc") != "" {
			return tool.Get_redirect("/w/" + tool.Url_parser(doc_name))
		}
		return tool.Get_redirect("/list/file")
	}

	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "file_delete", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	_, stored_exists := tool.Get_data_title(db, doc_name)
	if !stored_exists {
		return tool.Get_redirect("/list/file")
	}

	file_path, valid := File_storage_path(tool.Get_file_main_dir(db), doc_name)
	if !valid {
		return tool.Get_error_page(db, config, "invalid file")
	}
	body := "<p><a href='/image/" + tool.Url_parser(filepath.Base(file_path)) + "'>" + tool.HTML_escape(doc_name) + "</a></p>"
	body += "<form method='post'><input type='checkbox' name='with_doc' value='1' checked> " + tool.Get_language(db, "file_delete_with_document", true)
	body += "<hr class=\"main_hr\">"
	body += "<input name='send' placeholder='" + tool.Get_language(db, "why", true) + "'><hr class='main_hr'>"
	body += tool.Get_edit_check_box_ui(db) + "<button type='submit'>" + tool.Get_language(db, "delete", true) + "</button></form>"
	return tool.Get_template(db, config, doc_name, body, []any{"(" + tool.Get_language(db, "delete", true) + ")"}, [][]any{{"w/" + tool.Url_parser(doc_name), tool.Get_language(db, "return", true)}}, map[string]string{})
}
