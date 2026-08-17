package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_auth_group_delete(config tool.Config, name string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	default_groups := map[string]bool{"owner": true, "user": true, "ip": true, "ban": true}
	if default_groups[name] {
		return tool.Get_redirect("/auth/list")
	}
	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	if values != nil {
		used := ""
		if tool.QueryRow_DB(db, "select id from user_set where name = 'acl' and data = ? limit 1", []any{&used}, name) {
			return tool.Get_error_page(db, config, "error")
		}
		tool.Exec_DB(db, "delete from alist where name = ?", name)
		tool.Do_insert_auth_history(db, config.IP, "auth_group_delete ("+name+")")
		return tool.Get_redirect("/auth/list")
	}

	data := `<form method="post"><button type="submit">` + tool.Get_language(db, "delete", true) + `</button></form>`
	return tool.Get_template(db, config, tool.Get_language(db, "delete_admin_group", true), data, []any{"(" + tool.HTML_escape(name) + ")"}, [][]any{{"auth/list", tool.Get_language(db, "return", true)}}, map[string]string{})
}
