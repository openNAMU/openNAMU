package route

import (
	"opennamu/route/tool"
	"strings"
)

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

	data_html := `<table id="main_table_set"><tr id="main_table_top_tr"><td>` + tool.Get_language(db, "editor", true) + `</td><td>` + tool.Get_language(db, "time", true) + `</td><td>` + tool.Get_language(db, "comment", true) + `</td></tr>`
	row_count := 0
	for rows.Next() {
		comment_set_id := ""
		comment_code := ""
		date := ""
		if rows.Scan(&comment_set_id, &comment_code, &date) != nil {
			continue
		}

		parts := strings.Split(comment_set_id, "-")
		if len(parts) < 2 {
			continue
		}
		row_count++

		bbs_id := parts[0]
		post_id := parts[1]
		comment_link := ""
		if len(parts) > 2 {
			comment_link = strings.Join(parts[2:], "-") + "-"
		}
		comment_link += comment_code

		comment_user := record_bbs_legacy_value(db, "comment_user_id", comment_set_id, comment_code)
		title := record_bbs_legacy_value(db, "title", bbs_id, post_id)
		bbs_name := record_bbs_legacy_board_name(db, bbs_id)
		title_link := `<a href="/bbs/w/` + tool.Url_parser(bbs_id) + `/` + tool.Url_parser(post_id) + `#` + tool.Url_parser(comment_link) + `">` + tool.HTML_escape(title) + `</a>`
		data_html += `<tr><td>` + tool.IP_parser(db, comment_user, config.IP) + `</td><td>` + tool.HTML_escape(date) + `</td><td>#` + tool.HTML_escape(comment_link) + `</td></tr>`
		data_html += `<tr><td colspan="3">` + title_link + ` (` + tool.HTML_escape(bbs_name) + `)</td></tr>`
	}
	data_html += `</table>`
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
