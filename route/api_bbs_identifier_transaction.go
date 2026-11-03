package route

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"database/sql/driver"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"

	"opennamu/route/tool"
)

const bbs_identifier_lock_timeout_seconds = 30

// Bbs_identifier_lock_name returns a bounded MariaDB advisory-lock name for a
// database and logical BBS write scope. MariaDB limits lock names to 64 bytes.
func Bbs_identifier_lock_name(database_name string, scope string) string {
	digest := sha256.Sum256([]byte(database_name + "\x00" + scope))
	return "onamu:bbs:" + hex.EncodeToString(digest[:27])
}

// Bbs_identifier_transaction serializes BBS identifier allocation and the
// associated write until commit. Top-level posts use a board lock; comment
// writes use a post lock and may also request the board lock when board Slow
// Mode is enabled. SQLite/Turso use existing EAV/config rows as transaction
// write anchors; MySQL/MariaDB use connection-owned named locks.
func Bbs_identifier_transaction(db *sql.DB, board_id string, post_code string, lock_board bool, callback func(*sql.Tx) error) error {
	lock_scopes := []string{}
	if lock_board {
		lock_scopes = append(lock_scopes, "board:"+board_id)
	}
	if post_code != "" {
		lock_scopes = append(lock_scopes, "post:"+board_id+":"+post_code)
	}
	if len(lock_scopes) == 0 {
		return tool.DB_transaction_retry(db, callback)
	}
	if tool.Get_DB_type() == "mysql" {
		return Bbs_mysql_identifier_transaction(db, lock_scopes, callback)
	}
	return tool.DB_transaction_retry(db, func(tx *sql.Tx) error {
		board_locked := false
		if lock_board {
			if err := Bbs_local_board_write_lock(tx, board_id); err != nil {
				return err
			}
			board_locked = true
		}
		if post_code != "" {
			if err := Bbs_local_post_write_lock(tx, board_id, post_code, board_locked); err != nil {
				return err
			}
		}
		return callback(tx)
	})
}

func Bbs_local_board_write_lock(tx *sql.Tx, board_id string) error {
	// Temporarily change and restore one character so drivers that optimize a
	// no-op UPDATE still establish a write lock for this existing row.
	result, err := tx.Exec(
		tool.DB_change("update bbs_set set set_data = set_data || char(1) where set_name = 'bbs_name' and set_id = ?"),
		board_id,
	)
	if err != nil {
		return err
	}
	rows_affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows_affected > 0 {
		_, err := tx.Exec(
			tool.DB_change("update bbs_set set set_data = substr(set_data, 1, length(set_data) - 1) where set_name = 'bbs_name' and set_id = ?"),
			board_id,
		)
		return err
	}

	// Special internal boards such as set_id 0 and the report board may lack
	// a bbs_name row. `other.ver` is initialized during application startup and
	// gives only those scopes a stable anchor without adding EAV state.
	if !Bbs_is_special_board(board_id) {
		return errors.New("bbs identifier lock anchor is missing")
	}
	result, err = tx.Exec(tool.DB_change("update other set data = data || char(1) where name = 'ver'"))
	if err != nil {
		return err
	}
	rows_affected, err = result.RowsAffected()
	if err != nil {
		return err
	}
	if rows_affected == 0 {
		return errors.New("bbs identifier lock anchor is missing")
	}
	_, err = tx.Exec(tool.DB_change("update other set data = substr(data, 1, length(data) - 1) where name = 'ver'"))
	return err
}

func Bbs_local_post_write_lock(tx *sql.Tx, board_id string, post_code string, board_locked bool) error {
	result, err := tx.Exec(
		tool.DB_change("update bbs_data set set_data = set_data || char(1) where set_name = 'comment_count' and set_id = ? and set_code = ?"),
		board_id,
		post_code,
	)
	if err != nil {
		return err
	}
	rows_affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows_affected > 0 {
		_, err := tx.Exec(
			tool.DB_change("update bbs_data set set_data = substr(set_data, 1, length(set_data) - 1) where set_name = 'comment_count' and set_id = ? and set_code = ?"),
			board_id,
			post_code,
		)
		return err
	}
	if board_locked {
		return nil
	}
	// A legacy post may lack comment_count. Fall back to its board anchor so
	// it still receives a stable lock without creating persistent metadata.
	return Bbs_local_board_write_lock(tx, board_id)
}

func Bbs_mysql_identifier_transaction(db *sql.DB, lock_scopes []string, callback func(*sql.Tx) error) (return_err error) {
	conn, err := db.Conn(context.Background())
	if err != nil {
		return err
	}
	committed := false
	acquired_locks := []string{}
	physical_connection_bad := false
	defer func() {
		if release_err := Bbs_mysql_identifier_locks_release(conn, acquired_locks); release_err != nil {
			physical_connection_bad = true
			if return_err == nil && !committed {
				return_err = release_err
			}
		}
		if physical_connection_bad {
			_ = conn.Raw(func(any) error { return driver.ErrBadConn })
		}
		if close_err := conn.Close(); close_err != nil && return_err == nil && !committed {
			return_err = close_err
		}
	}()

	var database_name string
	if err := conn.QueryRowContext(context.Background(), "select database()").Scan(&database_name); err != nil {
		return err
	}
	if database_name == "" {
		return errors.New("bbs identifier lock requires a selected database")
	}

	lock_names := make([]string, 0, len(lock_scopes))
	for _, scope := range lock_scopes {
		lock_names = append(lock_names, Bbs_identifier_lock_name(database_name, scope))
	}
	sort.Strings(lock_names)
	for index := 1; index < len(lock_names); index++ {
		if lock_names[index] == lock_names[index-1] {
			lock_names = append(lock_names[:index], lock_names[index+1:]...)
			index--
		}
	}

	for _, lock_name := range lock_names {
		var acquired sql.NullInt64
		if err := conn.QueryRowContext(
			context.Background(),
			"select get_lock(?, ?)",
			lock_name,
			bbs_identifier_lock_timeout_seconds,
		).Scan(&acquired); err != nil {
			// The server may have acquired the lock before the response failed.
			// Discard this session rather than returning a potentially locked
			// connection to database/sql's pool.
			physical_connection_bad = true
			return err
		}
		if !acquired.Valid {
			physical_connection_bad = true
			return fmt.Errorf("bbs identifier lock acquisition failed")
		}
		if acquired.Int64 != 1 {
			return fmt.Errorf("bbs identifier lock acquisition timed out")
		}
		acquired_locks = append(acquired_locks, lock_name)
	}

	if err := tool.DB_transaction_retry_conn(conn, callback); err != nil {
		return err
	}
	committed = true
	return nil
}

func Bbs_mysql_identifier_locks_release(conn *sql.Conn, lock_names []string) error {
	var return_err error
	for index := len(lock_names) - 1; index >= 0; index-- {
		var released sql.NullInt64
		err := conn.QueryRowContext(context.Background(), "select release_lock(?)", lock_names[index]).Scan(&released)
		if err != nil {
			if return_err == nil {
				return_err = err
			}
			continue
		}
		if !released.Valid || released.Int64 != 1 {
			if return_err == nil {
				return_err = errors.New("bbs identifier advisory lock was not released")
			}
		}
	}
	return return_err
}
