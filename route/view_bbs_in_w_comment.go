package route

import (
	"database/sql"
	"regexp"
	"strconv"
	"strings"

	"opennamu/route/tool"
)

func View_bbs_in_w_comment(db *sql.DB, config tool.Config, set_id string, set_code string, selected_comment string, page int) string {
	data_api := Api_bbs_w_comment(config, "around", set_id+"-"+set_code)
	data_api_in := data_api["data"].([]map[string]string)
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

	var re = regexp.MustCompile(`^[0-9]+-[0-9]+-`)
	comment_path := "/bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code) + "/comment/"
	if page > 1 {
		comment_path = "/bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code) + "/page/" + strconv.Itoa(page) + "/comment/"
	}

	for _, v := range data_api_in {
		comment_data := v["comment"]
		if v["blind"] == "O" && !comment_manage {
			comment_data = ""
		}
		if comment_data == "" && v["blind"] != "O" {
			continue
		}

		code_id := v["id"] + "-" + v["code"]
		code_id = re.ReplaceAllString(code_id, "")

		count := strings.Count(code_id, "-")

		selected := ""
		if selected_comment == code_id {
			selected = ` selected`
		}
		select_html += `<option value="` + tool.HTML_escape(code_id) + `"` + selected + `>` + tool.HTML_escape(code_id) + `</option>`

		color := "default"
		date := ""
		if v["code"] == "1" {
			color = "red"
		} else if v["comment_user_id"] == config.IP {
			color = "green"
		}

		date += `<a href="` + comment_path + tool.Url_parser(code_id) + `#opennamu_comment_select">(` + tool.Get_language(db, "comment", true) + `)</a> `
		date += `<a href="/bbs/tool/` + tool.Url_parser(set_id) + `/` + tool.Url_parser(set_code) + `/` + tool.Url_parser(code_id) + `">(` + tool.Get_language(db, "tool", true) + `)</a> `
		date += v["comment_date"]

		padding_str := strconv.Itoa(20 * count)

		data_html += `<span style="padding-left: ` + padding_str + `px;"></span>`
		rendered_data := Get_bbs_render(db, set_id, comment_data, "thread", config)
		rendered_data = render_topic_reference(rendered_data, set_code, set_id, set_code, "bbs")
		data_html += get_thread_ui_with_render(
			db,
			v["comment_user_id_render"],
			date,
			rendered_data,
			code_id,
			color,
			"",
			`width: calc(100% - `+padding_str+`px);`,
			set_code,
		)
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
                ` + tool.Get_editor_ui(db, config, "", "bbs_comment", "", "") + `
            </form>
        `
	}

	data_html += tool.Get_page_control(db, page, page_count, 50, "/bbs/w/"+tool.Url_parser(set_id)+"/"+tool.Url_parser(set_code)+"/page/{}")
	return data_html
}
