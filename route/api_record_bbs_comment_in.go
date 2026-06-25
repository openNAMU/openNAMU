package route

import (
	"opennamu/route/tool"
	"strings"
)

func Api_record_bbs_comment_in(config tool.Config, user_name string, bbs_id string, page string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    page_int := tool.Str_to_int(page)
    num := 0
    if page_int * 50 > 0 {
        num = page_int * 50 - 50
    }

    rows := tool.Query_DB(
        db,
        `select set_id, set_code from bbs_data where set_name = "comment_user_id" and set_id like ? and set_data = ? order by set_code desc limit ?, 50`,
        bbs_id + "-%",
        user_name,
        num,
    )
    defer rows.Close()

    data_list := [][]string{}

    for rows.Next() {
        var set_id string
        var set_code string

        err := rows.Scan(&set_id, &set_code)
        if err != nil {
            panic(err)
        }

        post_id := strings.Split(set_id, "-")[1]

        data_list = append(data_list, []string{
            bbs_id,
            post_id,
            set_id,
            set_code,
        })
    }

    result_data := make(map[string]any)
    result_data["response"] = "ok"
    result_data["data"] = data_list

    return result_data
}