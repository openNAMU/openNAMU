package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func setting_save_value(db *sql.DB, name string, coverage string, data string) {
	old_data := ""
	exists := tool.QueryRow_DB(
		db,
		"select data from other where name = ? and coverage = ?",
		[]any{&old_data},
		name,
		coverage,
	)
	if exists {
		tool.Exec_DB(
			db,
			"update other set data = ? where name = ? and coverage = ?",
			data,
			name,
			coverage,
		)
		return
	}
	tool.Exec_DB(
		db,
		"insert into other (name, data, coverage) values (?, ?, ?)",
		name,
		data,
		coverage,
	)
}

func setting_save_fields(db *sql.DB, fields []setting_field, form map[string]string) {
	for _, field := range fields {
		value, exists := form[field.name]
		if !exists {
			value = field.default_value
		}
		setting_save_value(db, field.name, "", value)
	}
}
