package route

import "opennamu/route/tool"

func Api_w_comment(config tool.Config, doc_name string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	db_code := tool.Get_document_setting(db, doc_name, "document_comment_code", "")

	db_code_str := ""
	if len(db_code) >= 1 {
		db_code_str = db_code[0][0]
	}

	if db_code_str == "" {
		db_code_str = Api_bbs_w_comment_make(config, doc_name)["data"].(string)

		tool.Exec_DB(
			db,
			"insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, '', 'document_comment_code', ?)",
			doc_name,
			db_code_str,
		)
	}

	return_data := make(map[string]any)
	return_data["response"] = "ok"
	return_data["data"] = db_code_str

	return return_data
}
