package route

import "opennamu/route/tool"

func Api_thread_comment_blind_post(config tool.Config, topic_num string, comment_num string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_acl(db, "", "", "toron_auth", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	block := ""
	if !tool.QueryRow_DB(db, "select block from topic where code = ? and id = ?", []any{&block}, topic_num, comment_num) {
		return_data["response"] = "not exist"
		return_data["data"] = "comment"
		return return_data
	}
	if block == "O" {
		block = ""
	} else {
		block = "O"
	}
	tool.Exec_DB(db, "update topic set block = ? where code = ? and id = ?", block, topic_num, comment_num)
	tool.Exec_DB(db, "update rd set date = ? where code = ?", tool.Get_time(), topic_num)
	tool.Do_insert_auth_history(db, config.IP, "blind (code "+topic_num+"#"+comment_num+")")

	return_data["response"] = "ok"
	return_data["data"] = block
	return return_data
}
