package route

import (
	"opennamu/route/tool"
)

func Api_list_recent_discuss(config tool.Config, limit string, num string, set_type string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    limit_int := tool.Str_to_int(limit)
    if limit_int > 50 || limit_int < 0 {
        limit_int = 50
    }

    page_int := tool.Str_to_int(num)
    if page_int > 0 {
        page_int = (page_int * limit_int) - limit_int
    } else {
        page_int = 0
    }

    query := ""
    switch set_type {
    case "normal":
        query = "select title, sub, date, code, stop, agree from rd order by date desc limit ?, ?"
    case "close":
        query = "select title, sub, date, code, stop, agree from rd where stop = 'O' order by date desc limit ?, ?"
    default:
        query = "select title, sub, date, code, stop, agree from rd where stop != 'O' order by date desc limit ?, ?"
    }

    rows := tool.Query_DB(
        db,
        query,
        page_int, limit_int,
    )
    defer rows.Close()

    data_list := [][]string{}
    ip_parser_temp := map[string][]string{}

    for rows.Next() {
        var title string
        var sub string
        var date string
        var code string
        var stop string
        var agree string

        err := rows.Scan(&title, &sub, &date, &code, &stop, &agree)
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
            title,
            sub,
            date,
            code,
            stop,
            ip_pre,
            ip_render,
            id,
            agree,
        })
    }

    return_data := make(map[string]any)
    return_data["response"] = "ok"
    return_data["data"] = data_list

    return return_data
}
