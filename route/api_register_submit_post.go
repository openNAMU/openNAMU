package route

import (
	"database/sql"
	stdjson "encoding/json"

	"opennamu/route/tool"
)

func Api_register_submit_post(config tool.Config, id string, pw string, email string, question string, answer string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if id == "" || pw == "" || question == "" {
		return_data["response"] = "error"
		return return_data
	}
	encode := tool.Get_main_encode(db)
	application, err := stdjson.Marshal(map[string]string{
		"id":       id,
		"pw_hash":  tool.Password_encode(db, pw, encode),
		"email":    email,
		"encode":   encode,
		"question": question,
		"answer":   answer,
		"date":     tool.Get_time(),
	})
	if err != nil {
		return_data["response"] = "error"
		return return_data
	}
	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		if _, err := tx.Exec(tool.DB_change("delete from user_set where id = ? and name = 'application'"), id); err != nil {
			return err
		}
		_, err := tx.Exec(tool.DB_change("insert into user_set (id, name, data) values (?, 'application', ?)"), id, string(application))
		return err
	}); err != nil {
		panic(err)
	}

	return_data["response"] = "ok"
	return return_data
}
