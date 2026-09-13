package route

import "opennamu/route/tool"

func View_category_manual_add(config tool.Config, add_type string, value string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	name, err := tool.Get_base64_decode(value)
	if err != nil || name == "" {
		return tool.Get_error_page(db, config, "not found")
	}

	category_name := ""
	document_name := ""
	return_name := name
	switch add_type {
	case "category":
		if tool.Check_permission(db, "edit", config.IP) {
			category_name = category_manual_name(name)
		} else {
			return tool.Get_error_page(db, config, "auth")
		}
	case "document":
		document_name = name
		if !tool.Check_acl(db, document_name, "", "document_edit", config.IP) {
			return tool.Get_error_page(db, config, "auth")
		}
	default:
		return tool.Get_error_page(db, config, "not found")
	}

	if add_type == "category" {
		if _, exists := tool.Get_data_title(db, category_name); !exists {
			return tool.Get_error_page(db, config, "not found")
		}
	} else if add_type == "document" {
		if _, exists := tool.Get_data_title(db, document_name); !exists {
			return tool.Get_error_page(db, config, "not found")
		}
	}

	data := category_manual_add_form(db, category_name, document_name, return_name)
	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, "category_manual", true),
		data,
		[]any{},
		[][]any{{"w/" + tool.Url_parser(return_name), tool.Get_language(db, "return", true)}},
		map[string]string{},
	)
}
