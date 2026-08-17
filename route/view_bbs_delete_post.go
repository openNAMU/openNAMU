package route

import "opennamu/route/tool"

func View_bbs_delete_post(config tool.Config, set_id string, set_code string, comment_code string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	var api_data map[string]any
	if set_code == "" {
		api_data = Api_bbs_delete(config, set_id)
	} else if comment_code == "" {
		api_data = Api_bbs_w_delete(config, set_id, set_code)
	} else {
		api_data = Api_bbs_w_comment_one_delete(config, set_id, set_code+"-"+comment_code)
	}

	if api_data["response"] == "ok" {
		if set_code == "" {
			return tool.Get_redirect("/bbs/main")
		}
		if comment_code == "" {
			return tool.Get_redirect("/bbs/in/" + tool.Url_parser(set_id))
		}
		return tool.Get_redirect("/bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code))
	}

	if api_data["response"] == "require auth" {
		return tool.Get_redirect("/bbs/in/" + tool.Url_parser(set_id))
	}

	return tool.Get_redirect("/bbs/main")
}
