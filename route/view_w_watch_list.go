package route

import (
	"opennamu/route/tool"
)

func View_w_watch_list(config tool.Config, doc_name string, num string, do_type string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	api_data := Api_w_watch_list(config, doc_name, num, do_type)
	data_html := ""

	if do_type != "watchlist" {
		do_type = "star_doc"
	}

	if api_data["response"] != "ok" {
		return tool.Get_error_page(db, config, "auth")
	} else {
		data_html += "<ul>"
		data_list := api_data["data"].([][]string)
		for _, user_data := range data_list {
			data_html += "<li>" + user_data[1] + "</li>"
		}

		data_html += "</ul>"
		route_name := "doc_watch_list"
		if do_type == "star_doc" {
			route_name = "doc_star_doc"
		}
		data_html += tool.Get_page_control(
			db,
			tool.Str_to_int(num),
			len(data_list),
			50,
			"/"+route_name+"/{}/"+tool.Url_parser(doc_name),
		)
	}

	title := tool.Get_language(db, "watch_user_list", true)
	if do_type == "star_doc" {
		title = tool.Get_language(db, "star_user_list", true)
	}

	out := tool.Get_template(
		db,
		config,
		title,
		data_html,
		[]any{"(" + doc_name + ")"},
		[][]any{
			{"w/" + tool.Url_parser(doc_name), tool.Get_language(db, "return", false)},
			{"doc_watch_list/1/" + tool.Url_parser(doc_name), tool.Get_language(db, "watch_user_list", false)},
			{"doc_star_doc/1/" + tool.Url_parser(doc_name), tool.Get_language(db, "star_user_list", false)},
		},
		map[string]string{},
	)

	return out
}
