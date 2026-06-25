package route

import "opennamu/route/tool"

func View_record_bbs_comment(config tool.Config, user_name string, page string) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    api_data := Api_record_bbs_comment(config, user_name, page)
    if api_data["response"].(string) != "ok" {
        return tool.Get_error_page(
            db,
            config,
            "auth",
        )
    }

    data_list := api_data["data"].([][]string)
    data_html := ""

    for _, data := range data_list {
        bbs_name := Api_bbs_num_to_name(db, data[0])
        set_id := data[0]
        date := data[1]

        link := `<a href="/record_bbs_comment_in/` + tool.Url_parser(set_id) + `/` + tool.Url_parser(user_name) + `">` + tool.HTML_escape(bbs_name) + `</a>`

        data_html += tool.Get_list_ui(link, date, "", "")
    }

    out := tool.Get_template(
        db,
        config,
        user_name,
        data_html,
        []any{ "(" + tool.Get_language(db, "bbs_comment_record", true) + ")" },
        [][]any{},
        map[string]string{},
    )

    return out
}