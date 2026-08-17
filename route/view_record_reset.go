package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_record_reset(config tool.Config, user_name string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	if user_name == "" {
		return tool.Get_error_page(db, config, "error")
	}
	if values != nil {
		tool.Exec_DB(db, "delete from history where ip = ?", user_name)
		tool.Do_insert_auth_history(db, config.IP, "record_reset ("+user_name+")")
		return tool.Get_redirect("/record/" + tool.Url_parser(user_name))
	}
	return history_destructive_page(db, config, tool.Get_language(db, "record_reset", true), tool.Get_language(db, "reset", true), "record/"+tool.Url_parser(user_name))
}
