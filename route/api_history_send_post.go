package route

import "opennamu/route/tool"

func Api_history_send_post(config tool.Config, doc_name string, rev string, send string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	revision, ok := history_revision_value(rev)
	if !ok || doc_name == "" {
		return_data["response"] = "error"
		return_data["data"] = "invalid revision"
		return return_data
	}

	tool.Exec_DB(db, "update history set send = ? where title = ? and id = ?", send, doc_name, revision)

	return_data["response"] = "ok"
	return return_data
}
