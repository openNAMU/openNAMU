package route

import "opennamu/route/tool"

func View_vote_list(config tool.Config, type_str string, num_str string) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    if type_str == "" {
        type_str = "open"
    }

    title := ""
    switch type_str {
    case "open":
        title = tool.Get_language(db, "open_vote_list", true)
    default:
        title = tool.Get_language(db, "close_vote_list", true)
    }

    api_data := Api_vote_list(config, type_str, num_str)

    data_html := ""
    for _, in_data := range api_data["data"].([][]string) {
        data_html += "<li>"
        data_html += "<a href=\"/vote/" + tool.Url_parser(in_data[1]) + "\">"
        data_html += tool.HTML_escape(in_data[0])
        data_html += "</a>"
        data_html += "</li>"
    }

    return tool.Get_template(
        db,
        config,
        title,
        data_html,
        []any{},
        [][]any{},
        map[string]string{},
    )
}