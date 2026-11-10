package route

import "opennamu/route/tool"

func View_user_watch_list(config tool.Config, num string, do_type string) string {
	if do_type == "thread_watchlist" {
		page := tool.Str_to_int(num)
		if page > 1 {
			return tool.Get_redirect("/bbs_watch_list_page/" + tool.Url_parser(num))
		}
		return tool.Get_redirect("/bbs_watch_list")
	}

	db := tool.DB_connect()
	defer tool.DB_close(db)

	if tool.IP_or_user(config.IP) {
		return tool.Get_redirect("/login")
	}
	api_data := Api_user_watch_list(config, config.IP, num, do_type)
	response, _ := api_data["response"].(string)
	if response != "ok" {
		return tool.Get_error_page(db, config, "auth")
	}

	data_list, ok := api_data["data"].([]User_watch_item)
	if !ok {
		return tool.Get_error_page(db, config, "error")
	}
	data_html := ""

	if len(data_list) > 0 {
		data_html += "<ul>"
		for _, item := range data_list {
			view_path := "/w/" + tool.Url_parser(item.Data)
			delete_path := "/star_doc/" + tool.Url_parser(item.Data)
			remove_label := tool.Get_language(db, "delete", true)
			if do_type == "watchlist" {
				delete_path = "/watch_list/" + tool.Url_parser(item.Data)
			} else if do_type == "bbs_watchlist" {
				view_path = "/bbs/w/" + tool.Url_parser(item.Set_id) + "/" + tool.Url_parser(item.Set_code)
				delete_path = "/bbs_watch/" + tool.Url_parser(item.Set_id) + "/" + tool.Url_parser(item.Set_code)
				remove_label = tool.Get_language(db, "bbs_comment_watch_remove", true)
			} else if do_type == bbs_bookmark_list_type {
				view_path = "/bbs/w/" + tool.Url_parser(item.Set_id) + "/" + tool.Url_parser(item.Set_code)
				delete_path = "/bbs_bookmark/" + tool.Url_parser(item.Set_id) + "/" + tool.Url_parser(item.Set_code)
				remove_label = tool.Get_language(db, "bbs_post_bookmark_remove", true)
			} else if do_type == bbs_board_watch_list_type {
				view_path = "/bbs/in/" + tool.Url_parser(item.Set_id)
				delete_path = "/bbs_board_watch/" + tool.Url_parser(item.Set_id)
				remove_label = tool.Get_language(db, "bbs_board_watch_remove", true)
			} else if do_type == bbs_tag_watch_list_type {
				view_path = ""
				delete_path = "/bbs_tag_watch/" + tool.Base64_encode(item.Tag)
				remove_label = tool.Get_language(db, "bbs_tag_watch_remove", true)
			}

			display := tool.HTML_escape(item.Display)
			if (do_type == bbs_bookmark_list_type || do_type == bbs_board_watch_list_type) && !item.Available {
				display = tool.Get_language(db, "bbs_subscription_unavailable", true)
				view_path = ""
			}
			if do_type == bbs_tag_watch_list_type {
				display = tool.HTML_escape(item.Display)
			}

			date_html := ""
			if item.Date != "" {
				date_html = "(" + tool.HTML_escape(item.Date) + ") "
			}
			if view_path != "" {
				display = `<a href="` + view_path + `">` + display + `</a>`
			}
			data_html += `<li>` + display + ` ` + date_html + `<a href="` + delete_path + `">(` + remove_label + `)</a></li>`
		}
		data_html += "</ul><hr class=\"main_hr\">"
	}

	page_url := "/star_doc_page/{}"
	if do_type == "watchlist" {
		page_url = "/watch_list_page/{}"
	} else if do_type == "bbs_watchlist" {
		page_url = "/bbs_watch_list_page/{}"
	} else if do_type == bbs_bookmark_list_type {
		page_url = "/bbs_bookmark_list_page/{}"
	} else if do_type == bbs_board_watch_list_type {
		page_url = "/bbs_board_watch_list_page/{}"
	} else if do_type == bbs_tag_watch_list_type {
		page_url = "/bbs_tag_watch_list_page/{}"
	}

	if do_type == "bbs_watchlist" && len(data_list) == 50 {
		page := tool.Str_to_int(num)
		if page < 1 {
			page = 1
		}
		total_count := bbs_watch_merged_count(db, config.IP)
		has_next := page*50 < total_count
		data_html += tool.Get_page_control(db, page, len(data_list), 50, page_url, has_next)
	} else if has_next, ok := api_data["has_next"].(bool); ok {
		data_html += tool.Get_page_control(db, tool.Str_to_int(num), len(data_list), 50, page_url, has_next)
	} else {
		data_html += tool.Get_page_control(db, tool.Str_to_int(num), len(data_list), 50, page_url)
	}

	manager_url := "/manager/16"
	if do_type == "watchlist" {
		manager_url = "/manager/13"
	} else if do_type == "bbs_watchlist" || do_type == bbs_bookmark_list_type || do_type == bbs_board_watch_list_type || do_type == bbs_tag_watch_list_type {
		manager_url = ""
	}
	if manager_url != "" {
		data_html += `<a href="` + manager_url + `">(` + tool.Get_language(db, "add", true) + `)</a>`
	}

	title := tool.Get_language(db, "watchlist", true)
	if do_type == "star_doc" {
		title = tool.Get_language(db, "star_doc", true)
	} else if do_type == "bbs_watchlist" {
		title = tool.Get_language(db, "bbs_watchlist", true)
	} else if do_type == bbs_bookmark_list_type {
		title = tool.Get_language(db, "bbs_post_bookmark_list", true)
	} else if do_type == bbs_board_watch_list_type {
		title = tool.Get_language(db, "bbs_board_watch_list", true)
	} else if do_type == bbs_tag_watch_list_type {
		title = tool.Get_language(db, "bbs_tag_watch_list", true)
	}

	return tool.Get_template(
		db,
		config,
		title,
		data_html,
		[]any{},
		[][]any{
			{"user/" + tool.Url_parser(config.IP), tool.Get_language(db, "return", false)},
			{"watch_list", tool.Get_language(db, "watchlist", false)},
			{"star_doc", tool.Get_language(db, "star_doc", false)},
			{"bbs_watch_list", tool.Get_language(db, "bbs_watchlist", false)},
			{"bbs_bookmark_list", tool.Get_language(db, "bbs_post_bookmark_list", false)},
			{"bbs_board_watch_list", tool.Get_language(db, "bbs_board_watch_list", false)},
			{"bbs_tag_watch_list", tool.Get_language(db, "bbs_tag_watch_list", false)},
		},
		map[string]string{},
	)
}
