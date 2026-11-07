package route

import (
	"path/filepath"
	"strings"

	"opennamu/route/tool"
)

func File_rename_parts(doc_name string) (string, string, bool) {
	if !strings.HasPrefix(doc_name, "file:") {
		return "", "", false
	}

	file_name := strings.TrimPrefix(doc_name, "file:")
	extension := strings.TrimPrefix(strings.ToLower(filepath.Ext(file_name)), ".")
	base_name := strings.TrimSuffix(file_name, filepath.Ext(file_name))
	if base_name == "" || extension == "" || strings.ContainsAny(base_name, `/\\`) || strings.Contains(base_name, ".") {
		return "", "", false
	}

	return base_name, extension, true
}

func File_storage_path(db_file_dir string, doc_name string) (string, bool) {
	base_name, extension, valid := File_rename_parts(doc_name)
	if !valid {
		return "", false
	}
	return filepath.Join(db_file_dir, tool.File_name_to_dir(base_name, extension)), true
}
