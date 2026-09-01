package route

import (
	"strings"

	"opennamu/route/tool"
)

func View_list_no_link_page(config tool.Config, page string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	page_num := list_extra_page_number(page)
	offset := (page_num - 1) * 50
	rows := tool.Get_no_link_page_rows(db, offset)
	body := strings.Builder{}
	count := 0
	for rows.Next() {
		name, value := "", ""
		if rows.Scan(&name, &value) == nil {
			body.WriteString(tool.Get_list_ui(`<a href="/w/`+tool.Url_parser(name)+`">`+tool.HTML_escape(name)+`</a>`, tool.HTML_escape(value), "", ""))
			count++
		}
	}
	rows.Close()
	body.WriteString(tool.Get_page_control(db, page_num, count, 50, "/list/document/no_link/{}"))
	return list_extra_page(db, config, tool.Get_language(db, "no_link_document", true), body.String())
}
