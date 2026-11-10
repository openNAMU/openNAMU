package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_bbs_bookmark(config tool.Config, set_id string, set_code string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	bookmark_data := Api_bbs_bookmark_view(config, set_id, set_code)
	response, _ := bookmark_data["response"].(string)
	if response == "require auth" {
		return tool.Get_redirect("/login")
	}
	if response != "ok" {
		return tool.Get_redirect("/bbs/main")
	}
	if tool.IP_or_user(config.IP) {
		return tool.Get_redirect("/login")
	}
	post_info, ok := bookmark_data["data"].(map[string]string)
	if !ok {
		return tool.Get_redirect("/bbs/main")
	}

	post_path := "/bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code)
	available := post_info["available"] == "true"
	bookmark_label := tool.Get_language(db, "bbs_post_bookmark", true)
	if post_info["watched"] == "true" {
		bookmark_label = tool.Get_language(db, "bbs_post_bookmark_remove", true)
	}
	if values != nil {
		api_data := Api_bbs_bookmark_post(config, set_id, set_code)
		response, _ := api_data["response"].(string)
		if response == "require auth" {
			return tool.Get_redirect("/login")
		}
		if response != "ok" {
			return tool.Get_error_page(db, config, "error")
		}
		return tool.Get_redirect(post_path)
	}

	data := `<form method="post" action="/bbs_bookmark/` + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code) + `">
            <button id="opennamu_save_button" type="submit">` + bookmark_label + `</button>
        </form>`
	display_info := "(" + tool.HTML_escape(post_info["bbs_name"]) + " - " + tool.HTML_escape(post_info["title"]) + ")"
	if !available {
		display_info = "(" + tool.Get_language(db, "bbs_subscription_unavailable", true) + ")"
	}
	return tool.Get_template(
		db,
		config,
		bookmark_label,
		data,
		[]any{display_info},
		[][]any{{"bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code), tool.Get_language(db, "return", true)}},
		map[string]string{},
	)
}

func View_bbs_board_watch(config tool.Config, set_id string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	watch_data := Api_bbs_board_watch_view(config, set_id)
	response, _ := watch_data["response"].(string)
	if response == "require auth" {
		return tool.Get_redirect("/login")
	}
	if response != "ok" {
		return tool.Get_redirect("/bbs/main")
	}
	if tool.IP_or_user(config.IP) {
		return tool.Get_redirect("/login")
	}
	board_info, ok := watch_data["data"].(map[string]string)
	if !ok {
		return tool.Get_redirect("/bbs/main")
	}

	board_path := "/bbs/in/" + tool.Url_parser(set_id)
	available := board_info["available"] == "true"
	board_watch_label := tool.Get_language(db, "bbs_board_watch", true)
	if board_info["watched"] == "true" {
		board_watch_label = tool.Get_language(db, "bbs_board_watch_remove", true)
	}
	if values != nil {
		api_data := Api_bbs_board_watch_post(config, set_id)
		response, _ := api_data["response"].(string)
		if response == "require auth" {
			return tool.Get_redirect("/login")
		}
		if response != "ok" {
			return tool.Get_error_page(db, config, "error")
		}
		return tool.Get_redirect(board_path)
	}

	data := `<form method="post" action="/bbs_board_watch/` + tool.Url_parser(set_id) + `">
            <button id="opennamu_save_button" type="submit">` + board_watch_label + `</button>
        </form>`
	display_info := "(" + tool.HTML_escape(board_info["bbs_name"]) + ")"
	if !available {
		display_info = "(" + tool.Get_language(db, "bbs_subscription_unavailable", true) + ")"
	}
	return tool.Get_template(
		db,
		config,
		board_watch_label,
		data,
		[]any{display_info},
		[][]any{{"bbs/in/" + tool.Url_parser(set_id), tool.Get_language(db, "return", true)}},
		map[string]string{},
	)
}

func View_bbs_tag_watch(config tool.Config, set_id string, tag_token string, values url.Values) string {
	tag, err := tool.Get_base64_decode(tag_token)
	if err != nil || !Bbs_tag_watch_value_valid(tag) {
		if set_id == "" {
			return tool.Get_redirect("/bbs_tag_watch_list")
		}
		return tool.Get_redirect("/bbs/in/" + tool.Url_parser(set_id))
	}

	db := tool.DB_connect()
	defer tool.DB_close(db)

	watch_data := Api_bbs_tag_watch_view(config, set_id, tag)
	response, _ := watch_data["response"].(string)
	if response == "require auth" {
		return tool.Get_redirect("/login")
	}
	if response != "ok" {
		return tool.Get_redirect("/bbs/in/" + tool.Url_parser(set_id))
	}
	if tool.IP_or_user(config.IP) {
		return tool.Get_redirect("/login")
	}
	if _, ok := watch_data["data"].(map[string]string); !ok {
		return tool.Get_redirect("/bbs_tag_watch_list")
	}

	board_path := "/bbs/in/" + tool.Url_parser(set_id) + "/filter/tag/" + tool.Url_parser(tag) + "/1"
	return_path := "bbs/in/" + tool.Url_parser(set_id) + "/filter/tag/" + tool.Url_parser(tag) + "/1"
	if set_id == "" {
		board_path = "/bbs_tag_watch_list"
		return_path = "bbs_tag_watch_list"
	}
	tag_watch_label := tool.Get_language(db, "bbs_tag_watch", true)
	if tag_data, ok := watch_data["data"].(map[string]string); ok && tag_data["watched"] == "true" {
		tag_watch_label = tool.Get_language(db, "bbs_tag_watch_remove", true)
	}
	if values != nil {
		api_data := Api_bbs_tag_watch_post(config, set_id, tag)
		response, _ := api_data["response"].(string)
		if response == "require auth" {
			return tool.Get_redirect("/login")
		}
		if response != "ok" {
			return tool.Get_error_page(db, config, "error")
		}
		return tool.Get_redirect(board_path)
	}

	tag_watch_token := tool.Base64_encode(tag)
	data := `<form method="post" action="/bbs_tag_watch/` + tag_watch_token + `">
            <button id="opennamu_save_button" type="submit">` + tag_watch_label + `</button>
        </form>`
	return tool.Get_template(
		db,
		config,
		tag_watch_label,
		data,
		[]any{"(#" + tool.HTML_escape(tag) + ")"},
		[][]any{{return_path, tool.Get_language(db, "return", true)}},
		map[string]string{},
	)
}
