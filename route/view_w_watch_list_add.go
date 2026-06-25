package route

import "opennamu/route/tool"

func View_w_watch_list_add(config tool.Config, doc_name string, do_type string) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    switch do_type {
    case "watchlist_from":
        do_type = "watchlist"
    case "star_doc_from":
        do_type = "star_doc"
    }

    if do_type != "watchlist" {
        do_type = "star_doc"
    }

    out := tool.Get_template(
        db,
        config,
        doc_name,
        `<form method="post">
            <button id="opennamu_save_button" type="submit">` + tool.Get_language(db, "send", true) + `</button>
        </form>`,
        []any{ "(" + tool.Get_language(db, do_type, true) + ")" },
        [][]any{
            { "w/" + tool.Url_parser(doc_name), tool.Get_language(db, "return", true) },
        },
        map[string]string{},
    )

    return out
}
