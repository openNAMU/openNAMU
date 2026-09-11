package route

import (
	"strings"

	"opennamu/route/tool"
)

func View_list_no_category_page(config tool.Config, page string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	page_num := list_extra_page_number(page)
	offset := (page_num - 1) * 50
	rows := tool.Get_no_category_document_rows(db, offset)
	body := strings.Builder{}
	body.WriteString(`<ul>`)
	count := 0
	for rows.Next() {
		name := ""
		if rows.Scan(&name) != nil {
			continue
		}
		body.WriteString(`<li><a href="/w/` + tool.Url_parser(name) + `">` + tool.HTML_escape(name) + `</a></li>`)
		count++
	}
	rows.Close()
	body.WriteString(`</ul>`)
	body.WriteString(tool.Get_page_control(db, page_num, count, 50, "/list/document/no_category/{}"))

	return list_extra_page(db, config, tool.Get_language(db, "no_category_document", true), body.String())
}
