package route

import (
	"strings"

	"opennamu/route/tool"
)

func Api_history_add_post(config tool.Config, doc_name string, content string, get_ip string, send string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if doc_name == "" {
		return_data["response"] = "error"
		return_data["data"] = "invalid document"
		return return_data
	}

	content = strings.ReplaceAll(content, "\r", "")
	tool.Do_add_history(
		db,
		doc_name,
		content,
		tool.Get_time(),
		"Add:"+get_ip,
		send,
		"+"+tool.Get_edit_length_diff("", content),
		"",
		"add",
	)

	return_data["response"] = "ok"
	return return_data
}
