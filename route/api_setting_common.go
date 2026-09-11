package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func setting_save_value(db tool.DB_runner, name string, coverage string, data string) {
	if sql_db, ok := db.(*sql.DB); ok {
		if err := tool.DB_transaction(sql_db, func(tx *sql.Tx) error {
			setting_save_value(tx, name, coverage, data)
			return nil
		}); err != nil {
			panic(err)
		}
		return
	}

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

func setting_save_fields(db tool.DB_runner, fields []setting_field, form map[string]string) {
	for _, field := range fields {
		value, exists := form[field.name]
		if !exists {
			value = field.default_value
		}
		setting_save_value(db, field.name, "", value)
	}
}
