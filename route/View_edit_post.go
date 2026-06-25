package route

import (
	"opennamu/route/tool"
)

func View_edit_post(config tool.Config, doc_name string, data string, send string, agree string) string {   
    db := tool.DB_connect()
    defer tool.DB_close(db)

    return_data := Api_edit_post(config, doc_name, data, send, agree)

    result_html := ""
    if return_data["response"].(string) == "ok" {
        result_html = tool.Get_redirect("/w/" + tool.Url_parser(doc_name))
    } else {
        result_html = tool.Get_error_page(db, config, return_data["data"].(string))
    }

    return result_html
}