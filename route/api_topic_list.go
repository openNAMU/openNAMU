package route

import (
	"database/sql"
	"opennamu/route/tool"
)

func Api_topic_list(config tool.Config, num string, doc_name string, do_type string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    page_int := tool.Str_to_int(num)
    if page_int > 0 {
        page_int = (page_int * 50) - 50
    } else {
        page_int = 0
    }

    var rows *sql.Rows

    switch do_type {
    case "close":
        rows = tool.Query_DB(
            db,
            "select code, sub, stop, agree, date from rd where title = ? and stop = 'O' order by sub asc limit ?, 50",
            doc_name,
            page_int,
        )
    case "open":
        rows = tool.Query_DB(
            db,
            "select code, sub, stop, agree, date from rd where title = ? and agree = 'O' order by sub asc limit ?, 50",
            doc_name,
            page_int,
        )
    default:
        rows = tool.Query_DB(
            db,
            "select code, sub, stop, agree, date from rd where title = ? order by sub asc limit ?, 50",
            doc_name,
            page_int,
        )
    }
    defer rows.Close()

    data_list := [][]string{}
    ip_parser_temp := map[string][]string{}

    for rows.Next() {
        var code string
        var sub string
        var stop string
        var agree string
        var date string

        err := rows.Scan(&code, &sub, &stop, &agree, &date)
        if err != nil {
            panic(err)
        }

        ip := ""
        id := ""
        tool.QueryRow_DB(
            db,
            "select ip, id from topic where code = ? order by id + 0 desc limit 1",
            []any{ &ip, &id },
            code,
        )

        var ip_pre string
        var ip_render string

        if _, ok := ip_parser_temp[ip]; ok {
            ip_pre = ip_parser_temp[ip][0]
            ip_render = ip_parser_temp[ip][1]
        } else {
            ip_pre = tool.IP_preprocess(db, ip, config.IP)[0]
            ip_render = tool.IP_parser(db, ip, config.IP)

            ip_parser_temp[ip] = []string{ip_pre, ip_render}
        }

        data_list = append(data_list, []string{
            code,
            sub,
            stop,
            agree,
            ip_pre,
            ip_render,
            date,
            id,
        })
    }

    return_data := make(map[string]any)
    return_data["data"] = data_list

    return return_data
}
