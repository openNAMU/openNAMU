package route

import "opennamu/route/tool"

func Api_history_hidden_post(config tool.Config, doc_name string, rev string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "hidel", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	revision, ok := history_revision_value(rev)
	if !ok || doc_name == "" {
		return_data["response"] = "error"
		return_data["data"] = "invalid revision"
		return return_data
	}

	hide := ""
	tool.QueryRow_DB(db, "select hide from history where title = ? and id = ?", []any{&hide}, doc_name, revision)
	if hide == "" {
		hide = "O"
	} else {
		hide = ""
	}
	tool.Exec_DB(db, "update history set hide = ? where title = ? and id = ?", hide, doc_name, revision)

	return_data["response"] = "ok"
	return_data["data"] = hide
	return return_data
}
