package route

import "opennamu/route/tool"

func View_topic_list(config tool.Config, doc_name string, do_type string, num string) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    data_api := Api_topic_list(config, num, doc_name, do_type)
    data_api_in := data_api["data"].([][]string)

    data_html := `
        <a href="/topic_close/1/` + tool.Url_parser(doc_name) + `">(` + tool.Get_language(db, "closed_discussion", true) + `)</a>
        <a href="/topic_agree/1/` + tool.Url_parser(doc_name) + `">(` + tool.Get_language(db, "agreed_discussion", true) + `)</a>
        <hr class="main_hr">
        <a href="/thread/0/` + tool.Url_parser(doc_name) + `">(` + tool.Get_language(db, "make_new_topic", true) + `)</a>
    `
    sub_title := ""
    menu := [][]any{
        { "topic/" + tool.Url_parser(doc_name), tool.Get_language(db, "return", true) },
    }
    url := ""

    switch do_type {
    case "close":
        sub_title = tool.Get_language(db, "closed_discussion", true)
        url = "/topic_close/{}/" + tool.Url_parser(doc_name)
    case "agree":
        sub_title = tool.Get_language(db, "agreed_discussion", true)
        url = "/topic_agree/{}/" + tool.Url_parser(doc_name)
    default:
        sub_title = tool.Get_language(db, "discussion_list", true)
        menu = [][]any{
            { "w/" + tool.Url_parser(doc_name), tool.Get_language(db, "return", true) },
        }
        url = "/topic_page/{}/" + tool.Url_parser(doc_name)
    }

    for _, in_data := range data_api_in {
        data_html += `<h2><a href="/thread/` + in_data[0] + `">` + in_data[0] + `. ` + tool.HTML_escape(in_data[1]) + `</a></h2>`
    }

    num_int := tool.Str_to_int(num)
    data_html += tool.Get_page_control(db, num_int, len(data_api_in), 50, url)

    out := tool.Get_template(
        db,
        config,
        doc_name,
        data_html,
        []any{ "(" + sub_title + ")" },
        menu,
        map[string]string{},
    )

    return out
}
