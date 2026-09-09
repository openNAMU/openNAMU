package tool

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
)

func Get_user_profile_image_name(db *sql.DB, value string) (string, bool) {
	value = strings.TrimSpace(value)
	if strings.HasPrefix(strings.ToLower(value), "file:") {
		value = value[5:]
	}
	if value == "" {
		return "", true
	}
	if strings.ContainsAny(value, `/\`) {
		return "", false
	}

	extension := strings.TrimPrefix(strings.ToLower(filepath.Ext(value)), ".")
	base_name := strings.TrimSuffix(value, filepath.Ext(value))
	if base_name == "" || extension == "" || strings.Contains(base_name, ".") || !profile_image_extension(extension) {
		return "", false
	}

	doc_name := "file:" + value
	if _, exists := Get_data_title(db, doc_name); !exists {
		return "", false
	}

	file_path := filepath.Join(Get_file_main_dir(db), File_name_to_dir(base_name, extension))
	file_info, err := os.Stat(file_path)
	if err != nil || file_info.IsDir() {
		return "", false
	}

	return value, true
}

func Get_user_profile_image_ui(db *sql.DB, user_name string) string {
	if IP_or_user(user_name) {
		return ""
	}

	file_name, valid := Get_user_profile_image_name(db, Get_user_set_data(db, user_name, "profile_image"))
	if !valid || file_name == "" {
		return ""
	}

	extension := strings.TrimPrefix(strings.ToLower(filepath.Ext(file_name)), ".")
	base_name := strings.TrimSuffix(file_name, filepath.Ext(file_name))
	storage_name := File_name_to_dir(base_name, extension)
	revision := Get_history_last_revision(db, "file:"+file_name)
	if revision == "" {
		revision = "1"
	}

	return `<img class="opennamu_user_profile_image" loading="lazy" width="32" height="32" style="width:32px;height:32px;object-fit:cover;vertical-align:middle;" src="/thumbnail/64/` + Url_parser(storage_name) + `.cache_v` + Url_parser(revision) + `" alt=""> `
}

func profile_image_extension(extension string) bool {
	switch extension {
	case "jpg", "jpeg", "png", "gif", "webp", "bmp", "tif", "tiff", "avif", "heic":
		return true
	default:
		return false
	}
}
