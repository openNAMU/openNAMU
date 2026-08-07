package route

import "opennamu/route/tool"

func View_bbs_delete(config tool.Config, set_id string, set_code string, comment_code string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	bbs_name_data := Api_bbs_num_to_name(db, set_id)
	bbs_name, _ := bbs_name_data["data"].(string)
	if bbs_name == "" {
		return tool.Get_redirect("/bbs/main")
	}

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_redirect("/bbs/in/" + tool.Url_parser(set_id))
	}

	title := tool.Get_language(db, "delete", true)
	sub := "(" + bbs_name + ")"

	if set_code != "" {
		if comment_code == "" {
			post_data := Api_bbs_w(config, set_id, set_code)
			post, _ := post_data["data"].(map[string]string)
			if post_data["response"] != "ok" || len(post) == 0 {
				return tool.Get_redirect("/bbs/main")
			}
			if _, exists := post["user_id"]; !exists {
				return tool.Get_redirect("/bbs/main")
			}

			title = tool.Get_language(db, "bbs_post_delete", true)
			sub += " (" + tool.HTML_escape(set_code) + ")"
		} else {
			comment_data := Api_bbs_w_comment_one(config, true, "", set_id+"-"+set_code+"-"+comment_code)
			comments, _ := comment_data["data"].([]map[string]string)
			if len(comments) == 0 {
				return tool.Get_redirect("/bbs/main")
			}

			title = tool.Get_language(db, "bbs_comment_delete", true)
			sub += " (" + tool.HTML_escape(set_code) + ") (" + tool.HTML_escape(comment_code) + ")"
		}
	}

	data_html := `<form method="post">
        <span>` + tool.Get_language(db, "delete_warning", true) + `</span>
        <hr class="main_hr">
        <button class="__ON_BUTTON__" type="submit">` + tool.Get_language(db, "delete", true) + `</button>
    </form>`

	menu := [][]any{
		{"bbs/set/" + tool.Url_parser(set_id), tool.Get_language(db, "return", true)},
	}
	if set_code != "" {
		menu = [][]any{
			{"bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code), tool.Get_language(db, "return", true)},
		}
	}

	return tool.Get_template(
		db,
		config,
		title,
		data_html,
		[]any{sub},
		menu,
		map[string]string{},
	)
}

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
