package route

import (
	"strconv"

	"opennamu/route/tool"
)

func View_bbs_search(config tool.Config, set_id string, keyword string, page string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	page_int := tool.Str_to_int(page)
	if page_int < 1 {
		page_int = 1
	}

	bbs_name := ""
	title := tool.Get_language(db, "search", true)
	search_path := "/bbs/search"
	bbs_id_to_name := map[string]string{}

	if set_id != "" {
		bbs_name = Api_bbs_num_to_name(db, set_id)["data"].(string)
		if bbs_name == "" {
			return tool.Get_redirect("/bbs/main")
		}

		title = bbs_name
		search_path += "/" + tool.Url_parser(set_id)
		bbs_id_to_name[set_id] = bbs_name
	} else {
		for name, id := range bbs_list(db) {
			bbs_id_to_name[id] = name
		}
	}

	data_html := `<form method="post" action="` + search_path + `">
        <input class="__ON_INPUT__" name="keyword" value="` + tool.HTML_escape(keyword) + `" placeholder="` + tool.Get_language(db, "search", true) + `">
        <button class="__ON_BUTTON__" type="submit">` + tool.Get_language(db, "search", true) + `</button>
    </form><hr class="main_hr">`

	if keyword != "" {
		data_api := Api_bbs_search(config, keyword, set_id, strconv.Itoa(page_int))
		data_list := data_api["data"].([]map[string]string)
		data_html += Get_bbs_list_ui(config, data_list, bbs_id_to_name)

		page_url := "/bbs/search_page/{}/" + tool.Url_parser(keyword)
		if set_id != "" {
			page_url = "/bbs/search_board_page/" + tool.Url_parser(set_id) + "/{}" + "/" + tool.Url_parser(keyword)
		}
		data_html += tool.Get_page_control(db, page_int, len(data_list), 50, page_url)
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
