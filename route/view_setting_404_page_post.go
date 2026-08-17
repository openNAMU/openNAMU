package route

import "opennamu/route/tool"

func View_setting_404_page_post(config tool.Config, form map[string]string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	page := setting_form_value(form, "page", "404_page")
	if page != "404_file" {
		page = "404_page"
	}
	content := setting_form_value(form, "content", "")

	if setting_form_value(form, "action", "") == "preview" {
		return view_setting_404_page_data(db, config, page, content, true)
	}

	setting_save_value(db, "manage_404_page", "", page)
	setting_save_value(db, "manage_404_page_content", "", content)
	tool.Do_insert_auth_history(db, config.IP, "edit_set (404_page)")

	return tool.Get_redirect("/setting/404_page")
}
