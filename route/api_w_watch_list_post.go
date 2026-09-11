package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func Api_w_watch_list_post(config tool.Config, name string, do_type string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if tool.IP_or_user(config.IP) {
		return map[string]any{"response": "require auth"}
	}
	if do_type != "watchlist" && do_type != "star_doc" && do_type != "thread_watchlist" && do_type != "bbs_watchlist" {
		do_type = "star_doc"
	}

	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		var data string
		exist := tool.QueryRow_DB(tx, "select data from user_set where name = ? and id = ? and data = ?", []any{&data}, do_type, config.IP, name)
		if exist {
			_, err := tx.Exec(tool.DB_change("delete from user_set where name = ? and id = ? and data = ?"), do_type, config.IP, name)
			return err
		}
		_, err := tx.Exec(tool.DB_change("insert into user_set (id, name, data) values (?, ?, ?)"), config.IP, do_type, name)
		return err
	}); err != nil {
		panic(err)
	}

	return map[string]any{"response": "ok"}
}
