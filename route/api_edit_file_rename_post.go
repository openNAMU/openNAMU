package route

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"opennamu/route/tool"
)

func file_rename_parts(doc_name string) (string, string, bool) {
	if !strings.HasPrefix(doc_name, "file:") {
		return "", "", false
	}

	file_name := strings.TrimPrefix(doc_name, "file:")
	extension := strings.TrimPrefix(strings.ToLower(filepath.Ext(file_name)), ".")
	base_name := strings.TrimSuffix(file_name, filepath.Ext(file_name))
	if base_name == "" || extension == "" || strings.ContainsAny(base_name, `/\`) || strings.Contains(base_name, ".") {
		return "", "", false
	}

	return base_name, extension, true
}

func Api_edit_file_rename_post(config tool.Config, doc_name string, values url.Values) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_acl(db, doc_name, "", "document_move", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	old_base_name, extension, valid := file_rename_parts(doc_name)
	if !valid {
		return_data["response"] = "error"
		return_data["data"] = "invalid file"
		return return_data
	}

	new_base_name := strings.TrimSpace(values.Get("name"))
	if new_base_name == "" || strings.ContainsAny(new_base_name, `/\`) || strings.Contains(new_base_name, ".") {
		return_data["response"] = "error"
		return_data["data"] = "unallowed file name"
		return return_data
	}

	new_file_name := new_base_name + "." + extension
	if tool.Get_file_name_unallow_check(db, new_file_name) {
		return_data["response"] = "error"
		return_data["data"] = "unallowed file name"
		return return_data
	}

	new_doc_name := "file:" + new_file_name
	if new_doc_name == doc_name {
		return_data["response"] = "error"
		return_data["data"] = "already exist"
		return return_data
	}
	if target_exists, _ := move_document_exists(db, new_doc_name); target_exists {
		return_data["response"] = "error"
		return_data["data"] = "already exist"
		return return_data
	}

	file_main_dir := tool.Get_file_main_dir(db)
	old_path := filepath.Join(file_main_dir, tool.File_name_to_dir(old_base_name, extension))
	new_path := filepath.Join(file_main_dir, tool.File_name_to_dir(new_base_name, extension))
	if _, err := os.Stat(old_path); err != nil {
		return_data["response"] = "error"
		return_data["data"] = "invalid file"
		return return_data
	}
	if _, err := os.Stat(new_path); err == nil {
		return_data["response"] = "error"
		return_data["data"] = "already exist"
		return return_data
	}

	if err := os.Rename(old_path, new_path); err != nil {
		return_data["response"] = "error"
		return_data["data"] = "file rename error"
		return return_data
	}

	move_values := url.Values{}
	for key, value := range values {
		move_values[key] = value
	}
	move_values.Set("title", new_doc_name)
	move_values.Set("move_option", "normal")
	move_values.Set("move_topic_option", "none")
	move_values.Set("document_set_option", "normal")

	move_data := Api_edit_move_post(config, doc_name, move_values)
	if move_data["response"] != "ok" {
		_ = os.Rename(new_path, old_path)
		return move_data
	}

	return_data["response"] = "ok"
	return_data["data"] = new_doc_name
	return return_data
}
