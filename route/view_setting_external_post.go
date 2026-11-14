package route

import "opennamu/route/tool"

func View_setting_external_post(config tool.Config, form map[string]string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	api_data := Api_setting_external_post(config, form)
	if api_data["response"] == "error" && api_data["data"] == "captcha configuration" {
		return tool.Get_error_page(db, config, "captcha configuration")
	}
	return tool.Api_post_redirect(db, config, api_data, "/setting/external")
}
