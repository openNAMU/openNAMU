package route

import (
	"opennamu/route/tool"
)

func Api_list_title_index(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    page_int := tool.Str_to_int(other_set["num"])
    if page_int > 0 {
        page_int = (page_int * 50) - 50
    } else {
        page_int = 0
    }

    rows := tool.Query_DB(
        db,
        "select title from data limit ?, 50",
        page_int,
    )
    defer rows.Close()

    data_list := []string{}

    for rows.Next() {
        var title string

        err := rows.Scan(&title)
        if err != nil {
            panic(err)
        }

        data_list = append(data_list, title)
    }

    return_data := make(map[string]any)
    return_data["data"] = data_list

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}
