package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func user_value(db *sql.DB, id string, name string) string {
	value := ""
	tool.QueryRow_DB(db, "select data from user_set where id = ? and name = ?", []any{&value}, id, name)
	return value
}

func user_save(db *sql.DB, id string, name string, value string) {
	existing_id := ""
	if tool.QueryRow_DB(db, "select id from user_set where id = ? and name = ?", []any{&existing_id}, id, name) {
		tool.Exec_DB(db, "update user_set set data = ? where id = ? and name = ?", value, id, name)
		return
	}
	tool.Exec_DB(db, "insert into user_set (id, name, data) values (?, ?, ?)", id, name, value)
}

func user_delete(db *sql.DB, id string, name string) {
	tool.Exec_DB(db, "delete from user_set where id = ? and name = ?", id, name)
}

func user_form_page(db *sql.DB, config tool.Config, title string, body string) string {
	return tool.Get_template(db, config, title, body, []any{}, [][]any{{"user", tool.Get_language(db, "return", true)}}, map[string]string{})
}

func user_auth(db *sql.DB, config tool.Config) bool {
	return !tool.IP_or_user(config.IP)
}
func View_user_safe(config tool.Config, user_name string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	is_self := user_name == ""
	if is_self {
		user_name = config.IP
	}

	body := `<h2>` + tool.Get_language(db, "state", true) + `</h2><div id="opennamu_get_user_info">` + tool.HTML_escape(user_name) + `</div>`
	if is_self {
		alarm_count := "0"
		tool.QueryRow_DB(db, "select count(*) from user_notice where name = ? and readme = ''", []any{&alarm_count}, config.IP)
		alarm_text := tool.Get_language(db, "alarm", true)
		if tool.Str_to_int(alarm_count) > 0 {
			alarm_text += " (" + tool.HTML_escape(alarm_count) + ")"
		}

		login_menu := `<li><a href="/logout">` + tool.Get_language(db, "logout", true) + `</a></li><li><a href="/change">` + tool.Get_language(db, "user_setting", true) + `</a></li>`
		tool_menu := `<li><a href="/alarm">` + alarm_text + `</a></li>`
		if tool.IP_or_user(config.IP) {
			login_menu = `<li><a href="/login">` + tool.Get_language(db, "login", true) + `</a></li><li><a href="/register">` + tool.Get_language(db, "register", true) + `</a></li><li><a href="/change">` + tool.Get_language(db, "user_setting", true) + `</a></li><li><a href="/login/find">` + tool.Get_language(db, "password_search", true) + `</a></li>`
		} else {
			tool_menu += `<li><a href="/watch_list">` + tool.Get_language(db, "watchlist", true) + `</a></li><li><a href="/star_doc">` + tool.Get_language(db, "star_doc", true) + `</a></li><li><a href="/challenge">` + tool.Get_language(db, "challenge_and_level_manage", true) + `</a></li><li><a href="/acl/user:` + tool.Url_parser(config.IP) + `">` + tool.Get_language(db, "user_document_acl", true) + `</a></li>`
		}
		body += `<h2>` + tool.Get_language(db, "login", true) + `</h2><ul>` + login_menu + `</ul><h2>` + tool.Get_language(db, "tool", true) + `</h2><ul>` + tool_menu + `</ul>`
	}

	body += `<h2>` + tool.Get_language(db, "other", true) + `</h2><ul>` +
		`<li><a href="/record/` + tool.Url_parser(user_name) + `">` + tool.Get_language(db, "edit_record", true) + `</a></li>` +
		`<li><a href="/record/topic/` + tool.Url_parser(user_name) + `">` + tool.Get_language(db, "discussion_record", true) + `</a></li>` +
		`<li><a href="/record/bbs/` + tool.Url_parser(user_name) + `">` + tool.Get_language(db, "bbs_record", true) + `</a></li>` +
		`<li><a href="/record/bbs_comment/` + tool.Url_parser(user_name) + `">` + tool.Get_language(db, "bbs_comment_record", true) + `</a></li>` +
		`<li><a href="/topic/user:` + tool.Url_parser(user_name) + `">` + tool.Get_language(db, "user_discussion", true) + `</a></li>` +
		`<li><a href="/count/` + tool.Url_parser(user_name) + `">` + tool.Get_language(db, "count", true) + `</a></li></ul>`

	if tool.Check_acl(db, "", "", "give_auth", config.IP) {
		ban_name := tool.Get_language(db, "ban", true)
		user_auth := tool.Get_user_auth(db, user_name)
		if tool.Auth_group_name_ban(user_auth) {
			ban_name = tool.Get_language(db, "release", true)
		}
		body += `<h2>` + tool.Get_language(db, "admin", true) + `</h2><ul><li><a href="/auth/give/` + tool.Url_parser(user_name) + `">` + ban_name + `</a></li><li><a href="/list/user/check_submit/` + tool.Url_parser(user_name) + `">` + tool.Get_language(db, "check", true) + `</a></li></ul>`
	}

	return user_form_page(db, config, tool.Get_language(db, "user_tool", true), body)
}
