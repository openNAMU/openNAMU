package route

import "opennamu/route/tool"

func View_bbs_pinned_post(config tool.Config, set_id string, set_code string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	api_data := Api_bbs_w_pinned(config, set_id, set_code, true)
	if api_data["response"] == "ok" {
		return tool.Get_redirect("/bbs/in/" + tool.Url_parser(set_id))
	}
	if api_data["response"] == "require auth" {
		return tool.Get_redirect("/bbs/in/" + tool.Url_parser(set_id))
	}

	return tool.Get_redirect("/bbs/main")
}
