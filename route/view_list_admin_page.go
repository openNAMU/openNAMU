package route

import (
	"strings"

	"opennamu/route/tool"
)

func View_list_admin_page(config tool.Config, page string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	page_num := list_extra_page_number(page)
	offset := (page_num - 1) * 50
	body := strings.Builder{}
	count := 0
	for _, auth_data := range tool.Get_auth_user_list(db, offset, 50) {
		name, auth := auth_data[0], auth_data[1]
		body.WriteString(tool.Get_list_ui(`<a href="/user/`+tool.Url_parser(name)+`">`+tool.IP_parser(db, name, config.IP)+`</a>`, tool.HTML_escape(auth), "", ""))
		count++
	}
	body.WriteString(tool.Get_page_control(db, page_num, count, 50, "/list/admin/{}"))
	return list_extra_page(db, config, tool.Get_language(db, "admin_list", true), body.String())
}
