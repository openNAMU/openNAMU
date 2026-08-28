package route

import (
	"net/url"
	"os"

	"opennamu/route/tool"
)

func Api_edit_file_delete_post(config tool.Config, doc_name string, values url.Values) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "file_delete", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	file_path, valid := file_storage_path(tool.Get_file_main_dir(db), doc_name)
	if !valid {
		return_data["response"] = "error"
		return_data["data"] = "invalid file"
		return return_data
	}

	if values.Get("with_doc") != "" {
		result := Api_edit_delete_post(config, doc_name, values.Get("send"), values.Get("copyright_agreement"))
		if result["response"] != "ok" {
			return_data["response"] = "error"
			return_data["data"] = "delete error"
			return return_data
		}
	}
	if err := os.Remove(file_path); err != nil && !os.IsNotExist(err) {
		return_data["response"] = "error"
		return_data["data"] = "file delete error"
		return return_data
	}
	tool.Do_insert_auth_history(db, config.IP, "file_delete ("+doc_name+")")

	return_data["response"] = "ok"
	return return_data
}
