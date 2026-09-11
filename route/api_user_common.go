package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func user_save(db tool.DB_runner, id string, name string, value string) {
	if sql_db, ok := db.(*sql.DB); ok {
		if err := tool.DB_transaction(sql_db, func(tx *sql.Tx) error {
			user_save(tx, id, name, value)
			return nil
		}); err != nil {
			panic(err)
		}
		return
	}

	existing_id := ""
	if tool.QueryRow_DB(db, "select id from user_set where id = ? and name = ?", []any{&existing_id}, id, name) {
		tool.Exec_DB(db, "update user_set set data = ? where id = ? and name = ?", value, id, name)
		return
	}
	tool.Exec_DB(db, "insert into user_set (id, name, data) values (?, ?, ?)", id, name, value)
}

func user_delete(db tool.DB_runner, id string, name string) {
	tool.Exec_DB(db, "delete from user_set where id = ? and name = ?", id, name)
}
