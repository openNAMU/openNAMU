package route

import (
	"database/sql"
	"strconv"
	"strings"

	"opennamu/route/tool"
)

func View_list_admin(config tool.Config, auth_use bool) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if auth_use && !tool.Check_permission(db, "auth_group_manage", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	body := strings.Builder{}
	if auth_use {
		rows := tool.Query_DB(db, "select who, what from re_admin order by time desc limit 200")
		for rows.Next() {
			who, what := "", ""
			if rows.Scan(&who, &what) == nil {
				body.WriteString(tool.Get_list_ui(tool.IP_parser(db, who, config.IP), tool.HTML_escape(what), "", ""))
			}
		}
		rows.Close()
		return list_extra_page(db, config, tool.Get_language(db, "auth_use", true), body.String())
	}
	for _, auth_data := range tool.Get_auth_user_list(db, 0, 0) {
		name, auth := auth_data[0], auth_data[1]
		body.WriteString(tool.Get_list_ui(`<a href="/user/`+tool.Url_parser(name)+`">`+tool.IP_parser(db, name, config.IP)+`</a>`, tool.HTML_escape(auth), "", ""))
	}
	return list_extra_page(db, config, tool.Get_language(db, "admin_list", true), body.String())
}

func list_extra_page(db *sql.DB, config tool.Config, title string, body string) string {
	return tool.Get_template(db, config, title, body, []any{}, [][]any{{"other", tool.Get_language(db, "return", true)}}, map[string]string{})
}

func list_extra_query(db *sql.DB, config tool.Config, title string, query string, args ...any) string {
	rows := tool.Query_DB(db, query, args...)
	body := strings.Builder{}
	for rows.Next() {
		name, value := "", ""
		if rows.Scan(&name, &value) != nil {
			continue
		}
		body.WriteString(tool.Get_list_ui(`<a href="/w/`+tool.Url_parser(name)+`">`+tool.HTML_escape(name)+`</a>`, tool.HTML_escape(value), "", ""))
	}
	rows.Close()
	return list_extra_page(db, config, title, body.String())
}
func list_page_path(base string, page int) string {
	return base + "/" + strconv.Itoa(page)
}
