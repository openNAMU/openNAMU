package route

import (
	"database/sql"
	"opennamu/route/tool"
)

func Api_bbs_num_to_name(db *sql.DB, set_id string) string {
    bbs_name := ""

    tool.QueryRow_DB(
        db,
        "select set_data from bbs_set where set_id = ? and set_name = 'bbs_name'",
        []any{ &bbs_name },
        set_id,
    )

    return bbs_name
}

func Api_bbs(config tool.Config, bbs_num string, page string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)
    
    rows_arr := []*sql.Rows{}
    if bbs_num == "" {
        rows := tool.Query_DB(
            db,
            "select set_code, set_id, '0' from bbs_data where set_name = 'date' and " + tool.Get_except_set_id_SQL() + " order by set_data desc limit 50",
        )

        rows_arr = append(rows_arr, rows)
    } else {
        page := tool.Str_to_int(page)
        num := 0
        if page * 50 > 0 {
            num = page * 50 - 50
        }

        rows := tool.Query_DB(
            db,
            "select set_code, set_id, '1' from bbs_data where set_name = 'pinned' and set_id like ? order by set_data desc",
            bbs_num,
        )
        
        rows_arr = append(rows_arr, rows)

        rows = tool.Query_DB(
            db,
            "select set_code, set_id, '0' from bbs_data where set_name = 'title' and set_id like ? order by set_code + 0 desc limit ?, 50",
            bbs_num,
            num,
        )

        rows_arr = append(rows_arr, rows)
    }

    data_list := []map[string]string{}
    ip_parser_temp := map[string][]string{}

    for for_a := 0; for_a < len(rows_arr); for_a++ {
        defer rows_arr[for_a].Close()

        for rows_arr[for_a].Next() {
            temp_data := make(map[string]string)

            var set_code string
            var set_id string
            var pinned string

            err := rows_arr[for_a].Scan(&set_code, &set_id, &pinned)
            if err != nil {
                panic(err)
            }

            temp_data["set_code"] = set_code
            temp_data["set_id"] = set_id
            temp_data["pinned"] = pinned

            rows := tool.Query_DB(
                db,
                "select set_name, set_data, set_code, set_id from bbs_data where set_code = ? and set_id = ?",
                set_code,
                set_id,
            )
            defer rows.Close()

            for rows.Next() {
                var set_name string
                var set_data string

                err := rows.Scan(&set_name, &set_data, &set_code, &set_id)
                if err != nil {
                    panic(err)
                }

                if set_name == "user_id" {
                    var ip_pre string
                    var ip_render string

                    if _, ok := ip_parser_temp[set_data]; ok {
                        ip_pre = ip_parser_temp[set_data][0]
                        ip_render = ip_parser_temp[set_data][1]
                    } else {
                        ip_pre = tool.IP_preprocess(db, set_data, config.IP)[0]
                        ip_render = tool.IP_parser(db, set_data, config.IP)

                        ip_parser_temp[set_data] = []string{ip_pre, ip_render}
                    }

                    set_data = ip_pre
                    temp_data["user_id_render"] = ip_render
                }

                if set_name != "data" && set_name != "pinned" {
                    temp_data[set_name] = set_data
                }
            }

            data_list = append(data_list, temp_data)
        }
    }

    return_data := make(map[string]any)
    return_data["response"] = "ok"
    return_data["data"] = data_list

    return return_data
}
