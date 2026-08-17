package route

import "opennamu/route/tool"

func View_setting_sitemap_post(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	api_data := Api_setting_sitemap_post(config)
	response, _ := api_data["response"].(string)
	if response == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	}
	if response != "ok" {
		return tool.Get_error_page(db, config, "error")
	}

	return tool.Get_redirect("/setting/sitemap")
}
