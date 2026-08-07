package route

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"opennamu/route/tool"
)

func Api_file_upload_make_document(db *sql.DB) {

}

func Api_file_upload_post(config tool.Config, file_name string, file_data string, file_ext string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	file_name = strings.TrimSpace(file_name)
	file_data = strings.TrimSpace(file_data)
	file_ext = strings.TrimSpace(file_ext)
	file_ext = strings.TrimPrefix(strings.ToLower(file_ext), ".")

	allowed_ext := tool.Get_ext_allow_list(db)
	// file_max_size := tool.Get_file_max_size(db)

	return_value := make(map[string]any)

	if file_data == "" || file_name == "" || file_ext == "" {
		return_value["response"] = "error"
		return_value["data"] = "invalid data"

		return return_value
	} else if _, ok := allowed_ext[file_ext]; !ok {
		return_value["response"] = "error"
		return_value["data"] = "unallowed ext"

		return return_value
	} else if tool.Get_file_name_unallow_check(db, file_name) {
		return_value["response"] = "error"
		return_value["data"] = "unallowed file name"

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

	out, err := os.Create(dst_path)
	if err != nil {
		return_value["response"] = "error"
		return_value["data"] = "file create fail"

		return return_value
	}
	defer out.Close()

	if _, err := io.Copy(out, base64.NewDecoder(base64.StdEncoding, strings.NewReader(file_data))); err != nil {
		_ = out.Close()
		_ = os.Remove(dst_path)

		return_value["response"] = "error"
		return_value["data"] = "file write fail"

		return return_value
	}

	return_value["response"] = "ok"

	return return_value
}
