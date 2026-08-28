package route

import "opennamu/route/tool"

func View_setting_404_page_post(config tool.Config, form map[string]string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	page := setting_form_value(form, "page", "404_page")
	if page != "404_file" {
		page = "404_page"
	}
	content := setting_form_value(form, "content", "")

	if setting_form_value(form, "action", "") == "preview" {
		if !tool.Check_permission(db, "setting_404", config.IP) {
			return tool.Get_error_page(db, config, "auth")
		}
		return view_setting_404_page_data(db, config, page, content, true)
	}

	api_data := Api_setting_404_page_post(config, page, content)
	response, _ := api_data["response"].(string)
	if response == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	}
	if response != "ok" {
		return tool.Get_error_page(db, config, "error")
	}
	return tool.Get_redirect("/setting/404_page")
}
