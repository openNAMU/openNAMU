package route

import "opennamu/route/tool"

func View_setting_top_menu_post(config tool.Config, form map[string]string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	api_data := Api_setting_top_menu_post(config, form)
	return tool.Api_post_redirect(db, config, api_data, "/setting/top_menu")
}
