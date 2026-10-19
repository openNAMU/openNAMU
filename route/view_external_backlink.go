package route

import (
	"strconv"
	"strings"

	"opennamu/route/tool"
)

func View_external_docs_for_host(config tool.Config, host string, page_str string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	api_data := Api_external_docs_for_host(config, page_str, host)
	response, _ := api_data["response"].(string)
	if response == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	} else if response != "ok" {
		return tool.Get_error_page(db, config, "error")
	}

	data_list, _ := api_data["data"].([][]string)
	data_html := ""
	page_url := ""

	data_html += `<a href="/ext_backlink/` + tool.Url_parser(strings.TrimSpace(host)) + `">(` + tool.Get_language(db, "backlink", true) + `)</a><hr class="main_hr">`
	page_url = "/ext_backlink_page/{}/" + tool.Url_parser(strings.TrimSpace(host))

	data_html += "<ul>"
	for _, data := range data_list {
		doc_name := data[0]
		host_text := data[1]
		data_html += `<li><a href="/w/` + tool.Url_parser(doc_name) + `">` + tool.HTML_escape(doc_name) + `</a> ` + tool.HTML_escape(host_text) + `</li>`
	}
	data_html += "</ul>"
	data_html += tool.Get_page_control(
		db,
		tool.Str_to_int(page_str),
		len(data_list),
		50,
		page_url,
	)

	return tool.Get_template(
		db,
		config,
		host,
		data_html,
		[]any{},
		[][]any{
			{"ext_backlink/" + tool.Url_parser(strings.TrimSpace(host)), tool.Get_language(db, "backlink", true)},
		},
		map[string]string{},
	)
}

func View_external_hosts_of_doc(config tool.Config, doc_name string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	api_data := Api_external_hosts_of_doc(config, doc_name)
	response, _ := api_data["response"].(string)
	if response == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	} else if response != "ok" {
		return tool.Get_error_page(db, config, "error")
	}

	hosts, _ := api_data["data"].([]string)
	data_html := ""

	data_html += `<a href="/w/` + tool.Url_parser(doc_name) + `">(` + tool.Get_language(db, "normal", true) + `)</a><hr class="main_hr">`
	data_html += "<ul>"
	for _, h := range hosts {
		data_html += `<li>` + tool.HTML_escape(h) + `</li>`
	}
	data_html += "</ul>"

	return tool.Get_template(
		db,
		config,
		doc_name,
		data_html,
		[]any{},
		[][]any{
			{"w/" + tool.Url_parser(doc_name), tool.Get_language(db, "normal", true)},
		},
		map[string]string{},
	)
}

func View_external_count_for_host(config tool.Config, host string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	api_data := Api_external_count_for_host(config, host)
	response, _ := api_data["response"].(string)
	if response != "ok" {
		return tool.Get_error_page(db, config, "error")
	}

	count, _ := api_data["count"].(int)
	data_html := ""

	data_html += `<a href="/ext_backlink/` + tool.Url_parser(strings.TrimSpace(host)) + `">(` + tool.Get_language(db, "backlink", true) + `)</a><hr class="main_hr">`
	data_html += "<ul>"
	data_html += "<li>" + tool.Get_language(db, "link_count", true) + " : " + strconv.Itoa(count) + "</li>"
	data_html += "</ul>"

	return tool.Get_template(
		db,
		config,
		host,
		data_html,
		[]any{},
		[][]any{
			{"ext_backlink/" + tool.Url_parser(strings.TrimSpace(host)), tool.Get_language(db, "backlink", true)},
		},
		map[string]string{},
	)
}

func View_external_all_hosts(config tool.Config, page_str string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	api_data := Api_external_all_hosts(config, page_str)
	response, _ := api_data["response"].(string)
	if response == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	} else if response != "ok" {
		return tool.Get_error_page(db, config, "error")
	}

	data_list, _ := api_data["data"].([][]string)
	data_html := ""
	page_url := "/ext_backlink_hosts/{}"

	data_html += `<a href="/ext_backlink_hosts">(` + tool.Get_language(db, "backlink", true) + `)</a><hr class="main_hr">`
	data_html += "<ul>"
	for _, data := range data_list {
		host := data[0]
		count := data[1]
		data_html += `<li><a href="/ext_backlink/` + tool.Url_parser(host) + `">` + tool.HTML_escape(host) + `</a> ` + tool.Get_language(db, "link_count", true) + ` : ` + count + `</li>`
	}
	data_html += "</ul>"
	data_html += tool.Get_page_control(
		db,
		tool.Str_to_int(page_str),
		len(data_list),
		50,
		page_url,
	)

	return tool.Get_template(
		db,
		config,
		"ext_backlink_hosts",
		data_html,
		[]any{},
		[][]any{},
		map[string]string{},
	)
}
