package route

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"opennamu/route/tool"
	"opennamu/route/tool/markup"
)

func Api_file_upload_post(config tool.Config, file_name string, file_data string, file_ext string) map[string]any {
	return api_file_upload_post(config, file_name, file_data, file_ext, "direct_input", "", "", false, false)
}

func api_file_upload_make_document(db *sql.DB, doc_name string, doc_data string, ip string) bool {
	if db == nil {
		return false
	}

	if _, err := db.Exec(tool.DB_change("insert into data (title, data) values (?, ?)"), doc_name, doc_data); err != nil {
		return false
	}
	markup.Get_render(db, doc_name, doc_data, "backlink")
	tool.Do_add_history(db, doc_name, doc_data, tool.Get_time(), ip, "", "0", "upload", "")

	return true
}

func api_file_upload_post(config tool.Config, file_name string, file_data string, file_ext string, license string, license_text string, captcha string, check_captcha bool, many_upload bool) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	file_name = strings.TrimSpace(file_name)
	file_data = strings.TrimSpace(file_data)
	file_ext = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(file_ext)), ".")

	allowed_ext := tool.Get_ext_allow_list(db)
	return_value := make(map[string]any)

	if file_data == "" || file_name == "" || file_ext == "" {
		return_value["response"] = "error"
		return_value["data"] = "invalid data"
		return return_value
	} else if strings.ContainsAny(file_name, `/\`) || strings.Contains(file_name, ".") {
		return_value["response"] = "error"
		return_value["data"] = "unallowed file name"
		return return_value
	} else if _, ok := allowed_ext[file_ext]; !ok {
		return_value["response"] = "error"
		return_value["data"] = "unallowed ext"
		return return_value
	} else if tool.Get_file_name_unallow_check(db, file_name+"."+file_ext) {
		return_value["response"] = "error"
		return_value["data"] = "unallowed file name"
		return return_value
	} else if !tool.Check_permission(db, "upload", config.IP) || (many_upload && !tool.Check_permission(db, "multiple_upload", config.IP)) {
		return_value["response"] = "require auth"
		return return_value
	} else if check_captcha && !tool.Captcha_check(db, config.Session, config.IP, captcha) {
		return_value["response"] = "error"
		return_value["data"] = "recaptcha"
		return return_value
	}

	decoded, err := tool.Get_base64_decode(file_data)
	if err != nil || len(decoded) == 0 {
		return_value["response"] = "error"
		return_value["data"] = "invalid data"
		return return_value
	}

	file_max_size := tool.Get_file_max_size_by_extension(db, file_ext)
	if file_max_size <= 0 {
		file_max_size = tool.Get_file_max_size(db)
	}
	if file_max_size <= 0 {
		file_max_size = 2
	}
	if len(decoded) > file_max_size*1000*1000 {
		return_value["response"] = "error"
		return_value["data"] = "file too large"
		return return_value
	}

	doc_name := "file:" + file_name + "." + file_ext
	var old_doc_name string
	if tool.QueryRow_DB(db, "select title from data where title = ?", []any{&old_doc_name}, doc_name) {
		return_value["response"] = "error"
		return_value["data"] = "already exist"
		return return_value
	}

	main_dir := tool.Get_file_main_dir(db)
	if err := os.MkdirAll(main_dir, 0o755); err != nil {
		return_value["response"] = "error"
		return_value["data"] = "directory create fail"
		return return_value
	}

	file_full_dir := tool.File_name_to_dir(file_name, file_ext)
	dst_path := filepath.Join(main_dir, file_full_dir)
	if _, err := os.Stat(dst_path); err == nil {
		return_value["response"] = "error"
		return_value["data"] = "already exist"
		return return_value
	} else if !errors.Is(err, os.ErrNotExist) {
		return_value["response"] = "error"
		return_value["data"] = "exist check fail"
		return return_value
	}

	out, err := os.OpenFile(dst_path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return_value["response"] = "error"
		return_value["data"] = "file create fail"
		return return_value
	}

	if _, err := out.Write([]byte(decoded)); err != nil {
		_ = out.Close()
		_ = os.Remove(dst_path)
		return_value["response"] = "error"
		return_value["data"] = "file write fail"
		return return_value
	}
	if err := out.Close(); err != nil {
		_ = os.Remove(dst_path)
		return_value["response"] = "error"
		return_value["data"] = "file write fail"
		return return_value
	}

	if license == "" {
		license = "direct_input"
	}
	license = strings.ReplaceAll(license, "]", "_")
	doc_data := license + "\n"
	if tool.Get_document_markup(db, doc_name, "document") == "namumark" {
		doc_data += "[[category:" + license + "]]\n"
	}
	doc_data += license_text

	if !api_file_upload_make_document(db, doc_name, doc_data, config.IP) {
		_ = os.Remove(dst_path)
		return_value["response"] = "error"
		return_value["data"] = "document create fail"
		return return_value
	}

	return_value["response"] = "ok"
	return_value["data"] = doc_name
	return return_value
}
