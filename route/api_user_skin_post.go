package route

import "opennamu/route/tool"

func Api_user_skin_post(config tool.Config, skin string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if tool.IP_or_user(config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if !tool.Arr_in_str(tool.Get_skin_list("ringo", true), skin) {
		return_data["response"] = "error"
		return return_data
	}
	user_save(db, config.IP, "skin", skin)
	return_data["response"] = "ok"
	return return_data
}
