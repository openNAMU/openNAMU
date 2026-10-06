package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func View_related_documents_html(db *sql.DB, config tool.Config, doc_name string, names []string) string {
	if len(names) == 0 {
		return "" // §8 empty list renders nothing (View_w also skips the call on an empty result, line 145)
	}

	title := tool.Get_language(db, "related_document", true)

	html := "<h2>" + title + "</h2><hr class=\"main_hr\"><ul>"
	for _, name := range names {
		html += `<li><a href="/w/` + tool.Url_parser(name) + `">` + tool.HTML_escape(name) + `</a></li>`
	}
	html += "</ul>"

	return html
}
