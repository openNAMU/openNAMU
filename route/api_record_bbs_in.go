package route

import "opennamu/route/tool"

func Api_record_bbs_in(config tool.Config, user_name string, set_id string, page string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    page_int := tool.Str_to_int(page)
    num := 0
    if page_int * 50 > 0 {
        num = page_int * 50 - 50
    }

    rows := tool.Query_DB(
        db,
        `select set_code from bbs_data where set_name = "user_id" and set_id = ? and set_data = ? order by set_code desc limit ?, 50`,
        set_id,
        user_name,
        num,
    )
    defer rows.Close()

    data_list := []string{}

    for rows.Next() {
        var set_code string

        err := rows.Scan(&set_code)
        if err != nil {
            panic(err)
        }

        data_list = append(data_list, set_code)
    }

    result_data := make(map[string]any)
    result_data["response"] = "ok"
    result_data["data"] = data_list

    return result_data
}