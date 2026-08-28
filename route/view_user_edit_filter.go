package route

import (
	"net/url"
	"regexp"

	"opennamu/route/tool"
)

func View_user_edit_filter(config tool.Config, user_name string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if user_name == "" {
		user_name = config.IP
	}
	if user_name != config.IP && !tool.Check_acl(db, "", "", "edit_filter_auth", config.IP) {
		return tool.Get_redirect("/recent_block")
	}
	if values != nil {
		if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
			return tool.Get_error_page(db, config, "auth")
		}
		tool.Exec_DB(db, "delete from user_set where name = 'edit_filter' and id = ?", user_name)
		return tool.Get_redirect("/edit_filter/" + tool.Url_parser(user_name))
	}

	filter_data := ""
	tool.QueryRow_DB(db, "select data from user_set where name = 'edit_filter' and id = ?", []any{&filter_data}, user_name)
	body := `<a href="/filter/edit_filter">(` + tool.Get_language(db, "edit_filter_rule", true) + `)</a><hr class="main_hr">`
	body += `<textarea readonly class="opennamu_textarea_500">` + tool.HTML_escape(filter_data) + `</textarea><ul>`
	rows := tool.Query_DB(db, "select plus from html_filter where kind = 'regex_filter' and plus != ''")
	for rows.Next() {
		pattern := ""
		if rows.Scan(&pattern) != nil {
			continue
		}
		re, err := regexp.Compile("(?i)" + pattern)
		if err == nil {
			if match := re.FindString(filter_data); match != "" {
				body += `<li>` + tool.HTML_escape(match) + `</li>`
			}
		}
	}
	rows.Close()
	body += `</ul>`
	if tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		body += `<hr class="main_hr"><form method="post"><button type="submit">` + tool.Get_language(db, "delete", true) + `</button></form>`
	}
	return tool.Get_template(db, config, user_name, body, []any{tool.Get_language(db, "edit_filter", true)}, [][]any{{"recent_block", tool.Get_language(db, "return", true)}}, map[string]string{})
}
