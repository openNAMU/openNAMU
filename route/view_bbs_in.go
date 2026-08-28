package route

import (
	"opennamu/route/tool"
)

func View_bbs_in(config tool.Config, set_id string, page_num string, sort_type string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, set_id, "", "bbs_view", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	bbs_name := Api_bbs_num_to_name(db, set_id)["data"].(string)
	if bbs_name == "" {
		return tool.Get_redirect("/bbs/main")
	}

	data_api := Api_bbs(config, set_id, page_num, sort_type)
	data_api_in := data_api["data"].([]map[string]string)

	data_html := Get_bbs_list_ui(config, data_api_in, map[string]string{})
	page_path := "/bbs/in/" + tool.Url_parser(set_id) + "/{}"
	if sort_type == "view" {
		page_path = "/bbs/in/" + tool.Url_parser(set_id) + "/view/{}"
	}
	data_html += tool.Get_page_control(
		db,
		tool.Str_to_int(page_num),
		len(data_api_in),
		50,
		page_path,
	)

	sort_path := "/bbs/in/" + tool.Url_parser(set_id) + "/view/1"
	sort_name := tool.Get_language(db, "page_view_sort", true)
	if sort_type == "view" {
		sort_path = "/bbs/in/" + tool.Url_parser(set_id) + "/1"
		sort_name = tool.Get_language(db, "recent", true)
	}

	out := tool.Get_template(
		db,
		config,
		bbs_name,
		data_html,
		[]any{},
		[][]any{
			{"bbs/main", tool.Get_language(db, "return", true)},
			{"bbs/search/" + tool.Url_parser(set_id), tool.Get_language(db, "search", true)},
			{"bbs/edit/" + tool.Url_parser(set_id), tool.Get_language(db, "add", true)},
			{sort_path, sort_name},
			{"bbs/set/" + tool.Url_parser(set_id), tool.Get_language(db, "bbs_set", true)},
		},
		map[string]string{},
	)

	return out
}
