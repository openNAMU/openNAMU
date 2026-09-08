package route

import (
	"strings"

	"opennamu/route/tool"
)

func View_user_watch_list(config tool.Config, num string, do_type string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if tool.IP_or_user(config.IP) {
		return tool.Get_redirect("/login")
	}
	api_data := Api_user_watch_list(config, config.IP, num, do_type)
	if api_data["response"].(string) != "ok" {
		return tool.Get_error_page(db, config, "auth")
	}

	data_list := api_data["data"].([]string)
	data_html := ""

	if len(data_list) > 0 {
		data_html += "<ul>"
		for _, title := range data_list {
			display_title := title
			view_path := "/w/" + tool.Url_parser(title)
			delete_path := "/star_doc/" + tool.Url_parser(title)
			date_data := ""

			if do_type == "watchlist" {
				date_data = tool.Get_history_date(db, title)
				delete_path = "/watch_list/" + tool.Url_parser(title)
			} else if do_type == "thread_watchlist" {
				thread_title := ""
				thread_sub := ""
				tool.QueryRow_DB(
					db,
					"select title, sub, date from rd where code = ?",
					[]any{&thread_title, &thread_sub, &date_data},
					title,
				)
				if thread_title != "" {
					display_title = thread_title
					if thread_sub != "" {
						display_title += " - " + thread_sub
					}
				}
				view_path = "/thread/" + tool.Url_parser(title)
				delete_path = "/thread_watch/" + tool.Url_parser(title)
			} else if do_type == "bbs_watchlist" {
				watch_data := strings.SplitN(title, "-", 2)
				if len(watch_data) != 2 {
					continue
				}

				set_id := watch_data[0]
				set_code := watch_data[1]
				post_title, _ := tool.Get_bbs_data_value(db, set_id, set_code, "title")
				bbs_name := ""
				tool.QueryRow_DB(
					db,
					"select set_data from bbs_set where set_name = 'bbs_name' and set_id = ?",
					[]any{&bbs_name},
					set_id,
				)
				if post_title != "" {
					display_title = post_title
					if bbs_name != "" {
						display_title = bbs_name + " - " + display_title
					}
				}
				tool.QueryRow_DB(
					db,
					"select set_data from bbs_data where set_name = 'date' and set_id = ? and set_code = ?",
					[]any{&date_data},
					set_id,
					set_code,
				)
				view_path = "/bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code)
				delete_path = "/bbs_watch/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code)
			}

			date_html := ""
			if date_data != "" {
				date_html = "(" + tool.HTML_escape(date_data) + ") "
			}
			data_html += `<li><a href="` + view_path + `">` + tool.HTML_escape(display_title) + `</a> ` + date_html + `<a href="` + delete_path + `">(` + tool.Get_language(db, "delete", true) + `)</a></li>`
		}
		data_html += "</ul><hr class=\"main_hr\">"
	}

	page_url := "/star_doc_page/{}"
	if do_type == "watchlist" {
		page_url = "/watch_list_page/{}"
	} else if do_type == "thread_watchlist" {
		page_url = "/thread_watch_list_page/{}"
	} else if do_type == "bbs_watchlist" {
		page_url = "/bbs_watch_list_page/{}"
	}
	data_html += tool.Get_page_control(db, tool.Str_to_int(num), len(data_list), 50, page_url)

	manager_url := "/manager/16"
	if do_type == "watchlist" {
		manager_url = "/manager/13"
	} else if do_type == "thread_watchlist" || do_type == "bbs_watchlist" {
		manager_url = ""
	}
	if manager_url != "" {
		data_html += `<a href="` + manager_url + `">(` + tool.Get_language(db, "add", true) + `)</a>`
	}

	title := tool.Get_language(db, "watchlist", true)
	if do_type == "star_doc" {
		title = tool.Get_language(db, "star_doc", true)
	} else if do_type == "thread_watchlist" {
		title = tool.Get_language(db, "thread_watchlist", true)
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
			{"thread_watch_list", tool.Get_language(db, "thread_watchlist", false)},
			{"bbs_watch_list", tool.Get_language(db, "bbs_watchlist", false)},
		},
		map[string]string{},
	)
}
