package route

import "opennamu/route/tool"

func View_record_bbs_comment_legacy(config tool.Config, user_name string, page string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	page_int := tool.Str_to_int(page)
	if page_int < 1 {
		page_int = 1
	}
	offset := (page_int - 1) * 50
	rows := tool.Get_bbs_comment_record_rows(db, user_name, offset)
	defer rows.Close()

	data_html := ""
	row_count := 0
	for rows.Next() {
		comment_set_id := ""
		comment_code := ""
		date := ""
		if rows.Scan(&comment_set_id, &comment_code, &date) != nil {
			continue
		}

		bbs_id, post_id, comment_link, valid := Bbs_comment_storage_location(comment_set_id, comment_code)
		if !valid {
			continue
		}
		row_count++

		comment_user := Record_bbs_legacy_value(db, "comment_user_id", comment_set_id, comment_code)
		title := Record_bbs_legacy_value(db, "title", bbs_id, post_id)
		bbs_name := Record_bbs_legacy_board_name(db, bbs_id)
		title_link := `<a href="/bbs/w/` + tool.Url_parser(bbs_id) + `/` + tool.Url_parser(post_id) + `#` + tool.Url_parser(comment_link) + `">` + tool.HTML_escape(title) + `</a>`
		left := `<strong>` + tool.Get_language(db, "editor", true) + `:</strong> ` + tool.IP_parser(db, comment_user, config.IP)
		right := `<strong>` + tool.Get_language(db, "time", true) + `:</strong> ` + tool.HTML_escape(date)
		bottom := `<div><strong>` + tool.Get_language(db, "comment", true) + `:</strong> #` + tool.HTML_escape(comment_link) + `</div>`
		bottom += `<div>` + title_link + ` (` + tool.HTML_escape(bbs_name) + `)</div>`
		data_html += tool.Get_list_ui(left, right, bottom, "")
	}
	data_html += tool.Get_page_control(db, page_int, row_count, 50, "/record/bbs_comment/"+tool.Url_parser(user_name)+"/{}")

	return tool.Get_template(
		db,
		config,
		user_name,
		data_html,
		[]any{"(" + tool.Get_language(db, "bbs_comment_record", true) + ")"},
		[][]any{{"user/" + tool.Url_parser(user_name), tool.Get_language(db, "user_tool", true)}},
		map[string]string{},
	)
}
