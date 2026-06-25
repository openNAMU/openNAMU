package route

import (
	"opennamu/route/tool"
)

func Api_list_recent_change(config tool.Config, set_type string, limit string, num string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    if set_type == "edit" {
        set_type = ""
    }

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

    rows := tool.Query_DB(
        db,
        "select id, title from rc where type = ? order by date desc limit ?, ?",
        set_type,
        page_int,
        limit_int,
    )
    defer rows.Close()

    data_list := [][]string{}

    admin_auth := tool.Check_acl(db, "", "", "hidel_auth", config.IP)
    ip_parser_temp := map[string][]string{}

    for rows.Next() {
        var id string
        var title string

        err := rows.Scan(&id, &title)
        if err != nil {
            panic(err)
        }

        date := ""
        ip := ""
        send := ""
        leng := ""
        hide := ""
        type_data := ""
        tool.QueryRow_DB(
            db,
            "select date, ip, send, leng, hide, type from history where id = ? and title = ?",
            []any{ &date, &ip, &send, &leng, &hide, &type_data },
            id, title,
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

        if hide == "" || admin_auth {
            data_list = append(data_list, []string{
                id,
                title,
                date,
                ip_pre,
                send,
                leng,
                hide,
                ip_render,
                type_data,
            })
        } else {
            data_list = append(data_list, []string{"", "", "", "", "", "", hide, "", ""})
        }
    }

    return_data := make(map[string]any)
    return_data["response"] = "ok"
    return_data["data"] = data_list

    return return_data
}
