package route

import "opennamu/route/tool"

func View_bbs_in_w_tabom_post(config tool.Config, set_id string, set_code string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	data := Api_bbs_w_tabom_post(config, set_id+"-"+set_code)
	if response, _ := data["response"].(string); response == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	}

	return tool.Get_redirect("/bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code))
}
