package route

import (
	"database/sql"
	"errors"

	"opennamu/route/tool"
)

func Api_calendar_event_add(config tool.Config, id string, date string, title string, memo string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !Calendar_view_allowed(db, config.IP) || !tool.Check_permission(db, "calendar_create", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	id, id_ok := Calendar_id_normalize(id)
	if !id_ok {
		return_data["response"] = "not exist"
		return_data["data"] = "calendar"
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

	code := ""
	if err := tool.DB_transaction_retry(db, func(tx *sql.Tx) error {
		locked, lock_err := Calendar_write_lock(tx, id)
		if lock_err != nil {
			return lock_err
		}
		if !locked {
			return calendar_missing_error
		}

		code = Calendar_next_event_code(tx, id)
		now := tool.Get_time()

		for _, value := range [][]string{
			{"date", date},
			{"title", title},
			{"memo", memo},
			{"user", config.IP},
			{"created", now},
			{"updated", now},
		} {
			if _, err := tx.Exec(
				tool.DB_change("insert into calendar (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)"),
				value[0],
				code,
				id,
				value[1],
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
		panic(err)
	}

	return_data["response"] = "ok"
	return_data["data"] = code
	return_data["month"] = date[:7]
	return return_data
}
