package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func Api_login_find_key_post(config tool.Config, user_id string, password string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if user_id == "" || password == "" {
		return_data["response"] = "error"
		return return_data
	}
	password_hash := tool.Password_encode(db, password, tool.Get_user_encode(db, user_id))
	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		user_save(tx, user_id, "pw", password_hash)
		user_delete(tx, user_id, "2fa")
		user_delete(tx, user_id, "2fa_pw")
		user_delete(tx, user_id, "2fa_pw_encode")
		user_delete(tx, user_id, "random_key")
		return nil
	}); err != nil {
		panic(err)
	}

	return_data["response"] = "ok"
	return return_data
}
