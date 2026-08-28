package route

import "opennamu/route/tool"

func Api_filter_delete_post(config tool.Config, kind string, name string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	spec, ok := get_filter_spec(kind)
	if !ok {
		return_data["response"] = "error"
		return return_data
	}
	if !tool.Check_permission(db, "filter_manage", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	tool.Exec_DB(db, "delete from html_filter where html = ? and kind = ?", name, spec.db_kind)
	if kind == "inter_wiki" {
		tool.Exec_DB(db, "delete from html_filter where html = ? and kind = 'inter_wiki_sub'", name)
	}
	tool.Do_insert_auth_history(db, config.IP, "filter_delete ("+kind+")")

	return_data["response"] = "ok"
	return return_data
}
