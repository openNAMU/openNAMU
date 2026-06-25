package route

import (
	"opennamu/route/tool"
)

func View_main_search_post(config tool.Config, search_type string, goto_document bool, keyword string) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    if goto_document {
        data_api_exist := Api_w_raw(config, keyword, "true", "")
        data_api_exist_in := data_api_exist["data"].(string)

        data_api := Api_func_search(config, keyword, "1", "title")
        data_api_in := data_api["data"].([]string)

        if data_api_exist_in != "" {
            return tool.Get_redirect("/w/" + tool.Url_parser(data_api_exist_in))
        } else if len(data_api_in) > 0 {
            return tool.Get_redirect("/w/" + tool.Url_parser(data_api_in[0]))
        }
    }

    switch search_type {
    case "":
        return tool.Get_redirect("/search/" + tool.Url_parser(keyword))
    case "title":
        return tool.Get_redirect("/search_page/1/" + tool.Url_parser(keyword))
    default:
        return tool.Get_redirect("/search_data_page/1/" + tool.Url_parser(keyword))
    }
}
