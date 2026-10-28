package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func Bbs_watch_key(set_id string, set_code string) string {
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
	if _, allowed := Bbs_post_view_auth(db, set_id, set_code, config.IP); !allowed {
		return map[string]any{"response": "require auth"}
	}

	key := Bbs_watch_key(set_id, set_code)
	if set_id != thread_bbs_id {
		return Api_w_watch_list_post(config, key, "bbs_watchlist")
	}

	err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		var bbs_count, thread_raw_count, thread_key_count int
		if err := tx.QueryRow("select count(*) from user_set where name = 'bbs_watchlist' and id = ? and data = ?", config.IP, key).Scan(&bbs_count); err != nil {
			return err
		}
		if err := tx.QueryRow("select count(*) from user_set where name = 'thread_watchlist' and id = ? and data = ?", config.IP, set_code).Scan(&thread_raw_count); err != nil {
			return err
		}
		if err := tx.QueryRow("select count(*) from user_set where name = 'thread_watchlist' and id = ? and data = ?", config.IP, key).Scan(&thread_key_count); err != nil {
			return err
		}

		watched := bbs_count > 0 || thread_raw_count > 0 || thread_key_count > 0
		if watched {
			if _, err := tx.Exec(tool.DB_change("delete from user_set where name = ? and id = ? and data = ?"), "bbs_watchlist", config.IP, key); err != nil {
				return err
			}
			if _, err := tx.Exec(tool.DB_change("delete from user_set where name = ? and id = ? and (data = ? or data = ?)"), "thread_watchlist", config.IP, set_code, key); err != nil {
				return err
			}
			return nil
		}

		if _, err := tx.Exec(tool.DB_change("insert into user_set (id, name, data) values (?, ?, ?)"), config.IP, "bbs_watchlist", key); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return map[string]any{"response": "ok"}
}

func Bbs_watch_notify(db *sql.DB, config tool.Config, set_id string, set_code string, comment_code string, bbs_name string, title string, post_user string, parent_user string) {
	alarm := `BBS <a href="/bbs/w/` + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code) + "#" + tool.Url_parser(comment_code) + `">` + tool.HTML_escape(bbs_name) + " - " + tool.HTML_escape(title) + "#" + tool.Url_parser(comment_code) + `</a>`
	skip := map[string]bool{
		post_user:   true,
		parent_user: true,
	}
	sent := map[string]bool{}

	key := Bbs_watch_key(set_id, set_code)

	rows := tool.Query_DB(
		db,
		"select id from user_set where name = 'bbs_watchlist' and data = ?",
		key,
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

	if set_id == thread_bbs_id {
		thread_rows := tool.Query_DB(
			db,
			"select id from user_set where name = 'thread_watchlist' and (data = ? or data = ?)",
			set_code,
			key,
		)
		defer thread_rows.Close()

		for thread_rows.Next() {
			target := ""
			if thread_rows.Scan(&target) != nil || target == "" || skip[target] || sent[target] {
				continue
			}
			sent[target] = true
			tool.Send_alarm(db, config.IP, target, alarm)
		}
	}
}
