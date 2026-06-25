package route

import (
	"database/sql"
	"opennamu/route/tool"
	"sort"
)

func bbs_list(db *sql.DB) map[string]string {
    rows := tool.Query_DB(
        db,
        "select set_data, set_id from bbs_set where set_name = 'bbs_name' and " + tool.Get_except_set_id_SQL(),
    )
    defer rows.Close()

    data_list := map[string]string{}

    for rows.Next() {
        var name string
        var id string

        err := rows.Scan(&name, &id)
        if err != nil {
            panic(err)
        }

        data_list[name] = id
    }

    return data_list
}

type BBS_item struct {
    Id string
    Name string
    Type string
    Date string
}

func Api_bbs_list(config tool.Config) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    data_list := bbs_list(db)
    items := make([]BBS_item, 0, len(data_list))

    for k, v := range data_list {
        bbs_type := ""
        tool.QueryRow_DB(
            db,
            "select set_data from bbs_set where set_name = 'bbs_type' and set_id = ?",
            []any{ &bbs_type },
            v,
        )

        bbs_date := ""
        tool.QueryRow_DB(
            db,
            "select set_data from bbs_data where set_id = ? and set_name = 'date' order by set_code + 0 desc limit 1",
            []any{ &bbs_date },
            v,
        )

        items = append(items, BBS_item{
            Id : v,
            Name : k,
            Type : bbs_type,
            Date : bbs_date,
        })
    }

    sort.Slice(items, func(i, j int) bool {
        return items[i].Date > items[j].Date
    })

    data_list_sub := make([][]string, 0, len(items))
    for _, item := range items {
        data_list_sub = append(data_list_sub, []string{ item.Name, item.Id, item.Type, item.Date })
    }

    return_data := make(map[string]any)
    return_data["response"] = "ok"
    return_data["data"] = data_list_sub

    return return_data
}
