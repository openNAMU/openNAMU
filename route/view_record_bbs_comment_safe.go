package route

import "opennamu/route/tool"

func View_record_bbs_comment_safe(config tool.Config, user_name string, page string) string {
	if !record_bbs_allowed(config, user_name) {
		db := tool.DB_connect()
		defer tool.DB_close(db)
		return tool.Get_error_page(db, config, "auth")
	}
	return View_record_bbs_comment_legacy(config, user_name, page)
}
