package route

import "opennamu/route/tool"

func Api_w_page_view(config tool.Config, doc_name string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    return_data := make(map[string]any)
    return_data["response"] = "ok"

    page_view := 0

    pv_continue := tool.Get_setting(db, "not_use_view_count", "")
    if len(pv_continue) == 0 || pv_continue[0][0] == "" {
        page_view_str := ""

        tool.QueryRow_DB(
            db,
            "select set_data from data_set where doc_name = ? and set_name = 'view_count' and doc_rev = ''",
            []any{ &page_view_str },
            doc_name,
        )

        page_view = tool.Str_to_int(page_view_str)        
    }

    return_data["data"] = page_view

    return return_data
}