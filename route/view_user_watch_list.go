package route

import (
	"opennamu/route/tool"
)

func View_user_watch_list(config tool.Config, num string, do_type string) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    api_data := Api_user_watch_list(config, config.IP, num, do_type)
    data_html := ""

    if api_data["response"] != "ok" {
        return tool.Get_error_page(db, config, "auth")
    } else {
        data_html += "<ul>"
        for _, title := range api_data["data"].([]string) {
            data_html += "<li><a href=\"/w/" + tool.Url_parser(title) + "\">" + tool.HTML_escape(title) + "</a></li>"
        }
        data_html += "</ul>"
    }

    title := tool.Get_language(db, "watchlist", true)
    if do_type == "star_doc" {
        title = tool.Get_language(db, "star_doc", true)
    }

    out := tool.Get_template(
        db,
        config,
        title,
        data_html,
        []any{},
        [][]any{
            { "user/" + tool.Url_parser(config.IP), tool.Get_language(db, "return", false) },
            { "watch_list", tool.Get_language(db, "watchlist", false) },
            { "star_doc", tool.Get_language(db, "star_doc", false) },
        },
        map[string]string{},
    )

    return out
}
