package route

import (
	"opennamu/route/tool"
	"strconv"
)

func Api_bbs_w_page_view_post(config tool.Config, set_id string, set_code string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    return_data := make(map[string]any)
    return_data["response"] = "ok"

    page_view_str := ""
    exist := tool.QueryRow_DB(
        db,
        "select set_data from bbs_data where set_name = 'view_count' and set_id = ? and set_code = ?",
        []any{ &page_view_str },
        set_id,
        set_code,
    )

    page_view_int := tool.Str_to_int(page_view_str) + 1
    page_view_str = strconv.Itoa(page_view_int)

    if exist {
        tool.Exec_DB(
            db,
            "update bbs_data set set_data = ? where set_name = 'view_count' and set_id = ? and set_code = ?",
            page_view_str,
            set_id,
            set_code,
        )
    } else {
        tool.Exec_DB(
            db,
            "insert into bbs_data (set_name, set_id, set_code, set_data) values ('view_count', ?, ?, ?)",
            set_id,
            set_code,
            page_view_str,
        )
    }

    return return_data
}