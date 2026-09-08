package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func bbs_watch_key(set_id string, set_code string) string {
	return set_id + "-" + set_code
}

func Api_bbs_watch_post(config tool.Config, set_id string, set_code string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	title := ""
	if !tool.QueryRow_DB(
		db,
		"select set_data from bbs_data where set_name = 'title' and set_id = ? and set_code = ?",
		[]any{&title},
		set_id,
		set_code,
	) {
		return map[string]any{"response": "not exist", "data": "post"}
	}
	if tool.IP_or_user(config.IP) {
		return map[string]any{"response": "require auth"}
	}

	return Api_w_watch_list_post(config, bbs_watch_key(set_id, set_code), "bbs_watchlist")
}

func bbs_watch_notify(db *sql.DB, config tool.Config, set_id string, set_code string, comment_code string, bbs_name string, title string, post_user string, parent_user string) {
	alarm := `BBS <a href="/bbs/w/` + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code) + "#" + tool.Url_parser(comment_code) + `">` + tool.HTML_escape(bbs_name) + " - " + tool.HTML_escape(title) + "#" + tool.Url_parser(comment_code) + `</a>`
	skip := map[string]bool{
		post_user:   true,
		parent_user: true,
	}
	sent := map[string]bool{}

	rows := tool.Query_DB(
		db,
		"select id from user_set where name = 'bbs_watchlist' and data = ?",
		bbs_watch_key(set_id, set_code),
	)
	defer rows.Close()

	for rows.Next() {
		target := ""
		if rows.Scan(&target) != nil || target == "" || skip[target] || sent[target] {
			continue
		}
		sent[target] = true
		tool.Send_alarm(db, config.IP, target, alarm)
	}
}
