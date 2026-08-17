package route

import (
	"opennamu/route/tool"

	"github.com/gin-gonic/gin"
)

func View_bbs_in_w(c *gin.Context, config tool.Config, set_id string, set_code string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, set_id, "", "bbs_view", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	bbs_name := Api_bbs_num_to_name(db, set_id)["data"].(string)

	data_api := Api_bbs_w(config, set_id, set_code)
	data_api_in := data_api["data"].(map[string]string)

	if len(data_api_in) == 0 {
		return tool.Get_redirect("/bbs/main")
	}

	post_data := get_render_setting_css(db, config) + Get_bbs_render(db, set_id, data_api_in["data"], "bbs", config)
	post_data = render_topic_reference(post_data, set_code, set_id, set_code, "bbs")
	post_title := tool.HTML_escape(data_api_in["title"])
	if data_api_in["prefix"] != "" {
		post_title = "[" + tool.HTML_escape(data_api_in["prefix"]) + "] " + post_title
	}

	tag_html := ""
	for _, tag := range bbs_tag_list(data_api_in["tags"]) {
		if tag_html != "" {
			tag_html += ", "
		}
		tag_html += "#" + tool.HTML_escape(tag)
	}
	if tag_html != "" {
		tag_html = "<div>" + tag_html + "</div>"
	}

	data_html := `
        <div class="opennamu_bbs_w_post_tab">
            <big><big><big>` + post_title + `</big></big></big>
            ` + tag_html + `
            <hr class="main_hr">
            ` + data_api_in["user_id_render"] + ` <span style="float: right;">` + data_api_in["date"] + `</span>
            <hr>
            <div class="opennamu_bbs_w_post_tab_content">
                ` + post_data + `
            </div>
        </div>
    `

	Api_bbs_w_page_view_post(config, set_id, set_code)

	view_count_api := Api_bbs_w_page_view(config, set_id, set_code)
	view_count_api_data := view_count_api["data"].(int)

	data_html += View_bbs_in_w_comment(db, config, set_id, set_code, c.Query("comment_select"))

	out := tool.Get_template(
		db,
		config,
		bbs_name,
		data_html,
		[]any{"(" + tool.Get_language(db, "bbs", true) + ")", data_api_in["date"], 0, 0, view_count_api_data},
		[][]any{
			{"bbs/in/" + tool.Url_parser(set_id), tool.Get_language(db, "return", true)},
			{"bbs/edit/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code), tool.Get_language(db, "edit", true)},
			{"bbs/tool/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code), tool.Get_language(db, "tool", true)},
		},
		map[string]string{
			"path": c.Request.URL.Path,
		},
	)

	return out
}
