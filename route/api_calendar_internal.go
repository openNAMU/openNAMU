package route

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"opennamu/route/tool"
)

var calendar_missing_error = errors.New("calendar does not exist")
var calendar_event_missing_error = errors.New("calendar event does not exist")

func Calendar_id_normalize(data string) (string, bool) {
	id := tool.Str_to_int(strings.TrimSpace(data))
	if id <= 0 {
		return "", false
	}

	return strconv.Itoa(id), true
}

func Calendar_code_normalize(data string) (string, bool) {
	code := tool.Str_to_int(strings.TrimSpace(data))
	if code <= 0 {
		return "", false
	}

	return strconv.Itoa(code), true
}

func Calendar_view_allowed(db *sql.DB, ip string) bool {
	return tool.Get_auth_info(db, ip)["site_view"]
}

func Calendar_text_normalize(data string) string {
	return strings.TrimSpace(strings.ReplaceAll(data, "\r", ""))
}

func Calendar_meta_get(db tool.DB_runner, id string) (map[string]string, bool) {
	rows := tool.Query_DB(
		db,
		"select set_name, set_data from calendar where set_id = ? and set_code = '' and set_name in ('name', 'subject', 'user', 'created')",
		id,
	)
	defer rows.Close()

	data := map[string]string{}
	for rows.Next() {
		name := ""
		value := ""
		if err := rows.Scan(&name, &value); err != nil {
			panic(err)
		}
		if name == "user" {
			name = "owner"
		}
		data[name] = value
	}

	if data["name"] == "" {
		return map[string]string{}, false
	}

	return data, true
}

func Calendar_event_get(db tool.DB_runner, id string, code string) (map[string]string, bool) {
	rows := tool.Query_DB(
		db,
		"select set_name, set_data from calendar where set_id = ? and set_code = ? and set_name in ('date', 'title', 'memo', 'user', 'created', 'updated')",
		id,
		code,
	)
	defer rows.Close()

	data := map[string]string{}
	for rows.Next() {
		name := ""
		value := ""
		if err := rows.Scan(&name, &value); err != nil {
			panic(err)
		}
		data[name] = value
	}

	if data["date"] == "" {
		return map[string]string{}, false
	}

	return data, true
}

func Calendar_manage_allowed(db *sql.DB, config tool.Config, calendar_owner string, event_user string) bool {
	if tool.Check_permission(db, "calendar_manage", config.IP) {
		return true
	}
	if calendar_owner != "" && calendar_owner == config.IP {
		return true
	}
	if event_user != "" && event_user == config.IP && tool.Check_permission(db, "calendar_create", config.IP) {
		return true
	}

	return false
}

// Calendar_write_lock verifies the calendar still exists inside the write
// transaction and serializes code allocation. MySQL locks the parent row with
// FOR UPDATE. SQLite/Turso take the write lock first with the same temporary
// append/restore update used by BBS identifier transactions, so two writers
// cannot read the same last event code before either writes.
func Calendar_write_lock(tx *sql.Tx, id string) (bool, error) {
	if tool.Get_DB_type() == "mysql" {
		name := ""
		exists := tool.QueryRow_DB(
			tx,
			"select set_data from calendar where set_id = ? and set_code = '' and set_name = 'name' for update",
			[]any{&name},
			id,
		)
		return exists, nil
	}

	result, err := tx.Exec(
		tool.DB_change("update calendar set set_data = set_data || char(1) where set_id = ? and set_code = '' and set_name = 'name'"),
		id,
	)
	if err != nil {
		return false, err
	}
	rows_affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if rows_affected == 0 {
		return false, nil
	}

	_, err = tx.Exec(
		tool.DB_change("update calendar set set_data = substr(set_data, 1, length(set_data) - 1) where set_id = ? and set_code = '' and set_name = 'name'"),
		id,
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

func Calendar_next_event_code(tx *sql.Tx, id string) string {
	last_code := ""
	tool.QueryRow_DB(
		tx,
		"select set_code from calendar where set_id = ? and set_name = 'date' order by set_code + 0 desc limit 1",
		[]any{&last_code},
		id,
	)

	return strconv.Itoa(tool.Str_to_int(last_code) + 1)
}
