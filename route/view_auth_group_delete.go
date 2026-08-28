package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_auth_group_delete(config tool.Config, name string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	default_groups := map[string]bool{"owner": true, "admin": true, "user": true, "ip": true, "ban": true, "ban_without_login": true, "ban_without_site": true, "email_verified": true, "up_to_level_10": true, "up_to_level_3": true, "trust_a": true, "trust_b": true, "trust_c": true, "trust_d": true}
	if default_groups[name] {
		return tool.Get_redirect("/auth/list")
	}
	if values == nil && !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	if values != nil {
		result := Api_auth_group_delete_post(config, name)
		if result["response"] != "ok" {
			return tool.Get_error_page(db, config, "error")
		}
		return tool.Get_redirect("/auth/list")
	}

	data := `<form method="post"><button type="submit">` + tool.Get_language(db, "delete", true) + `</button></form>`
	return tool.Get_template(db, config, tool.Get_language(db, "delete_admin_group", true), data, []any{"(" + tool.HTML_escape(name) + ")"}, [][]any{{"auth/list", tool.Get_language(db, "return", true)}}, map[string]string{})
}
