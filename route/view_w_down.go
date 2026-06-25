package route

import "opennamu/route/tool"

func View_w_down(config tool.Config, doc_name string) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    data_api := Api_w_down(config, doc_name)
    data_api_in := data_api["data"].([]string)

    data_html := "<ul>"

    for _, v := range data_api_in {
        data_html += `<li><a href="/w/` + tool.Url_parser(v) + `">` + v + `</a></li>`
    }

    data_html += "</ul>"

    out := tool.Get_template(
        db,
        config,
        doc_name,
        data_html,
        []any{ "(" + tool.Get_language(db, "sub", true) + ")" },
        [][]any{
            { "w/" + tool.Url_parser(doc_name), tool.Get_language(db, "return", true) },
        },
        map[string]string{},
    )

    return out
}
