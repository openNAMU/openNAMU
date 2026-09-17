package route

import "opennamu/route/tool"

func View_setting_404_page_post(config tool.Config, form map[string]string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	page := Setting_form_value(form, "page", "404_page")
	if page != "404_file" {
		page = "404_page"
	}
	content := Setting_form_value(form, "content", "")

	if Setting_form_value(form, "action", "") == "preview" {
		if !tool.Check_permission(db, "setting_404", config.IP) {
			return tool.Get_error_page(db, config, "auth")
		}
		return View_setting_404_page_data(db, config, page, content, true)
	}

	api_data := Api_setting_404_page_post(config, page, content)
	return tool.Api_post_redirect(db, config, api_data, "/setting/404_page")
}
