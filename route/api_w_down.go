package route

import "opennamu/route/tool"

func Api_w_down(config tool.Config, doc_name string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    rows := tool.Query_DB(
        db,
        "select title from data where title like ?",
        doc_name,
    )

    title_list := []string{}

    for rows.Next() {
        var title string

        err := rows.Scan(&title)
        if err != nil {
            panic(err)
        }

        title_list = append(title_list, title)
    }

    return_data := make(map[string]any)
    return_data["response"] = "ok"
    return_data["data"] = title_list

    return return_data
}