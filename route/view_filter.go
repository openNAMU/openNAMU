package route

import "opennamu/route/tool"

func View_filter(config tool.Config, kind string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	spec, ok := Get_filter_spec(kind)
	if !ok {
		return tool.Get_error_page(db, config, "error")
	}
	if kind == "edit_filter" && !tool.Check_permission(db, "edit_filter_view", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	data := ""
	headers := []string{"name", "", ""}
	if kind == "external_image" || kind == "html" {
		headers = []string{"domain"}
		if kind == "html" {
			headers[0] = "tag"
		}
	} else if kind == "replace_filter" {
		headers = []string{"regex", "replacement"}
	} else {
		switch kind {
		case "inter_wiki", "outer_link":
			headers = []string{"name", "link", "icon"}
		case "document":
			headers = []string{"name", "regex", "acl"}
		case "edit_filter":
			headers = []string{"name", "regex", "day"}
		case "template":
			headers = []string{"template", "explanation", ""}
		case "edit_top":
			headers = []string{"title", "markup", ""}
		case "email_filter":
			headers = []string{"email", "", ""}
		case "name_filter", "file_filter":
			headers = []string{"regex", "", ""}
		case "image_license":
			headers = []string{"license", "", ""}
		case "extension_filter":
			headers = []string{"extension", "max_file_size", ""}
		}
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

		name_html := `<strong>` + tool.Get_language(db, headers[0], true) + `:</strong> ` + tool.HTML_escape(name)
		if can_edit && kind != "email_filter" && kind != "name_filter" && kind != "file_filter" && kind != "extension_filter" && kind != "image_license" {
			name_html += ` <a href="/filter/` + kind + `/add/` + tool.Url_parser(name) + `">(` + tool.Get_language(db, "edit", true) + `)</a>`
		}
		action_html := ""
		if can_edit {
			action_html = `<a href="/filter/` + kind + `/del/` + tool.Url_parser(name) + `">(` + tool.Get_language(db, "delete", true) + `)</a>`
		}

		bottom := ""
		if len(headers) > 1 && headers[1] != "" {
			value_html := tool.HTML_escape(plus)
			if kind == "inter_wiki" {
				value_html = `<a class="opennamu_link_out" href="` + Filter_safe_link(plus) + `">` + tool.HTML_escape(plus) + `</a>`
			}
			bottom = `<div><strong>` + tool.Get_language(db, headers[1], true) + `:</strong> ` + value_html + `</div>`
		}
		if len(headers) > 2 && headers[2] != "" {
			bottom += `<div><strong>` + tool.Get_language(db, headers[2], true) + `:</strong> ` + tool.HTML_escape(plus_t) + `</div>`
		}
		data += tool.Get_list_ui(name_html, action_html, bottom, "")
	}

	if can_edit {
		data += `<hr class="main_hr"><a href="/filter/` + kind + `/add">(` + tool.Get_language(db, "add", true) + `)</a>`
	}

	return tool.Get_template(db, config, tool.Get_language(db, spec.title, true), data, []any{}, [][]any{{"manager", tool.Get_language(db, "return", true)}}, map[string]string{})
}
