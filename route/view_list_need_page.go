package route

import (
	"strings"

	"opennamu/route/tool"
)

func View_list_need_page(config tool.Config, page string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	page_num := list_extra_page_number(page)
	offset := (page_num - 1) * 50
	rows := tool.Query_DB(db, "select title, link from back where type = 'no' order by title limit ?, 50", offset)
	body := strings.Builder{}
	count := 0
	for rows.Next() {
		name, link := "", ""
		if rows.Scan(&name, &link) == nil {
			body.WriteString(tool.Get_list_ui(`<a href="/w/`+tool.Url_parser(name)+`">`+tool.HTML_escape(name)+`</a>`, tool.HTML_escape(link), "", ""))
			count++
		}
	}
	rows.Close()
	body.WriteString(tool.Get_page_control(db, page_num, count, 50, "/list/document/need/{}"))
	return list_extra_page(db, config, tool.Get_language(db, "need_document", true), body.String())
}
