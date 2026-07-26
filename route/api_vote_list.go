package route

import (
	"database/sql"
	"opennamu/route/tool"
)

func Api_vote_list(config tool.Config, set_type string, num_str string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

    page := tool.Str_to_int(num_str)
    num := 0
    if page * 50 > 0 {
        num = page * 50 - 50
    }

    if set_type == "" {
        set_type = "open"
    }

    var rows *sql.Rows
    if set_type == "open" {
        rows = tool.Query_DB(
            db,
            `select name, id, type from vote where type = "open" or type = "n_open" limit ?, 50`,
            num,
        )
    } else {
        rows = tool.Query_DB(
            db,
            `select name, id, type from vote where type = "close" or type = "n_close" limit ?, 50`,
            num,
        )
    }
    defer rows.Close()

    data_list := [][]string{}
    for rows.Next() {
        var name string
        var id string
        var type_str string

        err := rows.Scan(&name, &id, &type_str)
        if err != nil {
            panic(err)
        }

        data_list = append(data_list, []string{
            name,
            id,
            type_str,
        })
    }

	return_data := make(map[string]any)
	return_data["response"] = "ok"
    return_data["data"] = data_list

	return return_data
}
