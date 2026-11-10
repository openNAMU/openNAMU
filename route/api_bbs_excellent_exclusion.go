package route

import (
	"database/sql"
	"errors"
	"strconv"

	"opennamu/route/tool"
)

var bbs_excellent_post_missing_error = errors.New("bbs excellent post does not exist")

func Api_bbs_excellent_exclusion_set(config tool.Config, set_id string, set_code string, excluded bool) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "bbs_manage", config.IP) {
		return map[string]any{"response": "require auth"}
	}

	board_id, board_err := strconv.ParseInt(set_id, 10, 64)
	post_code, post_err := strconv.ParseInt(set_code, 10, 64)
	if board_err != nil || strconv.FormatInt(board_id, 10) != set_id {
		return map[string]any{"response": "not exist", "data": "bbs"}
	}
	if post_err != nil || post_code < 1 || strconv.FormatInt(post_code, 10) != set_code {
		return map[string]any{"response": "not exist", "data": "post"}
	}

	if err := Bbs_board_write_transaction(db, set_id, set_code, func(tx *sql.Tx) error {
		var board_exists int
		if !tool.QueryRow_DB(
			tx,
			"select 1 from bbs_set where set_id = ? and set_name = 'bbs_name'",
			[]any{&board_exists},
			set_id,
		) {
			return bbs_board_missing_error
		}

		var post_exists int
		if !tool.QueryRow_DB(
			tx,
			"select 1 from bbs_data where set_name = 'title' and set_id = ? and set_code = ?",
			[]any{&post_exists},
			set_id,
			set_code,
		) {
			return bbs_excellent_post_missing_error
		}

		if _, err := tx.Exec(
			tool.DB_change("delete from bbs_data where set_name = 'excellent_excluded' and set_id = ? and set_code = ?"),
			set_id,
			set_code,
		); err != nil {
			return err
		}
		if excluded {
			if _, err := tx.Exec(
				tool.DB_change("insert into bbs_data (set_name, set_code, set_id, set_data) values ('excellent_excluded', ?, ?, '1')"),
				set_code,
				set_id,
			); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		if errors.Is(err, bbs_board_missing_error) {
			return map[string]any{"response": "not exist", "data": "bbs"}
		}
		if errors.Is(err, bbs_excellent_post_missing_error) {
			return map[string]any{"response": "not exist", "data": "post"}
		}
		return map[string]any{"response": "error", "data": "database"}
	}

	state := "restored"
	if excluded {
		state = "excluded"
	}
	return map[string]any{"response": "ok", "data": state}
}
