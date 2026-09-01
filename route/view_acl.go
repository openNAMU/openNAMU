package route

import (
	"database/sql"
	"net/url"
	"strings"

	"opennamu/route/tool"
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
	value := tool.Get_acl_data_list(db, doc_name, field)
	if len(value) == 0 {
		return ""
	}
	return value[0]
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

	if values == nil {
		allowed := tool.Check_permission(db, "document_acl_manage", config.IP)
		if strings.HasPrefix(doc_name, "user:") && strings.TrimPrefix(doc_name, "user:") == config.IP {
			allowed = true
		}
		if !allowed {
			return tool.Get_error_page(db, config, "auth")
		}
	}

	if values != nil {
		api_data := Api_acl_post(config, doc_name, multiple, values)
		response, _ := api_data["response"].(string)
		if response == "require auth" {
			return tool.Get_error_page(db, config, "auth")
		}
		if response == "redirect" {
			return tool.Get_redirect("/acl")
		}
		if response != "ok" {
			error_name, _ := api_data["data"].(string)
			if error_name == "" {
				error_name = "error"
			}
			return tool.Get_error_page(db, config, error_name)
		}
		if multiple {
			return tool.Get_redirect("/manager")
		}
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
			date_value = tool.Get_document_setting_value(db, doc_name, "acl_date", field)
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
	if tool.Check_permission(db, "owner", config.IP) {
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

func document_set_value(db *sql.DB, doc_name string, set_name string) string {
	return tool.Get_document_setting_value(db, doc_name, set_name, "")
}
