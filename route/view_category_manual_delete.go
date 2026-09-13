package route

import "opennamu/route/tool"

func View_category_manual_delete(config tool.Config, category_value string, document_value string, return_value string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	category_name, category_err := tool.Get_base64_decode(category_value)
	document_name, document_err := tool.Get_base64_decode(document_value)
	return_name, return_err := tool.Get_base64_decode(return_value)
	if category_err != nil || document_err != nil || category_name == "" || document_name == "" {
		return tool.Get_error_page(db, config, "not found")
	}
	if return_err != nil || return_name == "" {
		return_name = document_name
	}

	category_name = category_manual_name(category_name)
	if _, exists := tool.Get_data_title(db, category_name); !exists {
		return tool.Get_error_page(db, config, "not found")
	}
	if _, exists := tool.Get_data_title(db, document_name); !exists {
		return tool.Get_error_page(db, config, "not found")
	}
	if !tool.Check_acl(db, document_name, "", "document_edit", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	data := `<form method="post" action="/category/delete"><input type="hidden" name="category" value="` + tool.HTML_escape(category_name) + `"><input type="hidden" name="document" value="` + tool.HTML_escape(document_name) + `"><input type="hidden" name="return" value="` + tool.HTML_escape(return_name) + `"><button type="submit">` + tool.Get_language(db, "delete", true) + `</button></form>`
	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, "delete", true),
		data,
		[]any{},
		[][]any{{"w/" + tool.Url_parser(return_name), tool.Get_language(db, "return", true)}},
		map[string]string{},
	)
}
