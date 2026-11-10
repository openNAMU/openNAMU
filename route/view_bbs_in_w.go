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

	bbs_name, _ := Api_bbs_num_to_name(db, set_id)["data"].(string)

	data_api := Api_bbs_w(config, set_id, set_code)
	data_api_in, _ := data_api["data"].(map[string]string)

	if len(data_api_in) == 0 {
		return tool.Get_redirect("/bbs/main")
	}

	post_data := Get_render_setting_css(db, config) + Get_bbs_render(db, set_id, data_api_in["data"], "bbs", config)
	post_data = Render_topic_reference(post_data, set_code, set_id, set_code, "bbs")
	post_title := tool.HTML_escape(data_api_in["title"])
	if data_api_in["prefix"] != "" {
		post_title = "[" + tool.HTML_escape(data_api_in["prefix"]) + "] " + post_title
	}

	tag_html := ""
	for _, tag := range Bbs_tag_list(data_api_in["tags"]) {
		if tag_html != "" {
			tag_html += ", "
		}
		tag_html += "#" + tool.HTML_escape(tag)
		if !tool.IP_or_user(config.IP) && Bbs_tag_watch_value_valid(tag) {
			tag_watch_label := tool.Get_language(db, "bbs_tag_watch", true)
			if watch_data := Api_bbs_tag_watch_view(config, "", tag); watch_data["response"] == "ok" {
				if watch_info, ok := watch_data["data"].(map[string]string); ok && watch_info["watched"] == "true" {
					tag_watch_label = tool.Get_language(db, "bbs_tag_watch_remove", true)
				}
			}
			tag_html += ` (<a href="/bbs_tag_watch/` + tool.Base64_encode(tag) + `">` + tag_watch_label + `</a>)`
		}
	}
	if tag_html != "" {
		tag_html = "<div>" + tag_html + "</div>"
	}

	report_target_html := ""
	if set_id == report_bbs_id && data_api_in["report_target"] != "" {
		report_target_html = `<div>` + tool.Get_language(db, "report_target", true) + `: <a href="` + tool.HTML_escape(data_api_in["report_target"]) + `">` + tool.HTML_escape(data_api_in["report_target"]) + `</a></div><hr class="main_hr">`
	}

	data_html := `
        <div class="opennamu_bbs_w_post_tab">
            <big><big><big>` + post_title + `</big></big></big>
            ` + tag_html + `
            <hr class="main_hr">
            ` + report_target_html + `
            ` + data_api_in["user_id_render"] + ` <span style="float: right;">` + data_api_in["date"] + `</span>
            <hr class="main_hr">
            <div class="opennamu_bbs_w_post_tab_content">
                ` + post_data + `
            </div>
        </div>
    `

	Api_bbs_w_page_view_post(config, set_id, set_code)

	view_count_api := Api_bbs_w_page_view(config, set_id, set_code)
	view_count_api_data, _ := view_count_api["data"].(int)

	data_html += View_bbs_in_w_comment(db, config, set_id, set_code, c.Param("comment_select"), tool.Str_to_int(c.Param("page")))

	menu := [][]any{
		{"bbs/in/" + tool.Url_parser(set_id), tool.Get_language(db, "return", true)},
		{"bbs/tool/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code), tool.Get_language(db, "tool", true)},
	}
	if set_id != thread_bbs_id && set_id != report_bbs_id {
		menu = append(menu, []any{"bbs/edit/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code), tool.Get_language(db, "edit", true)})
	}
	if set_id != report_bbs_id && tool.Check_permission(db, "bbs_comment", config.IP) {
		menu = append(menu, []any{"bbs/report/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code), tool.Get_language(db, "report", true)})
	}
	if !tool.IP_or_user(config.IP) {
		comment_watch_label := tool.Get_language(db, "bbs_comment_watch", true)
		if watch_data := Api_bbs_watch_view(config, set_id, set_code); watch_data["response"] == "ok" {
			if watch_info, ok := watch_data["data"].(map[string]string); ok && watch_info["watched"] == "true" {
				comment_watch_label = tool.Get_language(db, "bbs_comment_watch_remove", true)
			}
		}
		bookmark_label := tool.Get_language(db, "bbs_post_bookmark", true)
		if bookmark_data := Api_bbs_bookmark_view(config, set_id, set_code); bookmark_data["response"] == "ok" {
			if bookmark_info, ok := bookmark_data["data"].(map[string]string); ok && bookmark_info["watched"] == "true" {
				bookmark_label = tool.Get_language(db, "bbs_post_bookmark_remove", true)
			}
		}
		menu = append(menu,
			[]any{"bbs_watch/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code), comment_watch_label},
			[]any{"bbs_bookmark/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code), bookmark_label},
		)
	}

	out := tool.Get_template(
		db,
		config,
		bbs_name,
		data_html,
		[]any{"(" + tool.Get_language(db, "bbs", true) + ")", data_api_in["date"], 0, 0, view_count_api_data},
		menu,
		map[string]string{
			"path": c.Request.URL.Path,
		},
	)

	return out
}
