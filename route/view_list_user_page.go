package route

import (
	"strings"

	"opennamu/route/tool"
)

func View_list_user_page(config tool.Config, page string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !tool.Check_permission(db, "user_manage", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	page_num := list_extra_page_number(page)
	offset := (page_num - 1) * 50
	rows := tool.Get_user_date_rows(db, offset, true)
	body := strings.Builder{}
	count := 0
	for rows.Next() {
		name, date := "", ""
		if rows.Scan(&name, &date) == nil {
			body.WriteString(tool.Get_list_ui(`<a href="/user/`+tool.Url_parser(name)+`">`+tool.IP_parser(db, name, config.IP)+`</a>`, tool.HTML_escape(date), "", ""))
			count++
		}
	}
	rows.Close()
	body.WriteString(tool.Get_page_control(db, page_num, count, 50, "/list/user/{}"))
	return list_extra_page(db, config, tool.Get_language(db, "member_list", true), body.String())
}
