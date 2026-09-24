package route

import (
	"strings"

	"opennamu/route/tool"
)

func View_list_need_page(config tool.Config, page string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	page_num := List_extra_page_number(page)
	offset := (page_num - 1) * 50
	rows := tool.Get_need_document_rows(db, offset)
	body := strings.Builder{}
	count := 0
	has_next := false
	for rows.Next() {
		if count == 50 {
			has_next = true
			break
		}
		name, link_count := "", ""
		if rows.Scan(&name, &link_count) == nil {
			body.WriteString(tool.Get_list_ui(`<a href="/w/`+tool.Url_parser(name)+`">`+tool.HTML_escape(name)+`</a>`, tool.Get_language(db, "link_count", true)+" : "+tool.HTML_escape(link_count), "", ""))
			count++
		}
	}
	rows.Close()
	if count == 0 {
		body.WriteString(tool.Get_language(db, "data_missing", true))
	}
	body.WriteString(tool.Get_page_control(db, page_num, count, 50, "/list/document/need/{}", has_next))
	return List_extra_page(db, config, tool.Get_language(db, "need_document", true), body.String())
}
