package route

import "opennamu/route/tool"

func Api_bbs_w_page_view(config tool.Config, set_id string, set_code string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    return_data := make(map[string]any)
    return_data["response"] = "ok"

    page_view := 0
    page_view_str := ""

    tool.QueryRow_DB(
        db,
        "select set_data from bbs_data where set_name = 'view_count' and set_id = ? and set_code = ?",
        []any{ &page_view_str },
        set_id,
        set_code,
    )

    page_view = tool.Str_to_int(page_view_str)

    return_data["data"] = page_view

    return return_data
}