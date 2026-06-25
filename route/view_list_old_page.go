package route

import "opennamu/route/tool"

func View_list_old_page(config tool.Config, num string, set_type string) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

	api_data := Api_list_old_page(config, num, set_type)
    api_data_list := api_data["data"].([][]string)

    if set_type != "old" {
        set_type = "new"
    }

    title := tool.Get_language(db, "old_page", true)
    if set_type == "new" {
        title = tool.Get_language(db, "new_page", true)
    }

    data_html := ""

    for _, data := range api_data_list {
        doc_name := tool.Url_parser(data[0])
        doc_title := tool.HTML_escape(data[0])
        date := tool.HTML_escape(data[1])

        left := `<a href="/w/` + doc_name + `">` + doc_title + `</a>`
        right := date

        data_html += tool.Get_list_ui(left, right, "", "")
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