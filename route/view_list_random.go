package route

import (
	"opennamu/route/tool"
)

func View_list_random(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    data_list := Api_list_random(config, 50)

    data_html := "<ul>"
    for _, title := range data_list["data"].([]string) {
        data_html += "<li>" 
        data_html += "<a href=\"/w/" + tool.Url_parser(title) + "\">" + tool.HTML_escape(title) + "</a>"

        if tool.Get_document_is_redirect(db, title) {
            data_html += " (" + tool.Get_language(db, "redirect", false) + ")"
        }

        data_html += "</li>"
    }
    
    data_html += "</ul>"

    out := tool.Get_template(
        db,
        config,
        tool.Get_language(db, "random_list", true),
        data_html,
        []any{},
        [][]any{
            { "other", tool.Get_language(db, "return", true) },
        },
        map[string]string{},
    )

    return out
}
