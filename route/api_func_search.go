package route

import (
	"opennamu/route/tool"
)

func Api_func_search(config tool.Config, keyword string, num_str string, search_type string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    page := tool.Str_to_int(num_str)
    num := 0
    if page * 50 > 0 {
        num = page * 50 - 50
    }

    name := keyword
    query := ""
    
    if search_type == "title" {
        name = tool.Do_remove_spaces(name)
        query = "select title from data where replace(title, ' ', '') collate nocase like ? order by title limit ?, 50"
    } else {
        query = "select title from data where data collate nocase like ? order by title limit ?, 50"
    }

    title_list := []string{}

    rows := tool.Query_DB(
        db,
        query,
        "%" + name + "%", num,
    )
    defer rows.Close()

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
