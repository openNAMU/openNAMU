package route

import (
	"database/sql"
	"strconv"
	"strings"

	"opennamu/route/tool"
)

func View_bbs_in_w_comment(db *sql.DB, config tool.Config, set_id string, set_code string, selected_comment string, page int) string {
	data_api := Api_bbs_w_comment(config, "around", set_id+"-"+set_code)
	data_api_in := data_api["data"].([]map[string]string)
	post_user, _ := tool.Get_bbs_data_value(db, set_id, set_code, "user_id")
	if post_user != "" {
		post_user = tool.IP_preprocess(db, post_user, config.IP)[0]
	}
	all_data_api_in := data_api_in
	comment_prefix := set_id + "-" + set_code + "-"
	pinned_data_api_in := []map[string]string{}
	for _, v := range all_data_api_in {
		if v["pinned"] == "" || (v["comment"] == "" && v["blind"] != "O") {
			continue
		}
		code_id := strings.TrimPrefix(v["id"]+"-"+v["code"], comment_prefix)
		if bbs_comment_code_regex.MatchString(code_id) {
			pinned_data_api_in = append(pinned_data_api_in, v)
		}
	}
	if page < 1 {
		page = 1
	}
	start := (page - 1) * 50
	end := start + 50
	if start >= len(data_api_in) {
		data_api_in = []map[string]string{}
	} else {
		if end > len(data_api_in) {
			end = len(data_api_in)
		}
		data_api_in = data_api_in[start:end]
	}
	page_count := len(data_api_in)

	bbs_comment_acl := tool.Check_acl(db, set_id, "", "bbs_comment", config.IP)
	comment_manage := tool.Check_permission(db, "bbs_comment_manage", config.IP)
	comment_closed := bbs_comment_closed(db, set_id, set_code)
	can_comment := bbs_comment_acl && !comment_closed
	bbs_comment_placeholder := bbs_set_value(db, set_id, "bbs_comment_placeholder")

	select_html := `
        <select id="opennamu_comment_select" name="comment_select">
            <option value="0">` + tool.Get_language(db, "normal", true) + `</option>
    `
	data_html := ""
	if set_id == "0" && tool.Check_permission(db, "bbs_comment_manage", config.IP) {
		data_html += `<a href="/bbs/tool/` + tool.Url_parser(set_id) + `/` + tool.Url_parser(set_code) + `">(` + tool.Get_language(db, "comment_manage", true) + `)</a>`
	}

	tabom_count_api := Api_bbs_w_tabom(config, set_id, set_code)
	tabom_count := tabom_count_api["data"].(string)
	tabom_down_count := tabom_count_api["down_data"].(string)

	if bbs_comment_acl {
		data_html += `
            <hr class="main_hr">
            <form method="post" action="/bbs/w/` + tool.Url_parser(set_id) + `/` + tool.Url_parser(set_code) + `/tabom" style="display: inline">
                <input type="hidden" name="vote_type" value="up">
                <button type="submit">` + tool.Get_language(db, "upvote", true) + ` ` + tool.HTML_escape(tabom_count) + `</button>
            </form>
            <form method="post" action="/bbs/w/` + tool.Url_parser(set_id) + `/` + tool.Url_parser(set_code) + `/tabom" style="display: inline">
                <input type="hidden" name="vote_type" value="down">
                <button type="submit">` + tool.Get_language(db, "downvote", true) + ` ` + tool.HTML_escape(tabom_down_count) + `</button>
            </form>
        `
	}

	data_html += "<hr>"
	if comment_closed {
		data_html += `<span>` + tool.Get_language(db, "comment_closed", true) + `</span><hr>`
	}

	comment_path := "/bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code) + "/comment/"
	if page > 1 {
		comment_path = "/bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code) + "/page/" + strconv.Itoa(page) + "/comment/"
	}

	if page == 1 && len(pinned_data_api_in) > 0 {
		for _, v := range pinned_data_api_in {
			comment_html, _, exists := get_bbs_comment_ui(db, config, post_user, set_id, set_code, comment_prefix, comment_path, v, comment_manage, bbs_comment_acl, true)
			if exists {
				data_html += comment_html
			}
		}
		data_html += `<hr class="main_hr">`
	}

	for _, v := range data_api_in {
		comment_html, code_id, exists := get_bbs_comment_ui(db, config, post_user, set_id, set_code, comment_prefix, comment_path, v, comment_manage, bbs_comment_acl, false)
		if !exists {
			continue
		}

		selected := ""
		if selected_comment == code_id {
			selected = ` selected`
		}
		if v["blind"] != "O" || comment_manage {
			select_html += `<option value="` + tool.HTML_escape(code_id) + `"` + selected + `>` + tool.HTML_escape(code_id) + `</option>`
		}
		data_html += comment_html
	}

	return_anchor := "opennamu_comment_select"
	if selected_comment != "" {
		return_anchor = selected_comment
	}
	select_html += `</select> <a href="#` + tool.Url_parser(return_anchor) + `">(` + tool.Get_language(db, "return", true) + `)</a>`
	select_html += `<hr class="main_hr">`

	if can_comment {
		data_html += `
            <form method="post" action="/bbs/w/` + tool.Url_parser(set_id) + `/` + tool.Url_parser(set_code) + `">
                <div id="opennamu_bbs_w_post_select">` + select_html + `</div>
                ` + tool.Get_editor_ui(db, config, "", "bbs_comment", "", "", bbs_comment_placeholder) + `
            </form>
        `
	}

	data_html += tool.Get_page_control(db, page, page_count, 50, "/bbs/w/"+tool.Url_parser(set_id)+"/"+tool.Url_parser(set_code)+"/page/{}")
	return data_html
}
func get_bbs_comment_ui(db *sql.DB, config tool.Config, post_user string, set_id string, set_code string, comment_prefix string, comment_path string, v map[string]string, comment_manage bool, comment_acl bool, copy_comment bool) (string, string, bool) {
	comment_data := v["comment"]
	if v["blind"] == "O" && !comment_manage {
		comment_data = ""
	}
	if comment_data == "" && v["blind"] != "O" {
		return "", "", false
	}

	code_id := strings.TrimPrefix(v["id"]+"-"+v["code"], comment_prefix)
	if !bbs_comment_code_regex.MatchString(code_id) {
		return "", "", false
	}

	color := "default"
	if v["pinned"] != "" {
		color = "red"
	} else if post_user != "" && v["comment_user_id"] == post_user {
		color = "blue"
	} else if v["comment_user_id"] == config.IP {
		color = "green"
	}

	date := `<a href="` + comment_path + tool.Url_parser(code_id) + `#opennamu_comment_select">(` + tool.Get_language(db, "comment", true) + `)</a> `
	date += `<a href="/bbs/tool/` + tool.Url_parser(set_id) + `/` + tool.Url_parser(set_code) + `/` + tool.Url_parser(code_id) + `">(` + tool.Get_language(db, "tool", true) + `)</a> `
	date += v["comment_date"]
	if comment_acl {
		tabom_count := v["tabom_count"]
		if tabom_count == "" {
			tabom_count = "0"
		}
		tabom_down_count := v["tabom_down_count"]
		if tabom_down_count == "" {
			tabom_down_count = "0"
		}
		date += ` <a href="/bbs/w/` + tool.Url_parser(set_id) + `/` + tool.Url_parser(set_code) + `/comment_tabom/` + tool.Url_parser(code_id) + `/up">(` + tool.Get_language(db, "upvote", true) + ` ` + tool.HTML_escape(tabom_count) + `)</a>`
		date += ` <a href="/bbs/w/` + tool.Url_parser(set_id) + `/` + tool.Url_parser(set_code) + `/comment_tabom/` + tool.Url_parser(code_id) + `/down">(` + tool.Get_language(db, "downvote", true) + ` ` + tool.HTML_escape(tabom_down_count) + `)</a>`
	}

	padding_str := "0"
	if !copy_comment {
		padding_str = strconv.Itoa(20 * strings.Count(code_id, "-"))
	}

	rendered_data := Get_bbs_render(db, set_id, comment_data, "thread", config)
	rendered_data = render_topic_reference(rendered_data, set_code, set_id, set_code, "bbs")
	comment_ui := get_thread_ui_with_render(
		db,
		v["comment_user_id_render"],
		date,
		rendered_data,
		code_id,
		color,
		v["blind"],
		`width: calc(100% - `+padding_str+`px);`,
		set_code,
	)
	if copy_comment {
		comment_ui = get_thread_ui_with_render_copy(
			db,
			v["comment_user_id_render"],
			date,
			rendered_data,
			code_id,
			color,
			v["blind"],
			`width: calc(100% - `+padding_str+`px);`,
			set_code,
		)
	}

	return `<span style="padding-left: ` + padding_str + `px;"></span>` + comment_ui, code_id, true
}
