package route

import "opennamu/route/tool"

func View_bbs_raw(config tool.Config, set_id string, set_code string, comment_code string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, set_id, "", "bbs_view", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	title := ""
	if !tool.QueryRow_DB(
		db,
		"select set_data from bbs_data where set_name = 'title' and set_id = ? and set_code = ?",
		[]any{&title},
		set_id,
		set_code,
	) {
		return tool.Get_redirect("/bbs/main")
	}

	user_id := ""
	if !tool.QueryRow_DB(
		db,
		"select set_data from bbs_data where set_name = 'user_id' and set_id = ? and set_code = ?",
		[]any{&user_id},
		set_id,
		set_code,
	) {
		return tool.Get_redirect("/bbs/main")
	}
	if !bbs_post_view_allowed(db, set_id, user_id, config.IP, nil) {
		return tool.Get_error_page(db, config, "auth")
	}

	raw_data := ""
	if comment_code == "" {
		if !tool.QueryRow_DB(
			db,
			"select set_data from bbs_data where set_name = 'data' and set_id = ? and set_code = ?",
			[]any{&raw_data},
			set_id,
			set_code,
		) {
			return tool.Get_redirect("/bbs/main")
		}
	} else {
		comment_data := Api_bbs_w_comment_one(
			config,
			true,
			"",
			set_id+"-"+set_code+"-"+comment_code,
		)
		comment_list, ok := comment_data["data"].([]map[string]string)
		if !ok || len(comment_list) == 0 {
			return tool.Get_redirect("/bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code))
		}

		raw_data = comment_list[0]["comment"]
	}

	sub := "(" + tool.Get_language(db, "raw", true) + ") (" + tool.Get_language(db, "bbs", true) + ")"
	if comment_code != "" {
		sub += " (" + tool.HTML_escape(comment_code) + ")"
	}

	data_html := `
        <div id="opennamu_preview_area">
            <textarea readonly id="opennamu_edit_textarea" class="opennamu_textarea_500 __ON_TEXTAREA__">` + tool.HTML_escape(raw_data) + `</textarea>
        </div>
    `

	return tool.Get_template(
		db,
		config,
		title,
		data_html,
		[]any{sub},
		[][]any{
			{"bbs/tool/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code), tool.Get_language(db, "return", true)},
		},
		map[string]string{},
	)
}
