package route

import "opennamu/route/tool"

func Api_thread_watch_post(config tool.Config, topic_num string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	title := ""
	if !tool.QueryRow_DB(db, "select title from rd where code = ?", []any{&title}, topic_num) {
		return map[string]any{"response": "not exist", "data": "thread"}
	}
	if tool.IP_or_user(config.IP) {
		return map[string]any{"response": "require auth"}
	}

	return Api_w_watch_list_post(config, topic_num, "thread_watchlist")
}
