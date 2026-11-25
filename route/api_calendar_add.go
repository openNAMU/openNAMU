package route

import (
	"database/sql"
	"strconv"

	"opennamu/route/tool"
)

func Api_calendar_add(config tool.Config, name string, subject string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !Calendar_view_allowed(db, config.IP) || !tool.Check_permission(db, "calendar_create", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	name = Calendar_text_normalize(name)
	subject = Calendar_text_normalize(subject)
	if name == "" {
		return_data["response"] = "error"
		return_data["data"] = "empty data"
		return return_data
	}
	if !tool.Do_edit_filter(db, config, "", name) || !tool.Do_edit_filter(db, config, "", subject) {
		return_data["response"] = "error"
		return_data["data"] = "edit filter (content)"
		return return_data
	}

	id := ""
	if err := tool.DB_transaction_retry(db, func(tx *sql.Tx) error {
		if tool.Get_DB_type() == "mysql" {
			anchor := ""
			tool.QueryRow_DB(tx, "select data from other where name = 'ver' for update", []any{&anchor})
		}

		last_id := ""
		tool.QueryRow_DB(
			tx,
			"select set_id from calendar where set_name = 'name' order by set_id + 0 desc limit 1",
			[]any{&last_id},
		)
		id = strconv.Itoa(tool.Str_to_int(last_id) + 1)

		for _, value := range [][]string{
			{"name", name},
			{"subject", subject},
			{"user", config.IP},
			{"created", tool.Get_time()},
		} {
			if _, err := tx.Exec(
				tool.DB_change("insert into calendar (set_name, set_code, set_id, set_data) values (?, '', ?, ?)"),
				value[0],
				id,
				value[1],
			); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		panic(err)
	}

	return_data["response"] = "ok"
	return_data["data"] = id
	return return_data
}
