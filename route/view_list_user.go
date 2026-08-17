package route

import (
	"strings"

	"opennamu/route/tool"
)

func View_list_user(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	body := strings.Builder{}
	rows := tool.Query_DB(db, "select id, data from user_set where name = 'date' order by data desc")
	for rows.Next() {
		name, date := "", ""
		if rows.Scan(&name, &date) == nil {
			body.WriteString(tool.Get_list_ui(`<a href="/user/`+tool.Url_parser(name)+`">`+tool.IP_parser(db, name, config.IP)+`</a>`, tool.HTML_escape(date), "", ""))
		}
	}
	rows.Close()
	return list_extra_page(db, config, tool.Get_language(db, "user_list", true), body.String())
}
