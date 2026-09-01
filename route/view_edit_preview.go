package route

import "opennamu/route/tool"

func View_edit_preview(config tool.Config, doc_name string, data string, mode string, send string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, doc_name, "", "document_edit", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	darkmode := "0"
	preview_name := "preview_normal"
	if mode == "dark" {
		darkmode = "1"
		preview_name = "preview_dark"
	}

	preview_config := config
	preview_config.Cookies += "; main_css_darkmode=" + darkmode

	api_data := Api_w_render(preview_config, doc_name, data, "api_view", "")
	if api_data["response"] != "ok" {
		return tool.Get_error_page(db, config, "error")
	}

	rendered_data, ok := api_data["data"].(string)
	if !ok {
		return tool.Get_error_page(db, config, "error")
	}

	return view_edit_page(
		db,
		preview_config,
		doc_name,
		"",
		data,
		send,
		tool.Get_language(db, preview_name, true),
		rendered_data,
	)
}
