package route

import (
	"strconv"
	"strings"

	"opennamu/route/tool"
)

func View_bbs_in(config tool.Config, set_id string, page_num string, sort_type string) string {
	return view_bbs_in(config, set_id, page_num, sort_type, bbs_filter{}, "", false)
}

func View_bbs_in_filter(config tool.Config, set_id string, filter_data string) string {
	page_num, filter_data := bbs_filter_path_data(filter_data)
	filter := bbs_filter_parse(filter_data)
	filter_path := bbs_filter_path(filter)
	return view_bbs_in(config, set_id, page_num, "", filter, filter_path, true)
}

func View_bbs_in_filter_post(set_id string, comment_min string, tabom_min string, tag string) string {
	filter := bbs_filter{
		comment_min: bbs_filter_number(comment_min),
		tabom_min:   bbs_filter_number(tabom_min),
		tag:         strings.TrimSpace(tag),
	}
	filter_path := bbs_filter_path(filter)
	target := "/bbs/in/" + tool.Url_parser(set_id) + "/filter/"
	if filter_path != "" {
		target += filter_path + "/"
	}
	target += "1"

	return tool.Get_redirect(target)
}

func view_bbs_in(config tool.Config, set_id string, page_num string, sort_type string, filter bbs_filter, filter_path string, show_filter bool) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, set_id, "", "bbs_view", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	bbs_name := Api_bbs_num_to_name(db, set_id)["data"].(string)
	if bbs_name == "" {
		return tool.Get_redirect("/bbs/main")
	}

	var data_api map[string]any
	if show_filter {
		data_api = api_bbs(config, set_id, page_num, sort_type, filter)
	} else {
		data_api = Api_bbs(config, set_id, page_num, sort_type)
	}
	data_api_in := data_api["data"].([]map[string]string)

	data_html := ""
	if show_filter {
		data_html += `<form method="post" action="/bbs/in/` + tool.Url_parser(set_id) + `/filter">
        <label>` + tool.Get_language(db, "comment", true) + ` <input name="comment_min" value="` + strconv.Itoa(filter.comment_min) + `"></label>
        <label>` + tool.Get_language(db, "upvote", true) + ` <input name="tabom_min" value="` + strconv.Itoa(filter.tabom_min) + `"></label>
        <label>` + tool.Get_language(db, "tag", true) + ` <input name="tag" value="` + tool.HTML_escape(filter.tag) + `"></label>
        <button class="__ON_BUTTON__" type="submit">` + tool.Get_language(db, "filter", true) + `</button>
    </form><hr class="main_hr">`
	}
	data_html += Get_bbs_list_ui(config, data_api_in, map[string]string{})
	page_path := "/bbs/in/" + tool.Url_parser(set_id) + "/{}"
	if show_filter {
		page_path = "/bbs/in/" + tool.Url_parser(set_id) + "/filter/"
		if filter_path != "" {
			page_path += filter_path + "/"
		}
		page_path += "{}"
	} else if sort_type == "view" {
		page_path = "/bbs/in/" + tool.Url_parser(set_id) + "/view/{}"
	} else if sort_type == "comment" {
		page_path = "/bbs/in/" + tool.Url_parser(set_id) + "/comment/{}"
	} else if sort_type == "tabom" {
		page_path = "/bbs/in/" + tool.Url_parser(set_id) + "/tabom/{}"
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
	filter_menu_path := "/bbs/in/" + tool.Url_parser(set_id) + "/filter/"
	if filter_path != "" {
		filter_menu_path += filter_path + "/"
	}
	filter_menu_path += "1"

	menu := [][]any{
		{"bbs/main", tool.Get_language(db, "return", true)},
		{"bbs/search/" + tool.Url_parser(set_id), tool.Get_language(db, "search", true)},
		{"bbs/edit/" + tool.Url_parser(set_id), tool.Get_language(db, "add", true)},
		{sort_path, sort_name},
		{"bbs/in/" + tool.Url_parser(set_id) + "/comment/1", tool.Get_language(db, "comment_sort", true)},
		{"bbs/in/" + tool.Url_parser(set_id) + "/tabom/1", tool.Get_language(db, "upvote_sort", true)},
		{filter_menu_path, tool.Get_language(db, "filter", true)},
		{"bbs/set/" + tool.Url_parser(set_id), tool.Get_language(db, "bbs_set", true)},
	}

	out := tool.Get_template(
		db,
		config,
		bbs_name,
		data_html,
		[]any{},
		menu,
		map[string]string{},
	)

	return out
}
