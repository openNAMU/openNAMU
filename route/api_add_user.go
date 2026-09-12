package route

import (
	"database/sql"
	"errors"

	"opennamu/route/tool"
)

func Api_add_user(config tool.Config, id string, password string, email string, encode string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if encode == "" {
		encode = tool.Get_main_encode(db)
	}

	password_hash := tool.Password_encode(db, password, encode)
	return api_add_user_hash(config, id, password_hash, email, encode, "", false)
}

func Api_add_user_invite(config tool.Config, id string, password string, email string, encode string, invite_hash string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if encode == "" {
		encode = tool.Get_main_encode(db)
	}

	password_hash := tool.Password_encode(db, password, encode)
	return api_add_user_hash(config, id, password_hash, email, encode, invite_hash, true)
}

func api_add_user_hash(config tool.Config, id string, password_hash string, email string, encode string, invite_hash string, require_invite bool) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

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

	if require_invite && tool.Invite_required(db) && invite_hash == "" {
		return map[string]any{"response": "error", "data": "invite error"}
	}

	invite_error := errors.New("invite error")
	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		if invite_hash != "" && !tool.Invite_consume(tx, invite_hash) {
			return invite_error
		}

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
		if errors.Is(err, invite_error) {
			return map[string]any{"response": "error", "data": "invite error"}
		}
		panic(err)
	}

	return_data := make(map[string]any)
	return_data["response"] = "ok"

	return return_data
}
