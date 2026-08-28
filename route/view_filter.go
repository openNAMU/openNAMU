package route

import "opennamu/route/tool"

func View_filter(config tool.Config, kind string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	spec, ok := get_filter_spec(kind)
	if !ok {
		return tool.Get_error_page(db, config, "error")
	}
	if kind == "edit_filter" && !tool.Check_acl(db, "", "", "edit_filter_view", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	data := `<table id="main_table_set">`
	if kind == "external_image" || kind == "html" {
		header := "domain"
		if kind == "html" {
			header = "tag"
		}
		data += `<tr id="main_table_top_tr"><td>` + tool.Get_language(db, header, true) + `</td></tr>`
	} else {
		data += `<tr id="main_table_top_tr">` +
			`<td id="main_table_width">A</td><td id="main_table_width">B</td><td id="main_table_width">C</td></tr>`
	}
	rows := tool.Query_DB(db, "select html, plus, plus_t from html_filter where kind = ?", spec.db_kind)
	defer rows.Close()

	can_edit := tool.Check_acl(db, "", "", "owner_auth", config.IP)
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
