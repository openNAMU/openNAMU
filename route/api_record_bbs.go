package route

import "opennamu/route/tool"

func Api_record_bbs(config tool.Config, user_name string, page string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    page_int := tool.Str_to_int(page)
    num := 0
    if page_int * 50 > 0 {
        num = page_int * 50 - 50
    }

    rows := tool.Query_DB(
        db,
        `select distinct set_id from bbs_data where set_name = "user_id" and set_data = ? order by set_id desc limit ?, 50`,
        user_name,
        num,
    )
    defer rows.Close()

    data_list := [][]string{}

    for rows.Next() {
        var set_id string

        err := rows.Scan(&set_id)
        if err != nil {
            panic(err)
        }

        date := ""
        tool.QueryRow_DB(
            db,
            `select set_data from bbs_data where set_name = "date" and set_id = ? order by set_data desc limit 1`,
            []any{ &date },
            set_id,
        )

        data_list = append(data_list, []string{set_id, date})
    }

    result_data := make(map[string]any)
    result_data["response"] = "ok"
    result_data["data"] = data_list

    return result_data
}