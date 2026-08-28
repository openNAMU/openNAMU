package route

import (
	"net/url"
	"strconv"

	"opennamu/route/tool"
)

func View_auth_group(config tool.Config, name string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	new_group := name == ""

	if values == nil && !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	if name == "" && values == nil {
		data := `<form method="post"><input name="group_name" placeholder="` + tool.Get_language(db, "name", true) + `"><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "go", true) + `</button></form>`
		return tool.Get_template(db, config, tool.Get_language(db, "add_admin_group", true), data, []any{}, [][]any{{"auth/list", tool.Get_language(db, "return", true)}}, map[string]string{})
	}
	if values != nil && name == "" {
		name = values.Get("group_name")
	}
	if name == "" || tool.HTML_escape(name) != name || (tool.Auth_group_name_reserved(name) && (new_group || !tool.Auth_group_name_default(name))) {
		return tool.Get_error_page(db, config, "error")
	}

	if values != nil {
		result := Api_auth_group_post(config, name, values)
		if result["response"] != "ok" {
			return tool.Get_error_page(db, config, "error")
		}
		return tool.Get_redirect("/auth/list/add/" + tool.Url_parser(name))
	}

	selected := map[string]bool{}
	for _, value := range tool.Get_auth_permission_list(db, name) {
		selected[value] = true
	}

	data := `<form method="post"><input type="hidden" name="group_name" value="` + tool.HTML_escape(name) + `">`
	for _, choice := range tool.Auth_choices() {
		checked := ""
		if selected[choice.Key] {
			checked = ` checked`
		}
		label := tool.Get_language(db, choice.Lang, true)
		data += `<div class="opennamu_list_1" style="margin-left:` + strconv.Itoa((choice.Level-1)*20) + `px"><label><input type="checkbox" name="` + choice.Key + `"` + checked + `> ` + label + `</label></div>`
	}
	data += `<hr class="main_hr"><button type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`

	return tool.Get_template(db, config, name, data, []any{"(" + tool.Get_language(db, "admin_group", true) + ")"}, [][]any{{"auth/list", tool.Get_language(db, "return", true)}}, map[string]string{})
}
