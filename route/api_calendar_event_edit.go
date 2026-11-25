package route

import (
	"database/sql"
	"errors"

	"opennamu/route/tool"
)

func Api_calendar_event_edit(config tool.Config, id string, code string, date string, title string, memo string) map[string]any {
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

	date, date_ok := tool.Calendar_date_normalize(Calendar_text_normalize(date))
	if !date_ok {
		return_data["response"] = "error"
		return_data["data"] = "invalid date"
		return return_data
	}

	title = Calendar_text_normalize(title)
	memo = Calendar_text_normalize(memo)
	if title == "" {
		return_data["response"] = "error"
		return_data["data"] = "empty data"
		return return_data
	}
	if !tool.Do_edit_filter(db, config, "", title) || !tool.Do_edit_filter(db, config, "", memo) {
		return_data["response"] = "error"
		return_data["data"] = "edit filter (content)"
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
		if _, exists := Calendar_event_get(tx, id, code); !exists {
			return calendar_event_missing_error
		}

		for _, value := range [][]string{
			{"date", date},
			{"title", title},
			{"memo", memo},
			{"updated", tool.Get_time()},
		} {
			if _, err := tx.Exec(
				tool.DB_change("update calendar set set_data = ? where set_id = ? and set_code = ? and set_name = ?"),
				value[1],
				id,
				code,
				value[0],
			); err != nil {
				return err
			}
		}

		return nil
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
	return_data["month"] = date[:7]
	return return_data
}
