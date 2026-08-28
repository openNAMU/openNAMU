package route

import "opennamu/route/tool"

func Api_easter_egg_post(config tool.Config) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if tool.IP_or_user(config.IP) {
		return_data["response"] = "ok"
		return return_data
	}
	found := ""
	if !tool.QueryRow_DB(db, "select id from user_set where id = ? and name = 'get_🥚' limit 1", []any{&found}, config.IP) {
		tool.Exec_DB(db, "insert into user_set (name, id, data) values ('get_🥚', ?, 'Y')", config.IP)
	}
	return_data["response"] = "ok"
	return return_data
}
