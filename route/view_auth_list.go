package route

import "opennamu/route/tool"

func View_auth_list(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	can_edit := tool.Check_acl(db, "", "", "owner_auth", config.IP)
	default_groups := map[string]bool{"owner": true, "admin": true, "user": true, "ip": true, "ban": true, "ban_without_login": true, "ban_without_site": true, "email_verified": true, "up_to_level_10": true, "up_to_level_3": true, "trust_a": true, "trust_b": true, "trust_c": true, "trust_d": true}
	data := `<ul>`
	for _, group := range auth_groups(db) {
		data += `<li><a href="/auth/list/add/` + tool.Url_parser(group) + `">` + tool.HTML_escape(group) + `</a>`
		if can_edit && !default_groups[group] {
			data += ` <a href="/auth/list/delete/` + tool.Url_parser(group) + `">(` + tool.Get_language(db, "delete", true) + `)</a>`
		}
		data += `</li>`
	}
	data += `</ul><hr class="main_hr"><a href="/auth/list/add">(` + tool.Get_language(db, "add", true) + `)</a>`

	return tool.Get_template(db, config, tool.Get_language(db, "admin_group_list", true), data, []any{}, [][]any{{"manager", tool.Get_language(db, "return", true)}}, map[string]string{})
}
