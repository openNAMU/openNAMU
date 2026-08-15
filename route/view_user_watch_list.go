package route

import "opennamu/route/tool"

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
			last_date := ""
			tool.QueryRow_DB(db, "select date from history where title = ? order by id + 0 desc limit 1", []any{&last_date}, title)
			date_html := ""
			if last_date != "" {
				date_html = "(" + tool.HTML_escape(last_date) + ") "
			}
			path := "/star_doc/"
			if do_type == "watchlist" {
				path = "/watch_list/"
			}
			data_html += `<li><a href="/w/` + tool.Url_parser(title) + `">` + tool.HTML_escape(title) + `</a> ` + date_html + `<a href="` + path + tool.Url_parser(title) + `">(` + tool.Get_language(db, "delete", true) + `)</a></li>`
		}
		data_html += "</ul><hr class=\"main_hr\">"
	}

	page_url := "/star_doc?num={}"
	if do_type == "watchlist" {
		page_url = "/watch_list?num={}"
	}
	data_html += tool.Get_page_control(db, tool.Str_to_int(num), len(data_list), 50, page_url)

	manager_url := "/manager/16"
	if do_type == "watchlist" {
		manager_url = "/manager/13"
	}
	data_html += `<a href="` + manager_url + `">(` + tool.Get_language(db, "add", true) + `)</a>`

	title := tool.Get_language(db, "watchlist", true)
	if do_type == "star_doc" {
		title = tool.Get_language(db, "star_doc", true)
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
		},
		map[string]string{},
	)
}
