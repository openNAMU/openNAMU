package route

import (
	"bytes"
	"database/sql"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"opennamu/route/tool"
	"opennamu/route/tool/markup"
)

func Api_file_upload_post(config tool.Config, file_name string, file_data string, file_ext string) map[string]any {
	decoded, err := tool.Get_base64_decode(strings.TrimSpace(file_data))
	if err != nil || len(decoded) == 0 {
		return map[string]any{"response": "error", "data": "invalid data"}
	}
	return api_file_upload_post(config, file_name, []byte(decoded), file_ext, "direct_input", "", "", false, false, false)
}

func api_file_upload_make_document(db *sql.DB, doc_name string, doc_data string, ip string) bool {
	if db == nil {
		return false
	}

	err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		if _, err := tx.Exec(tool.DB_change("insert into data (title, data) values (?, ?)"), doc_name, doc_data); err != nil {
			return err
		}
		tool.Do_add_history(tx, doc_name, doc_data, tool.Get_time(), ip, "", "0", "upload", "")
		return nil
	})
	if err != nil {
		return false
	}
	tool.Search_index_update(doc_name, doc_data)
	markup.Get_render(db, doc_name, doc_data, "backlink")
	return true
}

func api_file_upload_post(config tool.Config, file_name string, file_data []byte, file_ext string, license string, license_text string, captcha string, check_captcha bool, many_upload bool, replace bool) map[string]any {
	return api_file_upload_post_reader(config, file_name, bytes.NewReader(file_data), file_ext, license, license_text, captcha, check_captcha, many_upload, replace)
}

func api_file_upload_post_reader(config tool.Config, file_name string, file_reader io.Reader, file_ext string, license string, license_text string, captcha string, check_captcha bool, many_upload bool, replace bool) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	file_name = strings.TrimSpace(file_name)
	file_ext = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(file_ext)), ".")

	allowed_ext := tool.Get_ext_allow_list(db)
	return_value := make(map[string]any)

	if file_reader == nil || file_name == "" || file_ext == "" {
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

	file_max_size := tool.Get_file_max_size_by_extension(db, file_ext)
	if file_max_size <= 0 {
		file_max_size = tool.Get_file_max_size(db)
	}
	if file_max_size <= 0 {
		file_max_size = 2
	}
	file_max_bytes := int64(file_max_size) * 1000 * 1000

	doc_name := "file:" + file_name + "." + file_ext
	var old_doc_name string
	old_doc_exists := tool.QueryRow_DB(db, "select title from data where title = ?", []any{&old_doc_name}, doc_name)
	if old_doc_exists && !replace {
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
	if _, err := os.Stat(dst_path); err == nil && !replace {
		return_value["response"] = "error"
		return_value["data"] = "already exist"
		return return_value
	} else if !errors.Is(err, os.ErrNotExist) {
		return_value["response"] = "error"
		return_value["data"] = "exist check fail"
		return return_value
	}

	temp_file, err := os.CreateTemp(main_dir, ".opennamu-upload-*")
	if err != nil {
		return_value["response"] = "error"
		return_value["data"] = "file create fail"
		return return_value
	}
	temp_path := temp_file.Name()
	temp_closed := false
	temp_removed := true
	defer func() {
		if !temp_closed {
			_ = temp_file.Close()
		}
		if temp_removed {
			_ = os.Remove(temp_path)
		}
	}()

	written, err := io.Copy(temp_file, io.LimitReader(file_reader, file_max_bytes+1))
	if err != nil {
		return_value["response"] = "error"
		return_value["data"] = "file write fail"
		return return_value
	}
	if written == 0 {
		return_value["response"] = "error"
		return_value["data"] = "invalid data"
		return return_value
	}
	if written > file_max_bytes {
		return_value["response"] = "error"
		return_value["data"] = "file too large"
		return return_value
	}
	if err := temp_file.Chmod(0o644); err != nil {
		return_value["response"] = "error"
		return_value["data"] = "file write fail"
		return return_value
	}
	if err := temp_file.Close(); err != nil {
		temp_closed = true
		return_value["response"] = "error"
		return_value["data"] = "file write fail"
		return return_value
	}
	temp_closed = true

	rename_err := os.Rename(temp_path, dst_path)
	if rename_err != nil && replace && errors.Is(rename_err, os.ErrExist) {
		if err := os.Remove(dst_path); err == nil {
			rename_err = os.Rename(temp_path, dst_path)
		}
	}
	if rename_err != nil {
		return_value["response"] = "error"
		return_value["data"] = "file write fail"
		return return_value
	}
	temp_removed = false

	if old_doc_exists {
		return_value["response"] = "ok"
		return_value["data"] = doc_name
		return return_value
	}

	if license == "" {
		license = "direct_input"
	}
	license_is_registered := license != "direct_input" && tool.Get_html_filter_value(db, license, "image_license")[0] != ""
	license = strings.ReplaceAll(license, "]", "_")
	doc_data := license + "\n"
	if tool.Get_document_markup(db, doc_name, "document") == "namumark" {
		if license_is_registered {
			doc_data += "[include(틀:" + license + ")]\n"
		}
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
