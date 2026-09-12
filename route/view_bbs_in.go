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

func View_bbs_in_filter_post(set_id string, comment_min string, commented string, comment_user string, tabom_min string, mine string, participate string, tabom_user string, author string, prefix string, tag string) string {
	commented_state := 0
	if commented == "1" {
		commented_state = 1
	} else if commented == "0" {
		commented_state = 2
	}
	filter := bbs_filter{
		comment_min:  bbs_filter_number(comment_min),
		commented:    commented_state,
		comment_user: strings.TrimSpace(comment_user),
		tabom_min:    bbs_filter_number(tabom_min),
		mine:         mine == "1",
		participate:  participate == "1",
		tabom_user:   tabom_user == "1",
		author:       strings.TrimSpace(author),
		prefix:       strings.TrimSpace(prefix),
		tag:          strings.TrimSpace(tag),
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
		mine_checked := ""
		if filter.mine {
			mine_checked = " checked"
		}
		participate_checked := ""
		if filter.participate {
			participate_checked = " checked"
		}
		tabom_user_checked := ""
		if filter.tabom_user {
			tabom_user_checked = " checked"
		}
		commented_html := `<label>` + tool.Get_language(db, "bbs_comment_status", true) + ` <select name="commented"><option value="">` + tool.Get_language(db, "all", true) + `</option>`
		if filter.commented == 1 {
			commented_html += `<option value="1" selected>` + tool.Get_language(db, "bbs_has_comment", true) + `</option>`
		} else {
			commented_html += `<option value="1">` + tool.Get_language(db, "bbs_has_comment", true) + `</option>`
		}
		if filter.commented == 2 {
			commented_html += `<option value="0" selected>` + tool.Get_language(db, "bbs_no_comment", true) + `</option>`
		} else {
			commented_html += `<option value="0">` + tool.Get_language(db, "bbs_no_comment", true) + `</option>`
		}
		commented_html += `</select></label>`
		prefix_html := `<label>` + tool.Get_language(db, "bbs_prefix", true) + ` <select name="prefix"><option value="">` + tool.Get_language(db, "all", true) + `</option>`
		for _, prefix := range bbs_prefix_list(db, set_id) {
			selected := ""
			if prefix == filter.prefix {
				selected = " selected"
			}
			prefix_html += `<option value="` + tool.HTML_escape(prefix) + `"` + selected + `>` + tool.HTML_escape(prefix) + `</option>`
		}
		prefix_html += `</select></label>`
		data_html += `<form method="post" action="/bbs/in/` + tool.Url_parser(set_id) + `/filter">
        <div><label>` + tool.Get_language(db, "comment", true) + ` <input name="comment_min" value="` + strconv.Itoa(filter.comment_min) + `"></label></div><hr class="main_hr">
        <div>` + commented_html + `</div><hr class="main_hr">
        <div><label>` + tool.Get_language(db, "upvote", true) + ` <input name="tabom_min" value="` + strconv.Itoa(filter.tabom_min) + `"></label></div><hr class="main_hr">
        <div><label><input type="checkbox" name="mine" value="1"` + mine_checked + `>` + tool.Get_language(db, "my_bbs_post", true) + `</label></div><hr class="main_hr">
        <div><label><input type="checkbox" name="participate" value="1"` + participate_checked + `>` + tool.Get_language(db, "participate_bbs_post", true) + `</label></div><hr class="main_hr">
        <div><label><input type="checkbox" name="tabom_user" value="1"` + tabom_user_checked + `>` + tool.Get_language(db, "my_tabom_bbs_post", true) + `</label></div><hr class="main_hr">
        <div><label>` + tool.Get_language(db, "bbs_comment_author", true) + ` <input name="comment_user" value="` + tool.HTML_escape(filter.comment_user) + `"></label></div><hr class="main_hr">
        <div><label>` + tool.Get_language(db, "bbs_author", true) + ` <input name="author" value="` + tool.HTML_escape(filter.author) + `"></label></div><hr class="main_hr">
        <div>` + prefix_html + `</div><hr class="main_hr">
        <div><label>` + tool.Get_language(db, "tag", true) + ` <input name="tag" value="` + tool.HTML_escape(filter.tag) + `"></label></div>
        <div><button class="__ON_BUTTON__" type="submit">` + tool.Get_language(db, "filter", true) + `</button> <a href="/bbs/in/` + tool.Url_parser(set_id) + `/1">(` + tool.Get_language(db, "reset", true) + `)</a></div>
    </form><hr class="main_hr">`
	}
	data_html += bbs_list_example_ui(db)
	data_html += Get_bbs_list_ui(db, config, data_api_in, map[string]string{})
	page_path := "/bbs/in/" + tool.Url_parser(set_id) + "/{}"
	if show_filter {
		page_path = "/bbs/in/" + tool.Url_parser(set_id) + "/filter/"
		if filter_path != "" {
			page_path += filter_path + "/"
		}
		page_path += "{}"
	} else if sort_type == "activity" {
		page_path = "/bbs/in/" + tool.Url_parser(set_id) + "/activity/{}"
	} else if sort_type == "view" {
		page_path = "/bbs/in/" + tool.Url_parser(set_id) + "/view/{}"
	} else if sort_type == "comment" {
		page_path = "/bbs/in/" + tool.Url_parser(set_id) + "/comment/{}"
	} else if sort_type == "tabom" {
		page_path = "/bbs/in/" + tool.Url_parser(set_id) + "/tabom/{}"
	} else if sort_type == "excellent" {
		page_path = "/bbs/in/" + tool.Url_parser(set_id) + "/excellent/{}"
	}
	data_html += tool.Get_page_control(
		db,
		tool.Str_to_int(page_num),
		len(data_api_in),
		50,
		page_path,
	)

	sort_path := "bbs/in/" + tool.Url_parser(set_id) + "/view/1"
	sort_name := tool.Get_language(db, "page_view_sort", true)
	if sort_type == "view" {
		sort_path = "bbs/in/" + tool.Url_parser(set_id) + "/1"
		sort_name = tool.Get_language(db, "recent", true)
	}
	filter_menu_path := "bbs/in/" + tool.Url_parser(set_id) + "/filter/"
	if filter_path != "" {
		filter_menu_path += filter_path + "/"
	}
	filter_menu_path += "1"
	data_html += `<hr class="main_hr">`
	data_html += `(<a href="/bbs/in/` + tool.Url_parser(set_id) + `/activity/1">` + tool.Get_language(db, "activity_sort", false) + `</a>)&nbsp;`
	data_html += `(<a href="/bbs/in/` + tool.Url_parser(set_id) + `/comment/1">` + tool.Get_language(db, "comment_sort", false) + `</a>)&nbsp;`
	data_html += `(<a href="/bbs/in/` + tool.Url_parser(set_id) + `/tabom/1">` + tool.Get_language(db, "upvote_sort", false) + `</a>)&nbsp;`
	data_html += `(<a href="/bbs/in/` + tool.Url_parser(set_id) + `/excellent/1">` + tool.Get_language(db, "excellent_post", false) + `</a>)&nbsp;`
	data_html += `(<a href="/bbs/in/` + tool.Url_parser(set_id) + `/filter/mine/1">` + tool.Get_language(db, "my_bbs_post", false) + `</a>)&nbsp;`
	data_html += `(<a href="/bbs/in/` + tool.Url_parser(set_id) + `/filter/participate/1">` + tool.Get_language(db, "participate_bbs_post", false) + `</a>)&nbsp;`
	data_html += `(<a href="/bbs/in/` + tool.Url_parser(set_id) + `/filter/tabom_user/1">` + tool.Get_language(db, "my_tabom_bbs_post", false) + `</a>)&nbsp;`
	data_html += `(<a href="/` + filter_menu_path + `">` + tool.Get_language(db, "filter", false) + `</a>)&nbsp;`
	data_html += `(<a href="/bbs/set/` + tool.Url_parser(set_id) + `">` + tool.Get_language(db, "bbs_set", false) + `</a>)`

	menu := [][]any{
		{"bbs/main", tool.Get_language(db, "return", true)},
		{"bbs/search/" + tool.Url_parser(set_id), tool.Get_language(db, "search", true)},
	}
	if set_id != "0" {
		add_path := "bbs/edit/" + tool.Url_parser(set_id)
		if set_id == thread_bbs_id && filter.tag != "" {
			add_path += "/document/" + tool.Base64_encode(filter.tag)
		}
		menu = append(menu, []any{add_path, tool.Get_language(db, "add", true)})
	}
	menu = append(menu,
		[]any{sort_path, sort_name},
	)

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
