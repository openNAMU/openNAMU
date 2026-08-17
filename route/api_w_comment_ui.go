package route

import "opennamu/route/tool"

func Api_w_comment_ui(config tool.Config, doc_name string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := Api_w_comment(config, doc_name)
	db_code_str, _ := return_data["data"].(string)

	return_data = make(map[string]any)
	return_data["response"] = "ok"
	return_data["data"] = View_bbs_in_w_comment(db, config, "0", db_code_str, "")

	return return_data
}
