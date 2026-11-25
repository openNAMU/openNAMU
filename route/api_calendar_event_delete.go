package route

import (
	"database/sql"
	"errors"

	"opennamu/route/tool"
)

func Api_calendar_event_delete(config tool.Config, id string, code string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !Calendar_view_allowed(db, config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	id, id_ok := Calendar_id_normalize(id)
	code, code_ok := Calendar_code_normalize(code)
	if !id_ok || !code_ok {
		return_data["response"] = "not exist"
		return_data["data"] = "event"
		return return_data
	}

	event_data, event_exists := Calendar_event_get(db, id, code)
	if !event_exists {
		return_data["response"] = "not exist"
		return_data["data"] = "event"
		return return_data
	}

	calendar_data, calendar_exists := Calendar_meta_get(db, id)
	if !calendar_exists {
		return_data["response"] = "not exist"
		return_data["data"] = "calendar"
		return return_data
	}
	if !Calendar_manage_allowed(db, config, calendar_data["owner"], event_data["user"]) {
		return_data["response"] = "require auth"
		return return_data
	}

	month := ""
	if len(event_data["date"]) >= 7 {
		month = event_data["date"][:7]
	}

	if err := tool.DB_transaction_retry(db, func(tx *sql.Tx) error {
		locked, lock_err := Calendar_write_lock(tx, id)
		if lock_err != nil {
			return lock_err
		}
		if !locked {
			return calendar_missing_error
		}
		if _, exists := Calendar_event_get(tx, id, code); !exists {
			return calendar_event_missing_error
		}

		_, err := tx.Exec(
			tool.DB_change("delete from calendar where set_id = ? and set_code = ?"),
			id,
			code,
		)
		return err
	}); err != nil {
		if errors.Is(err, calendar_missing_error) {
			return_data["response"] = "not exist"
			return_data["data"] = "calendar"
			return return_data
		}
		if errors.Is(err, calendar_event_missing_error) {
			return_data["response"] = "not exist"
			return_data["data"] = "event"
			return return_data
		}
		panic(err)
	}

	return_data["response"] = "ok"
	return_data["month"] = month
	return return_data
}
