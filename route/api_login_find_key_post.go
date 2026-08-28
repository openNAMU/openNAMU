package route

import "opennamu/route/tool"

func Api_login_find_key_post(config tool.Config, user_id string, password string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if user_id == "" || password == "" {
		return_data["response"] = "error"
		return return_data
	}
	user_save(db, user_id, "pw", tool.Password_encode(db, password, tool.Get_user_encode(db, user_id)))
	user_delete(db, user_id, "2fa")
	user_delete(db, user_id, "2fa_pw")
	user_delete(db, user_id, "2fa_pw_encode")
	user_delete(db, user_id, "random_key")

	return_data["response"] = "ok"
	return return_data
}
