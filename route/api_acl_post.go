package route

import (
	"database/sql"
	"net/url"
	"regexp"
	"strings"

	"opennamu/route/tool"
	"opennamu/route/tool/markup"
)

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

func Api_acl_post(config tool.Config, doc_name string, multiple bool, values url.Values) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	allowed := tool.Check_permission(db, "document_acl_manage", config.IP)
	if strings.HasPrefix(doc_name, "user:") && strings.TrimPrefix(doc_name, "user:") == config.IP {
		allowed = true
	}
	if !allowed {
		return_data["response"] = "require auth"
		return return_data
	}

	if values.Get("acl_action") != "" {
		field := values.Get("acl_field")
		group := values.Get("acl_group")
		action := values.Get("acl_action")
		if !tool.Arr_in_str(document_acl_group_fields, field) || !tool.Arr_in_str([]string{"add", "delete"}, action) || group == "" {
			return_data["response"] = "error"
			return_data["data"] = "error"
			return return_data
		}
		if action == "add" && !acl_value_valid(db, group) {
			return_data["response"] = "error"
			return_data["data"] = "error"
			return return_data
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
		return_data["response"] = "ok"
		return return_data
	}

	if multiple {
		for _, name := range strings.Split(strings.ReplaceAll(values.Get("title_name"), "\r", ""), "\n") {
			name = strings.TrimSpace(name)
			if name != "" {
				save_acl(db, name, values)
				acl_history(db, config, name, values)
			}
		}
		return_data["response"] = "ok"
		return return_data
	}
	if doc_name == "" {
		return_data["response"] = "redirect"
		return return_data
	}

	old_markup := document_set_value(db, doc_name, "document_markup")
	if old_markup == "" {
		old_markup = tool.Get_document_markup(db, "", "document")
	}
	save_values := values
	if tool.Check_permission(db, "owner", config.IP) {
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

	return_data["response"] = "ok"
	return return_data
}
