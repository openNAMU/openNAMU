package route

import (
	"net/url"
	"path/filepath"
	"strings"

	"opennamu/route/tool"
)

func file_storage_path(db_file_dir string, doc_name string) (string, bool) {
	if !strings.HasPrefix(doc_name, "file:") {
		return "", false
	}
	file_name := strings.TrimPrefix(doc_name, "file:")
	extension := strings.TrimPrefix(strings.ToLower(filepath.Ext(file_name)), ".")
	base_name := strings.TrimSuffix(file_name, filepath.Ext(file_name))
	if base_name == "" || extension == "" || strings.ContainsAny(base_name, "/"+string(rune(92))) || strings.Contains(base_name, ".") {
		return "", false
	}
	return filepath.Join(db_file_dir, tool.File_name_to_dir(base_name, extension)), true
}

func View_edit_file_delete(config tool.Config, doc_name string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if values == nil && !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	var stored_name string
	if !tool.QueryRow_DB(db, "select title from data where title = ?", []any{&stored_name}, doc_name) && values == nil {
		return tool.Get_redirect("/list/file")
	}

	file_path, valid := file_storage_path(tool.Get_file_main_dir(db), doc_name)
	if !valid {
		return tool.Get_error_page(db, config, "invalid file")
	}
	if values != nil {
		result := Api_edit_file_delete_post(config, doc_name, values)
		if result["response"] == "require auth" {
			return tool.Get_error_page(db, config, "auth")
		}
		if result["response"] != "ok" {
			if result["data"] == "delete error" {
				return tool.Get_error_page(db, config, "delete error")
			}
			return tool.Get_error_page(db, config, "file delete error")
		}
		if values.Get("with_doc") != "" {
			return tool.Get_redirect("/w/" + tool.Url_parser(doc_name))
		}
		return tool.Get_redirect("/list/file")
	}

	body := "<p><a href='/image/" + tool.Url_parser(filepath.Base(file_path)) + "'>" + tool.HTML_escape(doc_name) + "</a></p>"
	body += "<form method='post'><input type='checkbox' name='with_doc' value='1' checked> " + tool.Get_language(db, "file_delete_with_document", true)
	body += "<input name='send' placeholder='" + tool.Get_language(db, "why", true) + "'><hr class='main_hr'>"
	body += tool.Get_edit_check_box_ui(db) + "<button type='submit'>" + tool.Get_language(db, "delete", true) + "</button></form>"
	return tool.Get_template(db, config, doc_name, body, []any{"(" + tool.Get_language(db, "delete", true) + ")"}, [][]any{{"w/" + tool.Url_parser(doc_name), tool.Get_language(db, "return", true)}}, map[string]string{})
}
