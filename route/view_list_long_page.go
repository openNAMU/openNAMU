package route

import "opennamu/route/tool"

func View_list_long_page(config tool.Config, num string, set_type string) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

	api_data := Api_list_long_page(config, num, set_type)
    api_data_list := api_data["data"].([][]string)

    if set_type != "long" {
        set_type = "short"
    }

    title := tool.Get_language(db, "long_page", true)
    if set_type == "short" {
        title = tool.Get_language(db, "short_page", true)
    }

    data_html := ""

    for _, data := range api_data_list {
        doc_name := tool.Url_parser(data[0])
        doc_title := tool.HTML_escape(data[0])
        length := tool.HTML_escape(data[1])

        right := `<a href="/w/` + doc_name + `">` + doc_title + `</a>`
        left := length

        data_html += tool.Get_list_ui(right, left, "", "")
    }

    data_html += tool.Get_page_control(
        db,
        tool.Str_to_int(num),
        len(api_data_list),
        50,
        "/list/document/" + set_type + "/{}",
    )

    out := tool.Get_template(
        db,
        config,
        title,
        data_html,
        []any{},
        [][]any{
            { "other", tool.Get_language(db, "return", true) },
        },
        map[string]string{},
    )

    return out
}