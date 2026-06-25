package route

import (
	"opennamu/route/tool"
)

func Api_w_xref(config tool.Config, num_str string, doc_name string, do_type string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    page := tool.Str_to_int(num_str)
    num := 0
    if page * 50 > 0 {
        num = page * 50 - 50
    }

    link_case_insensitive := ""
    tool.QueryRow_DB(
        db,
        "select data from other where name = 'link_case_insensitive'",
        []any{ &link_case_insensitive },
        doc_name,
    )

    if link_case_insensitive != "" {
        link_case_insensitive = " collate nocase"
    }

    query := ""
    if do_type == "1" {
        query = "select distinct link, type from back where title" + link_case_insensitive + " = ? and not type = 'no' and not type = 'nothing' order by type asc, link asc limit ?, 50"
    } else {
        query = "select distinct title, type from back where link" + link_case_insensitive + " = ? and not type = 'no' and not type = 'nothing' order by type asc, title asc limit ?, 50"
    }

    rows := tool.Query_DB(
        db,
        query,
        doc_name,
        num,
    )
    defer rows.Close()

    data_list := [][]string{}

    for rows.Next() {
        var name string
        var type_data string

        err := rows.Scan(&name, &type_data)
        if err != nil {
            panic(err)
        }

        data_list = append(data_list, []string{name, type_data})
    }

    return_data := make(map[string]any)
    return_data["response"] = "ok"
    return_data["data"] = data_list

    return return_data
}
