package tool

import (
	"database/sql"
	"strings"
)

func Get_user_info_ui(db *sql.DB, config Config, user_name string) string {
	auth_name := Get_user_auth(db, user_name)
	auth_date := Get_auth_date(db, user_name)
	if auth_date != "0" {
		auth_name += " (~" + auth_date + ")"
	}

	ban_state := Get_language(db, "normal", false)
	if Auth_group_name_ban(auth_name) {
		ban_state = "<a href=\"/recent_block\">" + Get_language(db, "ban", false) + "</a>"
	}

	level_data := Get_level(db, user_name)
	return `<table class="user_info_table"><tr><td>` + Get_language(db, "user_name", false) + `</td><td>` + IP_parser(db, user_name, config.IP) + `</td></tr><tr><td>` + Get_language(db, "authority", false) + `</td><td>` + HTML_escape(auth_name) + `</td></tr><tr><td>` + Get_language(db, "state", false) + `</td><td>` + ban_state + `</td></tr><tr><td>` + Get_language(db, "level", false) + `</td><td>` + HTML_escape(level_data[0]) + ` (` + HTML_escape(level_data[1]) + ` / ` + HTML_escape(level_data[2]) + `)</td></tr></table>`
}

func Replace_user_info_ui(db *sql.DB, config Config, data string) string {
	const start_tag = `<div id="opennamu_get_user_info">`
	const end_tag = `</div>`

	for {
		start_index := strings.Index(data, start_tag)
		if start_index < 0 {
			return data
		}
		content_start := start_index + len(start_tag)
		end_offset := strings.Index(data[content_start:], end_tag)
		if end_offset < 0 {
			return data
		}
		end_index := content_start + end_offset
		user_name := HTML_unescape(data[content_start:end_index])
		data = data[:start_index] + Get_user_info_ui(db, config, user_name) + data[end_index+len(end_tag):]
	}
}
