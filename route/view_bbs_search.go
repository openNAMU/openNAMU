package route

import (
	"strconv"

	"opennamu/route/tool"
)

func View_bbs_search(config tool.Config, set_id string, keyword string, page string) string {
	return View_bbs_search_internal(config, set_id, keyword, page, "title")
}

func View_bbs_search_data(config tool.Config, set_id string, keyword string, page string) string {
	return View_bbs_search_internal(config, set_id, keyword, page, "data")
}

func View_bbs_search_internal(config tool.Config, set_id string, keyword string, page string, search_type string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if set_id == "" && !tool.Check_permission(db, "bbs_main_view", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	if set_id != "" && !tool.Check_acl(db, set_id, "", "bbs_view", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	page_int := tool.Str_to_int(page)
	if page_int < 1 {
		page_int = 1
	}

	bbs_name := ""
	title := tool.Get_language(db, "search", true)
	search_path := "/bbs/search"
	search_switch_path := "/bbs/search_data"
	search_switch_text := tool.Get_language(db, "bbs_search_data", true)
	comment_search_path := "/bbs/search_comment"
	comment_search_text := tool.Get_language(db, "comment", true) + " " + tool.Get_language(db, "search", true)
	bbs_id_to_name := map[string]string{}
	if search_type == "data" {
		title = tool.Get_language(db, "bbs_search_data", true)
		search_path = "/bbs/search_data"
		search_switch_path = "/bbs/search"
		search_switch_text = tool.Get_language(db, "search", true)
	}

	if set_id != "" {
		bbs_name_data := Api_bbs_num_to_name(db, set_id)
		bbs_name, _ = bbs_name_data["data"].(string)
		if bbs_name == "" {
			return tool.Get_redirect("/bbs/main")
		}

		title = bbs_name
		search_path += "/" + tool.Url_parser(set_id)
		search_switch_path += "/" + tool.Url_parser(set_id)
		comment_search_path += "/" + tool.Url_parser(set_id)
		bbs_id_to_name[set_id] = bbs_name
	} else {
		for name, id := range Bbs_list(db) {
			if tool.Check_acl(db, id, "", "bbs_view", config.IP) {
				bbs_id_to_name[id] = name
			}
		}
	}

	data_html := `<form method="post" action="` + search_path + `">
        <div><label for="bbs_search_keyword">` + tool.Get_language(db, "search", true) + `</label></div>
        <div><input id="bbs_search_keyword" class="__ON_INPUT__" name="keyword" value="` + tool.HTML_escape(keyword) + `"></div><hr class="main_hr">
        <div><button class="__ON_BUTTON__" type="submit">` + tool.Get_language(db, "search", true) + `</button></div>
    </form><hr class="main_hr">
    <div>(<a href="` + search_switch_path + `">` + search_switch_text + `</a>)</div><hr class="main_hr">`
	data_html += `<div>(<a href="` + comment_search_path + `">` + comment_search_text + `</a>)</div><hr class="main_hr">`

	if keyword != "" {
		data_api := Api_bbs_search(config, keyword, set_id, strconv.Itoa(page_int))
		if search_type == "data" {
			data_api = Api_bbs_search_data(config, keyword, set_id, strconv.Itoa(page_int))
		}
		data_list, _ := data_api["data"].([]map[string]string)
		data_html += Get_bbs_list_ui(db, config, data_list, bbs_id_to_name)

		page_url := "/bbs/search_page/{}/" + tool.Url_parser(keyword)
		if search_type == "data" {
			page_url = "/bbs/search_data_page/{}/" + tool.Url_parser(keyword)
		}
		has_next, _ := data_api["has_next"].(bool)
		if len(data_list) == 0 {
			data_html += `<div>` + tool.Get_language(db, "search_no_result", true) + `</div><hr class="main_hr">`
		}
		if set_id != "" {
			page_url = "/bbs/search_board_page/" + tool.Url_parser(set_id) + "/{}" + "/" + tool.Url_parser(keyword)
			if search_type == "data" {
				page_url = "/bbs/search_data_board_page/" + tool.Url_parser(set_id) + "/{}" + "/" + tool.Url_parser(keyword)
			}
		}
		data_html += tool.Get_page_control(db, page_int, len(data_list), 50, page_url, has_next)
	}

	return_menu := "bbs/main"
	if set_id != "" {
		return_menu = "bbs/in/" + tool.Url_parser(set_id)
	}

	return tool.Get_template(
		db,
		config,
		title,
		data_html,
		[]any{},
		[][]any{
			{return_menu, tool.Get_language(db, "return", true)},
		},
		map[string]string{},
	)
}
