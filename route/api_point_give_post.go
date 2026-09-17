package route

import (
	"strconv"
	"strings"

	"opennamu/route/tool"
)

func Api_point_give_post(config tool.Config, user_name string, amount string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "admin", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	user_name = strings.TrimSpace(user_name)
	point_amount, err := strconv.Atoi(strings.TrimSpace(amount))
	if user_name == "" || err != nil || point_amount <= 0 {
		return_data["response"] = "error"
		return_data["data"] = "invalid data"
		return return_data
	}

	existing_id := ""
	if !tool.QueryRow_DB(db, "select id from user_set where id = ? and name = 'pw' limit 1", []any{&existing_id}, user_name) {
		return_data["response"] = "error"
		return_data["data"] = "user not found"
		return return_data
	}

	new_point := tool.Change_user_point(db, user_name, point_amount)
	return_data["response"] = "ok"
	return_data["data"] = map[string]any{"user_name": user_name, "point": new_point}
	return return_data
}
