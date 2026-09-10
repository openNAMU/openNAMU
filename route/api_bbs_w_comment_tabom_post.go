package route

import (
	"opennamu/route/tool"
)

func Api_bbs_w_comment_tabom_post(config tool.Config, set_id string, set_code string, comment_code string, vote_type string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_acl(db, set_id, "", "bbs_comment", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	comment_set_id, comment_set_code, exists := bbs_search_comment_location(set_id, set_code, comment_code)
	if !exists {
		return_data["response"] = "not exist"
		return_data["data"] = "comment"
		return return_data
	}

	comment := ""
	if !tool.QueryRow_DB(
		db,
		"select set_data from bbs_data where set_name = 'comment' and set_id = ? and set_code = ?",
		[]any{&comment},
		comment_set_id,
		comment_set_code,
	) || comment == "" {
		return_data["response"] = "not exist"
		return_data["data"] = "comment"
		return return_data
	}

	return api_bbs_tabom_post(db, config.IP, comment_set_id, comment_set_code, vote_type)
}
