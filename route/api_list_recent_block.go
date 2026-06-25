package route

import (
	"database/sql"
	"opennamu/route/tool"
)

func Api_list_recent_block(config tool.Config, num string, set_type string, why string, user_name string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    page_int := tool.Str_to_int(num)
    if page_int > 0 {
        page_int = (page_int * 50) - 50
    } else {
        page_int = 0
    }

    // private 공개 안되도록 조심할 것
    var rows *sql.Rows
    
    switch set_type {
    case "all":
        query := ""
        if why != "" {
            query = "select why, block, blocker, end, today, band, ongoing from rb where band != 'private' and why like ? order by today desc limit ?, 50"
        } else {
            query = "select why, block, blocker, end, today, band, ongoing from rb where band != 'private' order by today desc limit ?, 50"
        }

        if why != "" {
            rows = tool.Query_DB(
                db,
                query,
                why + "%", page_int,
            )
        } else {
            rows = tool.Query_DB(
                db,
                query,
                page_int,
            )
        }
    case "ongoing":
        rows = tool.Query_DB(
            db,
            "select why, block, blocker, end, today, band, ongoing from rb where ongoing = '1' and band != 'private' order by end desc limit ?, 50",
            page_int,
        )
    case "regex":
        rows = tool.Query_DB(
            db,
            "select why, block, blocker, end, today, band, ongoing from rb where band = 'regex' order by today desc limit ?, 50",
            page_int,
        )
    case "private":
        rows = tool.Query_DB(
            db,
            "select why, block, blocker, end, today, band, ongoing from rb where band = 'private' order by today desc limit ?, 50",
            page_int,
        )
    case "user":
        rows = tool.Query_DB(
            db,
            "select why, block, blocker, end, today, band, ongoing from rb where block = ? and band != 'private' order by today desc limit ?, 50",
            user_name, page_int,
        )
    case "cidr":
        rows = tool.Query_DB(
            db,
            "select why, block, blocker, end, today, band, ongoing from rb where band = 'cidr' order by today desc limit ?, 50",
            page_int,
        )
    default:
        rows = tool.Query_DB(
            db,
            "select why, block, blocker, end, today, band, ongoing from rb where blocker = ? and band != 'private' order by today desc limit ?, 50",
            user_name, page_int,
        )
    }
    defer rows.Close()

    data_list := [][]string{}
    ip_parser_temp := map[string][]string{}

    for rows.Next() {
        var why string
        var block string
        var blocker string
        var end string
        var today string
        var band string
        var ongoing string

        err := rows.Scan(&why, &block, &blocker, &end, &today, &band, &ongoing)
        if err != nil {
            panic(err)
        }

        var ip_pre_blocker string
        var ip_render_blocker string

        if _, ok := ip_parser_temp[blocker]; ok {
            ip_pre_blocker = ip_parser_temp[blocker][0]
            ip_render_blocker = ip_parser_temp[blocker][1]
        } else {
            ip_pre_blocker = tool.IP_preprocess(db, blocker, config.IP)[0]
            ip_render_blocker = tool.IP_parser(db, blocker, config.IP)

            ip_parser_temp[blocker] = []string{ip_pre_blocker, ip_render_blocker}
        }

        var ip_pre_block string
        var ip_render_block string

        if band == "" {
            if _, ok := ip_parser_temp[block]; ok {
                ip_pre_block = ip_parser_temp[block][0]
                ip_render_block = ip_parser_temp[block][1]
            } else {
                ip_pre_block = tool.IP_preprocess(db, block, config.IP)[0]
                ip_render_block = tool.IP_parser(db, block, config.IP)

                ip_parser_temp[block] = []string{ip_pre_block, ip_render_block}
            }
        } else {
            ip_pre_block = block
            ip_render_block = block
        }

        data_list = append(data_list, []string{
            why,
            ip_pre_block,
            ip_render_block,
            ip_pre_blocker,
            ip_render_blocker,
            end,
            today,
            band,
            ongoing,
        })
    }

    if set_type == "private" {
        if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
            data_list = [][]string{}
        }
    }

    return_data := make(map[string]any)
    return_data["response"] = "ok"
    return_data["data"] = data_list

    return return_data
}
