package route

import (
	"database/sql"
	"errors"

	"opennamu/route/tool"
)

func Api_calendar_delete(config tool.Config, id string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !Calendar_view_allowed(db, config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	id, id_ok := Calendar_id_normalize(id)
	if !id_ok {
		return_data["response"] = "not exist"
		return_data["data"] = "calendar"
		return return_data
	}

	calendar_data, calendar_exists := Calendar_meta_get(db, id)
	if !calendar_exists {
		return_data["response"] = "not exist"
		return_data["data"] = "calendar"
		return return_data
	}
	if !Calendar_manage_allowed(db, config, calendar_data["owner"], "") {
		return_data["response"] = "require auth"
		return return_data
	}

	if err := tool.DB_transaction_retry(db, func(tx *sql.Tx) error {
		locked, lock_err := Calendar_write_lock(tx, id)
		if lock_err != nil {
			return lock_err
		}
		if !locked {
			return calendar_missing_error
		}

		_, err := tx.Exec(tool.DB_change("delete from calendar where set_id = ?"), id)
		return err
	}); err != nil {
		if errors.Is(err, calendar_missing_error) {
			return_data["response"] = "not exist"
			return_data["data"] = "calendar"
			return return_data
		}
		panic(err)
	}

	return_data["response"] = "ok"
	return return_data
}
