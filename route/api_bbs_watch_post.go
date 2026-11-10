package route

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"log"
	"runtime"
	"strings"

	"github.com/go-sql-driver/mysql"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
	turso "turso.tech/database/tursogo"

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
		separator := strings.LastIndex(key, "-")
		if separator <= 0 || separator == len(key)-1 {
			return map[string]any{"response": "not exist", "data": "post"}
		}

		err := Bbs_board_write_transaction(db, key[:separator], key[separator+1:], func(tx *sql.Tx) error {
			return User_watch_list_toggle(tx, config.IP, "bbs_watchlist", key)
		})
		if err != nil {
			if err == bbs_board_missing_error {
				return map[string]any{"response": "not exist", "data": "bbs"}
			}
			panic(err)
		}
		return map[string]any{"response": "ok"}
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

func Bbs_notification_db_call(callback func()) (return_err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			recovered_err, is_error := recovered.(error)
			var runtime_error runtime.Error
			if !is_error || errors.As(recovered_err, &runtime_error) || !Bbs_notification_database_error(recovered_err) {
				panic(recovered)
			}
			return_err = recovered_err
		}
	}()

	callback()
	return nil
}

func Bbs_notification_database_error(err error) bool {
	var mysql_error *mysql.MySQLError
	if errors.As(err, &mysql_error) || errors.Is(err, driver.ErrBadConn) {
		return true
	}
	var sqlite_error *sqlite.Error
	if errors.As(err, &sqlite_error) {
		switch sqlite_error.Code() & 0xff {
		case sqlite3.SQLITE_BUSY, sqlite3.SQLITE_LOCKED, sqlite3.SQLITE_IOERR,
			sqlite3.SQLITE_NOMEM, sqlite3.SQLITE_READONLY, sqlite3.SQLITE_INTERRUPT,
			sqlite3.SQLITE_FULL, sqlite3.SQLITE_CANTOPEN, sqlite3.SQLITE_PROTOCOL,
			sqlite3.SQLITE_NOTADB, sqlite3.SQLITE_CORRUPT, sqlite3.SQLITE_CONSTRAINT,
			sqlite3.SQLITE_AUTH, sqlite3.SQLITE_PERM:
			return true
		}
		return false
	}

	for _, turso_error := range []error{
		turso.ErrTursoBusy,
		turso.ErrTursoInterrupt,
		turso.ErrTursoGeneric,
		turso.ErrTursoConstraint,
		turso.ErrTursoReadOnly,
		turso.ErrTursoDatabaseFull,
		turso.ErrTursoNotADb,
		turso.ErrTursoCorrupt,
	} {
		if errors.Is(err, turso_error) {
			return true
		}
	}
	return false
}

func Bbs_watch_notify(db *sql.DB, config tool.Config, set_id string, set_code string, comment_code string, bbs_name string, title string, post_user string, parent_user string) {
	key := Bbs_watch_key(set_id, set_code)
	targets := []string{}
	sent := map[string]bool{}
	skip := map[string]bool{
		post_user:   true,
		parent_user: true,
	}

	append_watchers := func(query string, values ...any) error {
		rows, err := db.Query(tool.DB_change(query), values...)
		if err != nil {
			return err
		}
		for rows.Next() {
			target := ""
			if err := rows.Scan(&target); err != nil {
				log.Printf("[BBS WATCH] watcher scan failed: %v", err)
				continue
			}
			if target == "" || skip[target] || sent[target] {
				continue
			}
			sent[target] = true
			targets = append(targets, target)
		}
		rows_err := rows.Err()
		close_err := rows.Close()
		if rows_err != nil {
			log.Printf("[BBS WATCH] watcher lookup iteration failed: %v", rows_err)
		}
		if close_err != nil {
			log.Printf("[BBS WATCH] watcher lookup cursor close failed: %v", close_err)
		}
		return nil
	}

	if err := append_watchers("select id from user_set where name = ? and data = ?", "bbs_watchlist", key); err != nil {
		log.Printf("[BBS WATCH] existing watcher lookup failed: %v", err)
		return
	}
	if set_id == thread_bbs_id {
		if err := append_watchers("select id from user_set where name = ? and (data = ? or data = ?)", "thread_watchlist", set_code, key); err != nil {
			log.Printf("[BBS WATCH] legacy thread watcher lookup failed: %v", err)
		}
	}

	for _, target := range targets {
		var allowed bool
		if err := Bbs_notification_db_call(func() {
			_, allowed = Bbs_post_view_auth(db, set_id, set_code, target)
		}); err != nil {
			log.Printf("[BBS WATCH] visibility lookup failed for recipient %q: %v", target, err)
			continue
		}
		if !allowed {
			continue
		}

		alarm := `BBS <a href="/bbs/w/` + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code) + "#" + tool.Url_parser(comment_code) + `">` + tool.HTML_escape(bbs_name) + " - " + tool.HTML_escape(title) + "#" + tool.Url_parser(comment_code) + `</a>`
		if err := Bbs_notification_db_call(func() {
			tool.Send_alarm(db, config.IP, target, alarm)
		}); err != nil {
			log.Printf("[BBS WATCH] notification delivery failed for recipient %q: %v", target, err)
		}
	}
}
