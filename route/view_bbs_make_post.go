package route

import "opennamu/route/tool"

func View_bbs_make_post(config tool.Config, bbs_name string, bbs_type string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	api_data := Api_bbs_make(config, bbs_name, bbs_type)

	response, _ := api_data["response"].(string)
	if response != "ok" {
		return tool.Get_error_page(
			db,
			config,
			"auth",
		)
	}

	return tool.Get_redirect("/bbs/main")
}
