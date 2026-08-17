package route

import (
	"strings"

	"opennamu/route/tool"
)

func View_list_file(config tool.Config, image_only bool) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	body := strings.Builder{}
	rows := tool.Query_DB(db, "select title, type from data where title like 'file:%' order by title")
	for rows.Next() {
		name, kind := "", ""
		if rows.Scan(&name, &kind) != nil {
			continue
		}
		if image_only && !strings.HasPrefix(kind, "image") && !strings.Contains(name, ".png") && !strings.Contains(name, ".jpg") && !strings.Contains(name, ".jpeg") && !strings.Contains(name, ".gif") && !strings.Contains(name, ".webp") {
			continue
		}
		body.WriteString(`<li><a href="/w/` + tool.Url_parser(name) + `">` + tool.HTML_escape(name) + `</a></li>`)
	}
	rows.Close()
	return list_extra_page(db, config, tool.Get_language(db, func() string {
		if image_only {
			return "image_list"
		}
		return "file_list"
	}(), true), `<ul>`+body.String()+`</ul>`)
}
