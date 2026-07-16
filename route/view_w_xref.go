package route

import (
	"opennamu/route/tool"
)

func View_w_xref(config tool.Config, doc_name string, do_type string, num string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	api_data := Api_w_xref(config, num, doc_name, do_type)
	response := api_data["response"].(string)
	if response == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	} else if response != "ok" {
		return tool.Get_error_page(db, config, "error")
	}

	data_list := api_data["data"].([][]string)
	data_html := ""
	data_sub := ""
	page_url := ""

	if do_type == "1" {
		data_html += `<a href="/xref_this/` + tool.Url_parser(doc_name) + `">(` + tool.Get_language(db, "link_in_this", true) + `)</a><hr class="main_hr">`
		data_sub = "(" + tool.Get_language(db, "backlink", true) + ")"
		page_url = "/xref_page/{}/" + tool.Url_parser(doc_name)
	} else {
		data_html += `<a href="/xref/` + tool.Url_parser(doc_name) + `">(` + tool.Get_language(db, "normal", true) + `)</a><hr class="main_hr">`
		data_sub = "(" + tool.Get_language(db, "link_in_this", true) + ")"
		page_url = "/xref_this_page/{}/" + tool.Url_parser(doc_name)
	}

	data_html += "<ul>"
	if do_type != "1" {
		link_count := api_data["link_count"].(string)
		if link_count == "" {
			link_count = tool.Get_language(db, "data_missing", true)
		} else {
			link_count = tool.HTML_escape(link_count)
		}

		data_html += "<li>" + tool.Get_language(db, "link_count", true) + " : " + link_count + "</li>"
	}

	for _, data := range data_list {
		data_html += `<li><a href="/w/` + tool.Url_parser(data[0]) + `">` + tool.HTML_escape(data[0]) + `</a>`

		if data[1] != "" {
			data_html += " (" + tool.HTML_escape(data[1]) + ")"
		}

		if len(data) > 2 && data[2] != "" {
			data_html += ` <a class="opennamu_link_inter" href="/xref/` + tool.Url_parser(data[0]) + `">(` + tool.Get_language(db, "backlink", true) + `)</a>`
		}

		data_html += "</li>"
	}
	data_html += "</ul>"
	data_html += tool.Get_page_control(
		db,
		tool.Str_to_int(num),
		len(data_list),
		50,
		page_url,
	)

	return tool.Get_template(
		db,
		config,
		doc_name,
		data_html,
		[]any{data_sub},
		[][]any{
			{"w/" + tool.Url_parser(doc_name), tool.Get_language(db, "return", true)},
			{"xref_reset/" + tool.Url_parser(doc_name), tool.Get_language(db, "reset_backlink", true)},
		},
		map[string]string{},
	)
}
