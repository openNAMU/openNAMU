package route

import (
	"database/sql"
	"regexp"
	"strconv"
	"strings"

	"opennamu/route/tool"
)

func View_bbs_in_w_comment(db *sql.DB, config tool.Config, set_id string, set_code string, selected_comment string) string {
	data_api := Api_bbs_w_comment(config, "around", set_id+"-"+set_code)
	data_api_in := data_api["data"].([]map[string]string)

	bbs_comment_acl := tool.Check_acl(db, set_id, "", "bbs_comment", config.IP)

	select_html := `
        <select id="opennamu_comment_select" name="comment_select">
            <option value="0">` + tool.Get_language(db, "normal", true) + `</option>
    `
	data_html := ""

	tabom_count_api := Api_bbs_w_tabom(config, set_id, set_code)
	tabom_count := tabom_count_api["data"].(string)

	if bbs_comment_acl {
		data_html += `
            <hr class="main_hr">
            <form method="post" action="/bbs/w/` + tool.Url_parser(set_id) + `/` + tool.Url_parser(set_code) + `/tabom">
                <button type="submit">` + tool.Get_language(db, "upvote", true) + ` ` + tool.HTML_escape(tabom_count) + `</button>
            </form>
        `
	}

	data_html += "<hr>"

	var re = regexp.MustCompile(`^[0-9]+-[0-9]+-`)

	for _, v := range data_api_in {
		if v["comment"] == "" {
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
		if v["comment_user_id"] == config.IP {
			color = "green"
		}

		date += `<a href="/bbs/w/` + tool.Url_parser(set_id) + `/` + tool.Url_parser(set_code) + `/comment/` + tool.Url_parser(code_id) + `#opennamu_comment_select">(` + tool.Get_language(db, "comment", true) + `)</a> `
		date += `<a href="/bbs/tool/` + tool.Url_parser(set_id) + `/` + tool.Url_parser(set_code) + `/` + tool.Url_parser(code_id) + `">(` + tool.Get_language(db, "tool", true) + `)</a> `
		date += v["comment_date"]

		padding_str := strconv.Itoa(20 * count)

		data_html += `<span style="padding-left: ` + padding_str + `px;"></span>`
		rendered_data := Get_bbs_render(db, set_id, v["comment"], "thread", config)
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

	if bbs_comment_acl {
		data_html += `
            <form method="post" action="/bbs/w/` + tool.Url_parser(set_id) + `/` + tool.Url_parser(set_code) + `">
                <div id="opennamu_bbs_w_post_select">` + select_html + `</div>
                ` + tool.Get_editor_ui(db, config, "", "bbs_comment", "", "") + `
            </form>
        `
	}

	return data_html
}
