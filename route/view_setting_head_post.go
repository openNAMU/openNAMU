package route

import "opennamu/route/tool"

func View_setting_head_post(config tool.Config, kind string, skin_name string, content string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	name, coverage, _, action, ok := setting_head_info(kind, skin_name)
	if !ok {
		return tool.Get_error_page(db, config, "error")
	}
	api_data := Api_setting_head_post(config, name, coverage, content)
	response, _ := api_data["response"].(string)
	if response == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	}
	if response != "ok" {
		return tool.Get_error_page(db, config, "error")
	}
	if skin_name != "" {
		action += "/" + tool.Url_parser(skin_name)
	}
	return tool.Get_redirect(action)
}
