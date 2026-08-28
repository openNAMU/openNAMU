package route

import (
	"database/sql"
	"net/url"
	"regexp"
	"strings"

	"opennamu/route/tool"
	"opennamu/route/tool/markup"
)

var document_acl_fields = []string{
	"view",
	"decu",
	"document_edit_acl",
	"document_move_acl",
	"document_delete_acl",
	"dis",
}

var document_acl_group_fields = []string{
	"view",
	"decu",
	"document_edit_acl",
	"document_move_acl",
	"document_delete_acl",
	"dis",
}

func acl_field_title(db *sql.DB, field string) string {
	key := map[string]string{
		"view":                "view_acl",
		"decu":                "document_acl",
		"document_edit_acl":   "document_edit_acl",
		"document_move_acl":   "document_move_acl",
		"document_delete_acl": "document_delete_acl",
		"dis":                 "discussion_acl",
	}[field]
	return tool.Get_language(db, key, true)
}

func acl_value(db *sql.DB, doc_name string, field string) string {
	value := ""
	tool.QueryRow_DB(db, "select data from acl where title = ? and type = ? limit 1", []any{&value}, doc_name, field)
	return value
}

func save_acl(db *sql.DB, doc_name string, values url.Values) {
	for _, field := range document_acl_group_fields {
		if _, ok := values[field]; !ok {
			continue
		}
		tool.Exec_DB(db, "delete from acl where title = ? and type = ?", doc_name, field)
		value := values.Get(field)
		if value != "" && value != "normal" {
			tool.Exec_DB(db, "insert into acl (title, data, type) values (?, ?, ?)", doc_name, value, field)
		}
	}

	for _, field := range document_acl_fields {
		tool.Exec_DB(db, "delete from data_set where doc_name = ? and doc_rev = ? and set_name = 'acl_date'", doc_name, field)
		date_value := values.Get(field + "_date")
		if regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}$`).MatchString(date_value) {
			tool.Exec_DB(db, "insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, ?, 'acl_date', ?)", doc_name, field, date_value)
		}
	}

	tool.Exec_DB(db, "delete from data_set where doc_name = ? and set_name = 'document_markup'", doc_name)
	tool.Exec_DB(db, "insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, '', 'document_markup', ?)", doc_name, values.Get("document_markup"))
	for _, field := range []string{"document_top", "document_editor_top"} {
		tool.Exec_DB(db, "delete from data_set where doc_name = ? and set_name = ?", doc_name, field)
		tool.Exec_DB(db, "insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, '', ?, ?)", doc_name, field, values.Get(field))
	}

	tool.Exec_DB(db, "delete from acl where title = ? and type = 'why'", doc_name)
	tool.Exec_DB(db, "insert into acl (title, data, type) values (?, ?, 'why')", doc_name, values.Get("why"))
}

func change_acl_group(db *sql.DB, doc_name string, field string, group string, action string) {
	if action == "add" {
		exist := ""
		if !tool.QueryRow_DB(db, "select data from acl where title = ? and type = ? and data = ? limit 1", []any{&exist}, doc_name, field, group) {
			tool.Exec_DB(db, "insert into acl (title, data, type) values (?, ?, ?)", doc_name, group, field)
		}
		tool.Exec_DB(db, "delete from acl where title = ? and type = ? and data = ''", doc_name, field)
		return
	}

	tool.Exec_DB(db, "delete from acl where title = ? and type = ? and data = ?", doc_name, field, group)
}

func acl_group_select(db *sql.DB) string {
	data := `<select name="acl_group">`
	for _, group := range acl_value_list(db, "") {
		data += `<option value="` + tool.HTML_escape(group) + `">` + tool.HTML_escape(group) + `</option>`
	}
	return data + `</select>`
}

func acl_group_setting(db *sql.DB, title string, field string) string {
	data := `<h3>` + acl_field_title(db, field) + `</h3>`
	groups := tool.Get_acl_data_list(db, title, field)
	if len(groups) == 0 {
		data += `<div>` + tool.Get_language(db, "normal", true) + `</div>`
	}
	for _, group := range groups {
		data += `<div>` + tool.HTML_escape(group) + ` <form method="post" style="display:inline"><input type="hidden" name="name" value="` + tool.HTML_escape(title) + `"><input type="hidden" name="acl_action" value="delete"><input type="hidden" name="acl_field" value="` + tool.HTML_escape(field) + `"><input type="hidden" name="acl_group" value="` + tool.HTML_escape(group) + `"><button type="submit">` + tool.Get_language(db, "delete", true) + `</button></form></div>`
	}
	data += `<form method="post"><input type="hidden" name="name" value="` + tool.HTML_escape(title) + `"><input type="hidden" name="acl_action" value="add"><input type="hidden" name="acl_field" value="` + tool.HTML_escape(field) + `">` + acl_group_select(db) + ` <button type="submit">` + tool.Get_language(db, "add", true) + `</button></form>`
	return data + `<hr class="main_hr">`
}

func acl_group_multiple_setting(db *sql.DB) string {
	data := `<h3>` + tool.Get_language(db, "acl", true) + `</h3><form method="post"><textarea class="opennamu_textarea_500" name="title_name" placeholder="` + tool.Get_language(db, "many_delete_help", true) + `"></textarea><hr class="main_hr"><select name="acl_field">`
	for _, field := range document_acl_group_fields {
		data += `<option value="` + tool.HTML_escape(field) + `">` + tool.HTML_escape(acl_field_title(db, field)) + `</option>`
	}
	data += `</select> ` + acl_group_select(db) + ` <button name="acl_action" value="add" type="submit">` + tool.Get_language(db, "add", true) + `</button> <button name="acl_action" value="delete" type="submit">` + tool.Get_language(db, "delete", true) + `</button></form>`
	return data
}

func View_acl(config tool.Config, doc_name string, multiple bool, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if values != nil && doc_name == "" && !multiple {
		doc_name = values.Get("name")
	}
	if doc_name == "" && !multiple && values == nil {
		data := `<form method="post"><input name="name" placeholder="` + tool.Get_language(db, "document_name", true) + `"><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "go", true) + `</button></form>`
		return tool.Get_template(db, config, tool.Get_language(db, "document_setting", true), data, []any{}, [][]any{{"manager", tool.Get_language(db, "return", true)}}, map[string]string{})
	}

	allowed := tool.Check_acl(db, "", "", "acl_auth", config.IP)
	if strings.HasPrefix(doc_name, "user:") && strings.TrimPrefix(doc_name, "user:") == config.IP {
		allowed = true
	}
	if !allowed {
		return tool.Get_error_page(db, config, "auth")
	}

	if values != nil && values.Get("acl_action") != "" {
		field := values.Get("acl_field")
		group := values.Get("acl_group")
		action := values.Get("acl_action")
		if !tool.Arr_in_str(document_acl_group_fields, field) || !tool.Arr_in_str([]string{"add", "delete"}, action) || group == "" {
			return tool.Get_error_page(db, config, "error")
		}
		if action == "add" && !acl_value_valid(db, group) {
			return tool.Get_error_page(db, config, "error")
		}

		names := []string{doc_name}
		if multiple {
			names = strings.Split(strings.ReplaceAll(values.Get("title_name"), "\r", ""), "\n")
		}
		for _, name := range names {
			name = strings.TrimSpace(name)
			if name == "" {
				continue
			}
			change_acl_group(db, name, field, group, action)
			tool.Do_insert_auth_history(db, config.IP, "document_acl_"+action+" ("+name+")")
		}
		if multiple {
			return tool.Get_redirect("/manager")
		}
		return tool.Get_redirect("/acl/" + tool.Url_parser(doc_name))
	}

	if values != nil {
		if multiple {
			for _, name := range strings.Split(strings.ReplaceAll(values.Get("title_name"), "\r", ""), "\n") {
				name = strings.TrimSpace(name)
				if name != "" {
					save_acl(db, name, values)
					acl_history(db, config, name, values)
				}
			}
			return tool.Get_redirect("/manager")
		}
		if doc_name == "" {
			return tool.Get_redirect("/acl")
		}
		old_markup := document_set_value(db, doc_name, "document_markup")
		if old_markup == "" {
			old_markup = tool.Get_document_markup(db, "", "document")
		}
		save_values := values
		if tool.Check_acl(db, "", "", "owner_auth", config.IP) {
			save_values = url.Values{}
			for key, value := range values {
				save_values[key] = append([]string{}, value...)
			}
			save_values.Set("document_top", document_set_value(db, doc_name, "document_top"))
			save_values.Set("document_editor_top", document_set_value(db, doc_name, "document_editor_top"))
		}
		save_acl(db, doc_name, save_values)
		new_markup := document_set_value(db, doc_name, "document_markup")
		if new_markup == "" {
			new_markup = tool.Get_document_markup(db, "", "document")
		}
		if old_markup != new_markup {
			data := ""
			if tool.QueryRow_DB(db, "select data from data where title = ?", []any{&data}, doc_name) {
				markup.Get_render(db, doc_name, data, "backlink")
			}
		}
		acl_history(db, config, doc_name, save_values)
		return tool.Get_redirect("/acl/" + tool.Url_parser(doc_name))
	}

	data := `<form method="post">`
	if multiple {
		data += `<textarea class="opennamu_textarea_500" name="title_name" placeholder="` + tool.Get_language(db, "many_delete_help", true) + `"></textarea><hr class="main_hr">`
	} else {
		data += `<input type="hidden" name="name" value="` + tool.HTML_escape(doc_name) + `">`
	}

	for _, field := range document_acl_fields {
		date_value := ""
		if !multiple {
			tool.QueryRow_DB(db, "select set_data from data_set where doc_name = ? and doc_rev = ? and set_name = 'acl_date' limit 1", []any{&date_value}, doc_name, field)
		}
		data += `<input type="date" name="` + field + `_date" value="` + tool.HTML_escape(date_value) + `"><hr class="main_hr">`
	}

	why := ""
	if !multiple {
		why = acl_value(db, doc_name, "why")
	}
	data += `<h3>` + tool.Get_language(db, "why", true) + `</h3><input name="why" value="` + tool.HTML_escape(why) + `"><hr class="main_hr">`

	markup := ""
	if !multiple {
		markup = tool.Get_document_markup(db, doc_name, "document")
	}
	data += `<h2>` + tool.Get_language(db, "markup", true) + `</h2>` + tool.Get_markup_select_ui(db, config, doc_name, markup, "", "")
	top_disabled := ""
	if tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		top_disabled = ` disabled`
	}
	data += `<h2>` + tool.Get_language(db, "document_top", true) + `</h2><textarea class="opennamu_textarea_100" name="document_top"` + top_disabled + `>` + tool.HTML_escape(document_set_value(db, doc_name, "document_top")) + `</textarea>`
	data += `<h2>` + tool.Get_language(db, "document_editor_top", true) + `</h2><textarea class="opennamu_textarea_100" name="document_editor_top"` + top_disabled + `>` + tool.HTML_escape(document_set_value(db, doc_name, "document_editor_top")) + `</textarea><hr class="main_hr">`
	data += `<button type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`
	if multiple {
		data += acl_group_multiple_setting(db)
	} else {
		for _, field := range document_acl_group_fields {
			data += acl_group_setting(db, doc_name, field)
		}
	}

	menu := [][]any{{"manager", tool.Get_language(db, "admin", true)}}
	if !multiple {
		menu = [][]any{{"w/" + tool.Url_parser(doc_name), tool.Get_language(db, "return", true)}, {"acl_multiple", tool.Get_language(db, "mutiple_document_setting", true)}, {"manager", tool.Get_language(db, "admin", true)}}
	}
	title := doc_name
	if multiple {
		title = tool.Get_language(db, "mutiple_document_setting", true)
	}
	return tool.Get_template(db, config, title, data, []any{}, menu, map[string]string{})
}

func acl_history(db *sql.DB, config tool.Config, doc_name string, values url.Values) {
	data := ""
	for _, field := range document_acl_fields {
		field_data := values.Get(field)
		if tool.Arr_in_str(document_acl_group_fields, field) {
			field_data = strings.Join(tool.Get_acl_data_list(db, doc_name, field), "\n")
		}
		data += field + "\n" + field_data + "\n\n"
	}
	for _, field := range []string{"why", "document_markup"} {
		data += field + "\n" + values.Get(field) + "\n\n"
	}
	tool.Do_add_history(db, doc_name, data, tool.Get_time(), config.IP, values.Get("why"), "0", "setting", "")
	tool.Do_insert_auth_history(db, config.IP, "document_set ("+doc_name+")")
}

func document_set_value(db *sql.DB, doc_name string, set_name string) string {
	value := ""
	tool.QueryRow_DB(db, "select set_data from data_set where doc_name = ? and set_name = ? limit 1", []any{&value}, doc_name, set_name)
	return value
}
