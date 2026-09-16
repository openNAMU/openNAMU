package route

import "opennamu/route/tool"

func Api_bbs_watch_view(config tool.Config, set_id string, set_code string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	bbs_name := ""
	title := ""
	if !tool.QueryRow_DB(
		db,
		"select set_data from bbs_set where set_name = 'bbs_name' and set_id = ?",
		[]any{&bbs_name},
		set_id,
	) || !tool.QueryRow_DB(
		db,
		"select set_data from bbs_data where set_name = 'title' and set_id = ? and set_code = ?",
		[]any{&title},
		set_id,
		set_code,
	) {
		return map[string]any{"response": "not exist"}
	}
	if _, allowed := bbs_post_view_auth(db, set_id, set_code, config.IP); !allowed {
		return map[string]any{"response": "require auth"}
	}

	return map[string]any{
		"response": "ok",
		"data": map[string]string{
			"bbs_name": bbs_name,
			"title":    title,
		},
	}
}
