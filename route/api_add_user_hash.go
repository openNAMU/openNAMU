package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func Api_add_user_hash(config tool.Config, id string, password_hash string, email string, encode string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if encode == "" {
		encode = tool.Get_main_encode(db)
	}

	temp := ""
	tool.QueryRow_DB(
		db,
		`select id from user_set limit 1`,
		[]any{&temp},
	)

	auth := "user"
	if temp == "" {
		auth = "owner"
	}

	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		for _, value := range [][]any{
			{id, "pw", password_hash},
			{id, "acl", auth},
			{id, "date", tool.Get_time()},
			{id, "encode", encode},
		} {
			if _, err := tx.Exec(
				tool.DB_change(`insert into user_set (id, name, data) values (?, ?, ?)`),
				value...,
			); err != nil {
				return err
			}
		}
		if email != "" {
			_, err := tx.Exec(
				tool.DB_change(`insert into user_set (id, name, data) values (?, 'email', ?)`),
				id,
				email,
			)
			return err
		}
		return nil
	}); err != nil {
		panic(err)
	}

	return_data := make(map[string]any)
	return_data["response"] = "ok"

	return return_data
}
