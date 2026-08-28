package route

import (
	"strings"

	"opennamu/route/tool"
)

func Api_edit_delete_multiple_post(config tool.Config, content string, send string, agree string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)

	if !tool.Check_acl(db, "", "", "acl_auth", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	for _, doc_name := range strings.Split(strings.ReplaceAll(content, "\r", ""), "\n") {
		doc_name = strings.TrimSpace(doc_name)
		if doc_name == "" {
			continue
		}
		result := Api_edit_delete_post(config, doc_name, send, agree)
		response, _ := result["response"].(string)
		if response == "not exist" {
			continue
		}
		if response != "ok" {
			return result
		}
	}

	return_data["response"] = "ok"
	return return_data
}
