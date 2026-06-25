package route

import "opennamu/route/tool"

func Api_list_long_page(config tool.Config, num string, set_type string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    page_int := tool.Str_to_int(num)
    if page_int > 0 {
        page_int = (page_int * 50) - 50
    } else {
        page_int = 0
    }

    query := ""
    if set_type == "long" {
        query = "select doc_name, set_data from data_set where set_name = 'length' and doc_rev = '' and " + tool.Get_except_document_name_SQL("doc_name") + " order by set_data + 0 desc limit ?, 50"
    } else {
        query = "select doc_name, set_data from data_set where set_name = 'length' and doc_rev = '' and " + tool.Get_except_document_name_SQL("doc_name") + " order by set_data + 0 asc limit ?, 50"
    }

    rows := tool.Query_DB(
        db,
        query,
        page_int,
    )
    defer rows.Close()

    data_list := [][]string{}

    for rows.Next() {
        var doc_name string
        var length string

        err := rows.Scan(&doc_name, &length)
        if err != nil {
            panic(err)
        }

        data_list = append(data_list, []string{ doc_name, length })
    }

    return_data := make(map[string]any)
    return_data["data"] = data_list

    return return_data
}