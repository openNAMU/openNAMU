package route

import (
	"database/sql"
	"errors"

	"opennamu/route/tool"
)

func Api_bbs_delete(config tool.Config, set_id string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)

	bbs_name := ""
	if !tool.QueryRow_DB(
		db,
		"select set_data from bbs_set where set_id = ? and set_name = 'bbs_name'",
		[]any{&bbs_name},
		set_id,
	) {
		return_data["response"] = "not exist"
		return_data["data"] = "bbs"

		return return_data
	}

	if Bbs_is_special_board(set_id) {
		return_data["response"] = "error"
		return_data["data"] = "not allowed"

		return return_data
	}

	if !tool.Check_permission(db, "bbs_delete", config.IP) {
		return_data["response"] = "require auth"

		return return_data
	}

	if err := Bbs_board_write_transaction(db, set_id, "", func(tx *sql.Tx) error {
		for _, value := range []struct {
			query string
			args  []any
		}{
			{
				"delete from user_set where name = 'bbs_watchlist' and data like ?",
				[]any{set_id + "-%"},
			},
			{
				"delete from user_set where name = 'bbs_bookmark' and data like ?",
				[]any{set_id + "-%"},
			},
			{
				"delete from user_set where name = 'bbs_board_watchlist' and data = ?",
				[]any{set_id},
			},
			{
				"delete from bbs_data where set_id = ?",
				[]any{set_id},
			},
			{
				"delete from bbs_set where set_id = ?",
				[]any{set_id},
			},
			{
				"delete from bbs_data where set_id like ?",
				[]any{set_id + "-%"},
			},
		} {
			if _, err := tx.Exec(tool.DB_change(value.query), value.args...); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		if errors.Is(err, bbs_board_missing_error) {
			return_data["response"] = "not exist"
			return_data["data"] = "bbs"
			return return_data
		}
		panic(err)
	}
	tool.Search_bbs_index_delete_set(db, set_id)

	return_data["response"] = "ok"

	return return_data
}
