package route

import (
	"database/sql"
	"opennamu/route/tool"
	"regexp"
	"strconv"
	"strings"
)

func View_bbs_in_w_comment(db *sql.DB, config tool.Config, user_name string, set_id string, set_code string) string {
    data_api := Api_bbs_w_comment(config, "around", set_id + "-" + set_code)
    data_api_in := data_api["data"].([]map[string]string)

    bbs_comment_acl := tool.Check_acl(db, set_id, "", "bbs_comment", config.IP)

    select_html := `
        <select id="opennamu_comment_select" name="comment_select">
            <option value="0">` + tool.Get_language(db, "normal", true) + `</option>
    `
    data_html := ""

    tabom_count_api := Api_bbs_w_tabom(config, set_id, set_code)
    tabom_count := tabom_count_api["data"]

    if bbs_comment_acl {
        data_html += `
            <hr class="main_hr">
            <div id="opennamu_bbs_w_post_tabom">
                <a href="javascript:void(0);" id="opennamu_tabom_button">
                    <span class="opennamu_bbs_w_post_tabom opennamu_svg opennamu_svg_tabom">&nbsp;</span>
                </a>
                <script>
                    window.addEventListener('DOMContentLoaded', function() {
                        document.getElementById('opennamu_tabom_button').addEventListener("click", function() {
                            opennamu_post_tabom("` + tool.JS_escape(set_id) + `", "` + tool.JS_escape(set_code) + `");
                        });
                    });
                </script>
                <hr class="main_hr">
                <span>` + tool.Get_language(db, "upvote", true) + `</span> <span class="opennamu_tabom_count">` + tabom_count + `</span>
            </div>
        `
    }

    data_html += "<hr>"

    var re = regexp.MustCompile(`^[0-9]+-[0-9]+-`)

    for _, v := range data_api_in {
        code_id := v["id"] + "-" + v["code"]
        code_id = re.ReplaceAllString(code_id, "")

        count := strings.Count(code_id, "-")

        select_html += `<option value="` + code_id + `">` + code_id + `</option>`

        color := "default"
        date := ""

        date += `<a href="javascript:opennamu_change_comment('` + code_id + `');">(` + tool.Get_language(db, "comment", true) + `)</a> `;
        date += `<a href="/bbs/tool/` + set_id + `/` + set_code + `/` + code_id + `">(` + tool.Get_language(db, "tool", true) + `)</a> `;
        date += v["comment_date"];

        padding_str := strconv.Itoa(20 * count)

        data_html += `<span style="padding-left: ` + padding_str + `px;"></span>`
        data_html += tool.Get_thread_ui(
            v["comment_user_id_render"],
            date,
            v["comment"],
            code_id,
            color,
            "",
            `width: calc(100% - ` + padding_str + `px);`,
            "",
        )
    }

    select_html += `</select> <a href="javascript:opennamu_return_comment();">(` + tool.Get_language(db, "return", true) + `)</a>`
    select_html += `<hr class="main_hr">`;

    if bbs_comment_acl {
        data_html += `
            <form method="post" action="/bbs/w/` + tool.Url_parser(set_id) + `/` + tool.Url_parser(set_code) + `">
                <div id="opennamu_bbs_w_post_select">` + select_html + `</div>
                ` + tool.Get_editor_ui(db, config, "", "bbs_comment", "", "") + `
            </form>
        `
    }

    data_html += `
        <script defer src="/views/main_css/js/route/topic.js` + tool.Cache_v() + `"></script>
        <script defer src="/views/main_css/js/route/bbs_w_post.js` + tool.Cache_v() + `"></script>
    `

    return data_html
}