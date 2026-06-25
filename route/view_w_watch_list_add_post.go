package route

import "opennamu/route/tool"

func View_w_watch_list_add_post(config tool.Config, doc_name string, do_type string) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    name_from := false
    switch do_type {
    case "watchlist_from":
        name_from = true
        do_type = "watchlist"
    case "star_doc_from":
        name_from = true
        do_type = "star_doc"
    }

    if do_type != "watchlist" {
        do_type = "star_doc"
    }

    api_data := Api_w_watch_list_post(config, doc_name, do_type)

    out := ""
    if api_data["response"] != "ok" {
        out = tool.Get_error_page(db, config, "error")
    } else {
        if name_from {
            out = tool.Get_redirect("/w/" + doc_name)
        } else if do_type == "watchlist" {
            out = tool.Get_redirect("/watch_list")
        } else {
            out = tool.Get_redirect("/star_doc")
        }
    }

    return out
}