package route

import "opennamu/route/tool"

func Api_user_key_post(config tool.Config) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if tool.IP_or_user(config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	value := ""
	for value == "" {
		value = tool.Get_random_key(128)
		existing := ""
		if tool.QueryRow_DB(db, "select data from user_set where name = ? and data = ?", []any{&existing}, "random_key", value) {
			value = ""
		}
	}
	user_save(db, config.IP, "random_key", value)

	return_data["response"] = "ok"
	return return_data
}
