package route

import (
	"strings"

	"opennamu/route/tool"
)

func View_list_document_acl(config tool.Config, page string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "document_acl_manage", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	page_num := list_extra_page_number(page)
	offset := (page_num - 1) * 50
	rows := tool.Query_DB(db, "select distinct title, data, type from acl where data != '' and title not like 'user:%' order by title desc limit ?, 50", offset)
	body := strings.Builder{}
	count := 0
	for rows.Next() {
		title, value, field := "", "", ""
		if rows.Scan(&title, &value, &field) != nil {
			continue
		}
		last_change := ""
		tool.QueryRow_DB(db, "select time from re_admin where what like ? order by time desc limit 1", []any{&last_change}, "acl ("+title+")%")
		why := ""
		tool.QueryRow_DB(db, "select data from acl where title = ? and type = 'why' limit 1", []any{&why}, title)
		body.WriteString(`<li>` + tool.HTML_escape(last_change) + ` | <a href="/acl/` + tool.Url_parser(title) + `">` + tool.HTML_escape(title) + `</a> | ` + tool.HTML_escape(value) + ` (` + tool.HTML_escape(field) + `)`)
		if why != "" {
			body.WriteString(` | ` + tool.HTML_escape(why))
		}
		body.WriteString(`</li>`)
		count++
	}
	rows.Close()
	body.WriteString(tool.Get_page_control(db, page_num, count, 50, "/list/document/acl/{}"))
	return list_extra_page(db, config, tool.Get_language(db, "acl_document_list", true), `<ul>`+body.String()+`</ul>`)
}
