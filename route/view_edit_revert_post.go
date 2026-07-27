package route

import "opennamu/route/tool"

func View_edit_revert_post(config tool.Config, doc_name string, rev string, send string, agree string) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    api_data := Api_edit_revert_post(config, doc_name, rev, send, agree)
    response := api_data["response"].(string)
    switch response {
    case "ok", "not exist":
        return tool.Get_redirect("/w/" + tool.Url_parser(doc_name))
    case "require auth":
        return tool.Get_error_page(db, config, "auth")
    }

    error_name, ok := api_data["data"].(string)
    if !ok {
        error_name = "error"
    }

    return tool.Get_error_page(db, config, error_name)
}
