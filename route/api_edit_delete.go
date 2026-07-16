package route

import (
	"opennamu/route/tool"
)

func Api_edit_delete(config tool.Config, doc_name string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_acl(db, doc_name, "", "document_delete", config.IP) {
		return_data["response"] = "require auth"

		return return_data
	}

	doc_title := ""
	exist := tool.QueryRow_DB(
		db,
		"select title from data where title = ?",
		[]any{&doc_title},
		doc_name,
	)
	if !exist {
		return_data["response"] = "not exist"

		return return_data
	}

	return_data["response"] = "ok"

	return return_data
}
