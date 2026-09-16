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
	if api_data["response"].(string) != "ok" {
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
			if do_type == "watchlist" {
				delete_path = "/watch_list/" + tool.Url_parser(item.Data)
			} else if do_type == "bbs_watchlist" {
				view_path = "/bbs/w/" + tool.Url_parser(item.Set_id) + "/" + tool.Url_parser(item.Set_code)
				delete_path = "/bbs_watch/" + tool.Url_parser(item.Set_id) + "/" + tool.Url_parser(item.Set_code)
			}

			date_html := ""
			if item.Date != "" {
				date_html = "(" + tool.HTML_escape(item.Date) + ") "
			}
			data_html += `<li><a href="` + view_path + `">` + tool.HTML_escape(item.Display) + `</a> ` + date_html + `<a href="` + delete_path + `">(` + tool.Get_language(db, "delete", true) + `)</a></li>`
		}
		data_html += "</ul><hr class=\"main_hr\">"
	}

	page_url := "/star_doc_page/{}"
	if do_type == "watchlist" {
		page_url = "/watch_list_page/{}"
	} else if do_type == "bbs_watchlist" {
		page_url = "/bbs_watch_list_page/{}"
	}
	data_html += tool.Get_page_control(db, tool.Str_to_int(num), len(data_list), 50, page_url)

	manager_url := "/manager/16"
	if do_type == "watchlist" {
		manager_url = "/manager/13"
	} else if do_type == "bbs_watchlist" {
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
		},
		map[string]string{},
	)
}
