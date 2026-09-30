package route

import "opennamu/route/tool"

func View_list_old_page(config tool.Config, num string, set_type string, filter string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	api_data := Api_list_old_page(config, num, set_type, filter)
	if api_data["response"] != "ok" {
		return tool.Get_error_page(db, config, "error")
	}
	api_data_list, _ := api_data["data"].([][]string)

	if set_type != "old" {
		set_type = "new"
	}

	title := tool.Get_language(db, "old_page", true)
	if set_type == "new" {
		title = tool.Get_language(db, "new_page", true)
	}

	data_html := ""

	if set_type == "old" {
		all_label := tool.Get_language(db, "all", true)
		normal_label := tool.Get_language(db, "normal", true)

		var all_url, normal_url string
		if filter == "normal" {
			// normal 탭 활성: pagination은 현재 num 유지; 전체로 전환(필터 변경)이면 1페이지 초기화
			all_url = "/list/document/" + set_type + "/" + "1"
			normal_url = "/list/document/" + set_type + "/" + num + "/normal"
		} else {
			// all 탭 활성: pagination은 현재 num 유지; 일반으로 전환(필터 변경)이면 1페이지 초기화
			all_url = "/list/document/" + set_type + "/" + num
			normal_url = "/list/document/" + set_type + "/" + "1" + "/normal"
		}

		data_html += `<a href="` + all_url + `">(` + all_label + `)</a> `
		data_html += `<a href="` + normal_url + `">(` + normal_label + `)</a> `
	}

	for _, data := range api_data_list {
		doc_name := tool.Url_parser(data[0])
		doc_title := tool.HTML_escape(data[0])
		if data[2] != "" {
			doc_title += " (" + tool.Get_language(db, "redirect", false) + ")"
		}
		date := tool.HTML_escape(data[1])

		left := `<a href="/w/` + doc_name + `">` + doc_title + `</a>`
		right := date

		data_html += tool.Get_list_ui(left, right, "", "")
	}
	if len(api_data_list) == 0 {
		data_html += tool.Get_language(db, "data_missing", true)
	}

	page_url := "/list/document/" + set_type + "/{}"
	if set_type == "old" && filter == "normal" {
		page_url = "/list/document/old/{}/normal"
	}

	data_html += tool.Get_page_control(
		db,
		tool.Str_to_int(num),
		len(api_data_list),
		50,
		page_url,
	)

	out := tool.Get_template(
		db,
		config,
		title,
		data_html,
		[]any{},
		[][]any{
			{"other", tool.Get_language(db, "return", true)},
		},
		map[string]string{},
	)

	return out
}
