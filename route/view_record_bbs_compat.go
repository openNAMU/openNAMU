package route

import "opennamu/route/tool"

func record_bbs_allowed(config tool.Config, user_name string) bool {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	return user_name == config.IP || tool.Check_permission(db, "hidel", config.IP)
}
