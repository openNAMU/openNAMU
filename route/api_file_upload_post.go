package route

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"io"
	"opennamu/route/tool"
	"os"
	"path/filepath"
	"strings"
)

func Api_file_upload_make_document(db *sql.DB) {
    
}

func Api_file_upload_post(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)
    
    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    file_name := strings.TrimSpace(other_set["file_name"])
    file_data := strings.TrimSpace(other_set["file_data"])
    file_ext := strings.TrimSpace(other_set["file_ext"])
    file_ext = strings.TrimPrefix(strings.ToLower(file_ext), ".")

    allowed_ext := tool.Get_ext_allow_list(db)
    // file_max_size := tool.Get_file_max_size(db)

    return_value := make(map[string]string)

    if file_data == "" || file_name == "" || file_ext == "" {
        return_value["response"] = "error"
        return_value["data"] = "invalid data"

        json_data, _ := json.Marshal(return_value)
        return string(json_data)
    } else if _, ok := allowed_ext[file_ext]; !ok {
        return_value["response"] = "error"
        return_value["data"] = "unallowed ext"
        
        json_data, _ := json.Marshal(return_value)
        return string(json_data)
    } else if tool.Get_file_name_unallow_check(db, file_name) {
        return_value["response"] = "error"
        return_value["data"] = "unallowed file name"
        
        json_data, _ := json.Marshal(return_value)
        return string(json_data)
    }

    main_dir := tool.Get_file_main_dir(db)

    if err := os.MkdirAll(main_dir, 0o755); err != nil {
        return_value["response"] = "error"
        return_value["data"] = "directory create fail"
        
        json_data, _ := json.Marshal(return_value)
        return string(json_data)
    }

    file_full_dir := tool.File_name_to_dir(file_name, file_ext)

    dst_path := filepath.Join(main_dir, file_full_dir)
    if _, err := os.Stat(dst_path); err == nil {
        return_value["response"] = "error"
        return_value["data"] = "already exist"

        json_data, _ := json.Marshal(return_value)
        return string(json_data)
    } else if !errors.Is(err, os.ErrNotExist) {
        return_value["response"] = "error"
        return_value["data"] = "exist check fail"

        json_data, _ := json.Marshal(return_value)
        return string(json_data)
    }

    out, err := os.Create(dst_path)
    if err != nil {
        return_value["response"] = "error"
        return_value["data"] = "file create fail"

        json_data, _ := json.Marshal(return_value)
        return string(json_data)
    }
    defer out.Close()

    if _, err := io.Copy(out, base64.NewDecoder(base64.StdEncoding, strings.NewReader(file_data))); err != nil {
        _ = out.Close()
        _ = os.Remove(dst_path)

        return_value["response"] = "error"
        return_value["data"] = "file write fail"

        json_data, _ := json.Marshal(return_value)
        return string(json_data)
    }

    return_value["response"] = "ok"
    
    json_data, _ := json.Marshal(return_value)
    return string(json_data)
}