package route

import (
	"opennamu/route/tool"
	"strings"
)

func View_main_search(config tool.Config, keyword string, num string, search_type string) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    if keyword == "" {
        return tool.Get_redirect("/")
    }

    data_html := `
        <form method="post">
            <input class="opennamu_width_200 __ON_INPUT__" name="search" value="` + tool.HTML_escape(keyword) + `">
            <button class="__ON_BUTTON__" type="submit">` + tool.Get_language(db, "search", true) + `</button>
        </form>
        <hr class="main_hr">
        <a href="/search_page/1/` + tool.Url_parser(keyword) + `">(` + tool.Get_language(db, "search_document_name", true) + `)</a>
        <a href="/search_data_page/1/` + tool.Url_parser(keyword) + `">(` + tool.Get_language(db, "search_document_data", true) + `)</a>
    `

    name_new := ""
    if strings.HasPrefix(keyword, "분류:") {
        name_new = strings.ReplaceAll(keyword, "분류:", "category:")
    } else if strings.HasPrefix(keyword, "사용자:") {
        name_new = strings.ReplaceAll(keyword, "사용자:", "user:")
    } else if strings.HasPrefix(keyword, "파일:") {
        name_new = strings.ReplaceAll(keyword, "파일:", "file:")
    }

    if name_new != "" {
        data_html += ` <a href="/search_page/1/` + tool.Url_parser(name_new) + `">(` + tool.HTML_escape(name_new) + `)</a>`
    }

    data_api_exist := Api_w_raw(config, keyword, "true", "")
    data_api_exist_in := data_api_exist["data"].(string)
    
    main_document_name := keyword
    link_id := `class="opennamu_not_exist_link"`
    if data_api_exist_in != "" {
        link_id = ""
        main_document_name = data_api_exist_in
    }

    data_html += `
        <ul>
            <li>
                ` + tool.Get_language(db, "go", true) + ` : <a ` + link_id + ` href="/w/` + tool.Url_parser(main_document_name) + `">` + tool.HTML_escape(main_document_name) + `</a>
            </li>
        </ul>
    `

    data_api := Api_func_search(config, keyword, num, search_type)
    data_api_in := data_api["data"].([]string)

    data_html += "<ul>"
    for _, v := range data_api_in {
        data_html += `<li><a href="/w/` + tool.Url_parser(v) + `">` + tool.HTML_escape(v) + `</a></li>`
    }

    data_html += "</ul>"

    num_int := tool.Str_to_int(num)

    if search_type == "title" {
        data_html += tool.Get_page_control(
            db,
            num_int,
            len(data_api_in),
            50,
            "/search_page/{}/" + tool.Url_parser(keyword),
        )
    } else {
        data_html += tool.Get_page_control(
            db,
            num_int,
            len(data_api_in),
            50,
            "/search_page/{}/" + tool.Url_parser(keyword),
        )
    }

    out := tool.Get_template(
        db,
        config,
        keyword,
        data_html,
        []any{ "(" + tool.Get_language(db, "search", true) + ")" },
        [][]any{
            { "other", tool.Get_language(db, "return", true) },
        },
        map[string]string{},
    )

    return out
}
