package route

import (
	"opennamu/route/tool"
	"strconv"
)

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
	if _, allowed := Bbs_post_view_auth(db, set_id, set_code, config.IP); !allowed {
		return map[string]any{"response": "require auth"}
	}

	watching := false
	if !tool.IP_or_user(config.IP) {
		key := Bbs_watch_key(set_id, set_code)
		watching = tool.Get_user_set_data_match(db, config.IP, "bbs_watchlist", key)
		if set_id == thread_bbs_id {
			watching = watching || tool.Get_user_set_data_match(db, config.IP, "thread_watchlist", set_code) || tool.Get_user_set_data_match(db, config.IP, "thread_watchlist", key)
		}
	}
	return map[string]any{
		"response": "ok",
		"data": map[string]string{
			"bbs_name": bbs_name,
			"title":    title,
			"watched":  strconv.FormatBool(watching),
		},
	}
}
