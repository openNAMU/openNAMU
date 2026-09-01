package route

import (
	"database/sql"
	"opennamu/route/tool"
)

func View_record_bbs_legacy(config tool.Config, user_name string, page string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	page_int := tool.Str_to_int(page)
	if page_int < 1 {
		page_int = 1
	}
	offset := (page_int - 1) * 50
	rows := tool.Get_bbs_record_rows(db, user_name, offset)
	defer rows.Close()

	data_html := `<table id="main_table_set"><tr id="main_table_top_tr"><td>` + tool.Get_language(db, "editor", true) + `</td><td>` + tool.Get_language(db, "time", true) + `</td><td>` + tool.Get_language(db, "last_comment_time", true) + `</td></tr>`
	row_count := 0
	for rows.Next() {
		set_code := ""
		set_id := ""
		date := ""
		if rows.Scan(&set_code, &set_id, &date) != nil {
			continue
		}
		row_count++

		post_user := record_bbs_legacy_value(db, "user_id", set_id, set_code)
		title := record_bbs_legacy_value(db, "title", set_id, set_code)
		comment_count := record_bbs_legacy_value(db, "comment_count", set_id, set_code)
		if comment_count == "" {
			comment_count = "0"
		}
		root_id := set_id + "-" + set_code
		last_comment_date := tool.Get_bbs_last_comment_date(db, root_id)

		bbs_name := record_bbs_legacy_board_name(db, set_id)
		title_link := `<a href="/bbs/w/` + tool.Url_parser(set_id) + `/` + tool.Url_parser(set_code) + `">` + tool.HTML_escape(title) + `</a>`
		data_html += `<tr><td>` + tool.IP_parser(db, post_user, config.IP) + `</td><td>` + tool.HTML_escape(date) + `</td><td>` + tool.HTML_escape(last_comment_date) + `</td></tr>`
		data_html += `<tr><td colspan="3">` + title_link + ` (` + tool.HTML_escape(comment_count) + `) (` + tool.HTML_escape(bbs_name) + `)</td></tr>`
	}
	data_html += `</table>`
	data_html += tool.Get_page_control(db, page_int, row_count, 50, "/record/bbs/"+tool.Url_parser(user_name)+"/{}")

	return tool.Get_template(
		db,
		config,
		user_name,
		data_html,
		[]any{"(" + tool.Get_language(db, "bbs_record", true) + ")"},
		[][]any{{"user/" + tool.Url_parser(user_name), tool.Get_language(db, "user_tool", true)}},
		map[string]string{},
	)
}

func record_bbs_legacy_value(db *sql.DB, set_name string, set_id string, set_code string) string {
	value, _ := tool.Get_bbs_data_value(db, set_id, set_code, set_name)
	return value
}

func record_bbs_legacy_board_name(db *sql.DB, set_id string) string {
	return tool.Get_bbs_set_data(db, set_id, "bbs_name")
}
