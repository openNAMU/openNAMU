package route

import "opennamu/route/tool"

func View_filter(config tool.Config, kind string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	spec, ok := get_filter_spec(kind)
	if !ok {
		return tool.Get_error_page(db, config, "error")
	}
	if kind == "edit_filter" && !tool.Check_permission(db, "edit_filter_view", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	data := `<table id="main_table_set">`
	if kind == "external_image" || kind == "html" {
		header := "domain"
		if kind == "html" {
			header = "tag"
		}
		data += `<tr id="main_table_top_tr"><td>` + tool.Get_language(db, header, true) + `</td></tr>`
	} else if kind == "replace_filter" {
		data += `<tr id="main_table_top_tr"><td>` + tool.Get_language(db, "regex", true) + `</td><td>` + tool.Get_language(db, "replacement", true) + `</td></tr>`
	} else {
		header := []string{"name", "", ""}
		switch kind {
		case "inter_wiki", "outer_link":
			header = []string{"name", "link", "icon"}
		case "document":
			header = []string{"name", "regex", "acl"}
		case "edit_filter":
			header = []string{"name", "regex", "day"}
		case "template":
			header = []string{"template", "explanation", ""}
		case "edit_top":
			header = []string{"title", "markup", ""}
		case "email_filter":
			header = []string{"email", "", ""}
		case "name_filter", "file_filter":
			header = []string{"regex", "", ""}
		case "image_license":
			header = []string{"license", "", ""}
		case "extension_filter":
			header = []string{"extension", "max_file_size", ""}
		}
		data += `<tr id="main_table_top_tr">`
		for _, value := range header {
			data += `<td id="main_table_width">`
			if value != "" {
				data += tool.Get_language(db, value, true)
			}
			data += `</td>`
		}
		data += `</tr>`
	}
	rows := tool.Get_html_filter_rows(db, spec.db_kind)
	defer rows.Close()

	can_edit := tool.Check_permission(db, "filter_manage", config.IP)
	for rows.Next() {
		name := ""
		plus := ""
		plus_t := ""
		if rows.Scan(&name, &plus, &plus_t) != nil {
			continue
		}

		data += `<tr><td>` + tool.HTML_escape(name)
		if can_edit && kind != "email_filter" && kind != "name_filter" && kind != "file_filter" && kind != "extension_filter" && kind != "image_license" {
			data += ` <a href="/filter/` + kind + `/add/` + tool.Url_parser(name) + `">(` + tool.Get_language(db, "edit", true) + `)</a>`
		}
		if can_edit {
			data += ` <a href="/filter/` + kind + `/del/` + tool.Url_parser(name) + `">(` + tool.Get_language(db, "delete", true) + `)</a>`
		}
		if kind == "external_image" || kind == "html" {
			data += `</td></tr>`
			continue
		}
		data += `</td><td>`
		if kind == "replace_filter" {
			data += tool.HTML_escape(plus) + `</td></tr>`
			continue
		}
		if kind == "inter_wiki" {
			data += `<a class="opennamu_link_out" href="` + filter_safe_link(plus) + `">` + tool.HTML_escape(plus) + `</a>`
		} else {
			data += tool.HTML_escape(plus)
		}
		data += `</td><td>` + tool.HTML_escape(plus_t) + `</td></tr>`
	}
	data += `</table>`

	if can_edit {
		data += `<hr class="main_hr"><a href="/filter/` + kind + `/add">(` + tool.Get_language(db, "add", true) + `)</a>`
	}

	return tool.Get_template(db, config, tool.Get_language(db, spec.title, true), data, []any{}, [][]any{{"manager", tool.Get_language(db, "return", true)}}, map[string]string{})
}
