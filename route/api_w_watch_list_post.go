package route

import (
	"database/sql"
	"strings"

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

	toggle := func(tx *sql.Tx) error {
		return User_watch_list_toggle(tx, config.IP, do_type, name)
	}

	var err error
	if do_type == "bbs_watchlist" {
		separator := strings.LastIndex(name, "-")
		if separator <= 0 || separator == len(name)-1 {
			return map[string]any{"response": "not exist", "data": "post"}
		}
		err = Bbs_board_write_transaction(db, name[:separator], name[separator+1:], toggle)
	} else {
		err = tool.DB_transaction(db, toggle)
	}
	if err != nil {
		if do_type == "bbs_watchlist" && err == bbs_board_missing_error {
			return map[string]any{"response": "not exist", "data": "bbs"}
		}
		panic(err)
	}

	return map[string]any{"response": "ok"}
}

func User_watch_list_toggle(tx *sql.Tx, user_id string, watch_type string, name string) error {
	var existing_data string
	if tool.QueryRow_DB(
		tx,
		"select data from user_set where name = ? and id = ? and data = ?",
		[]any{&existing_data},
		watch_type,
		user_id,
		name,
	) {
		_, err := tx.Exec(
			tool.DB_change("delete from user_set where name = ? and id = ? and data = ?"),
			watch_type,
			user_id,
			name,
		)
		return err
	}

	_, err := tx.Exec(
		tool.DB_change("insert into user_set (id, name, data) values (?, ?, ?)"),
		user_id,
		watch_type,
		name,
	)
	return err
}
