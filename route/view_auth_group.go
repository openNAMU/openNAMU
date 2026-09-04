package route

import (
	"database/sql"
	"net/url"
	"strconv"

	"opennamu/route/tool"
)

func View_auth_group(config tool.Config, name string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	new_group := name == ""
	can_edit := tool.Check_permission(db, "auth_group_manage", config.IP)

	if values == nil && name == "" && !can_edit {
		return tool.Get_error_page(db, config, "auth")
	}
	if values == nil && name != "" && !tool.Check_permission(db, "site_view", config.IP) {
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

	data := ""
	if can_edit {
		data = `<form method="post"><input type="hidden" name="group_name" value="` + tool.HTML_escape(name) + `">`
	}
	for _, choice := range tool.Auth_choices() {
		checked := ""
		if selected[choice.Key] {
			checked = ` checked`
		}
		disabled := ""
		if !can_edit {
			disabled = ` disabled`
		}
		label := tool.Get_language(db, choice.Lang, true)
		data += `<div class="opennamu_list_1" style="margin-left:` + strconv.Itoa((choice.Level-1)*20) + `px"><label><input type="checkbox" name="` + choice.Key + `"` + checked + disabled + `> ` + label + `</label></div>`
	}
	if can_edit {
		data += `<hr class="main_hr"><button type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`
	}
	data += auth_relation_view(db)

	return tool.Get_template(db, config, name, data, []any{"(" + tool.Get_language(db, "admin_group", true) + ")"}, [][]any{{"auth/list", tool.Get_language(db, "return", true)}}, map[string]string{})
}

func auth_relation_label(db *sql.DB, key string) string {
	if key == "user_default" || key == "ip_default" {
		return tool.HTML_escape(tool.Get_language(db, "auth_"+key, true))
	}

	for _, choice := range tool.Auth_choices() {
		if choice.Key == key {
			return tool.HTML_escape(tool.Get_language(db, choice.Lang, true))
		}
	}

	return tool.HTML_escape(key)
}

func auth_relation_view(db *sql.DB) string {
	data := `<hr class="main_hr"><h3>` + tool.Get_language(db, "auth_relation", true) + `</h3>`
	relations := tool.Auth_relations()

	for _, relation_type := range []string{tool.Auth_relation_group, tool.Auth_relation_permission, tool.Auth_relation_auto} {
		data += `<h3>` + tool.Get_language(db, "auth_relation_"+relation_type, true) + `</h3>`

		if relation_type == tool.Auth_relation_auto {
			for _, default_name := range []string{"user_default", "ip_default"} {
				source_name := "user"
				if default_name == "ip_default" {
					source_name = "ip"
				}
				data += `<p><b>` + auth_relation_label(db, default_name) + `</b></p><ul>`
				for _, relation := range relations {
					if relation.Type == relation_type && relation.From == source_name {
						data += `<li>` + auth_relation_label(db, relation.To) + `</li>`
					}
				}
				data += `</ul>`
			}
			continue
		}

		data += `<ul>`
		for _, relation := range relations {
			if relation.Type == relation_type {
				data += `<li>` + auth_relation_label(db, relation.From) + ` → ` + auth_relation_label(db, relation.To) + `</li>`
			}
		}
		data += `</ul>`
	}

	data += `<p>` + tool.Get_language(db, "auth_relation_note", true) + `</p>`
	return data
}
