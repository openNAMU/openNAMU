package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func Api_topic(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)
    
    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    if other_set["tool"] == "length" {
        length := "0"
        tool.QueryRow_DB(
            db,
            "select id from topic where code = ? order by id + 0 desc limit 1",
            []any{ &length },
            other_set["topic_num"],
        )

        new_data := map[string]string{}
        new_data["length"] = length

        json_data, _ := json.Marshal(new_data)
        return string(json_data)
    } else {
        var rows *sql.Rows

        if other_set["tool"] == "top" {
            rows = tool.Query_DB(
                db,
                "select id, data, date, ip, block, top from topic where code = ? and top = 'O' order by id + 0 asc",
                other_set["topic_num"],
            )
        } else {
            if other_set["s_num"] != "" && other_set["e_num"] != "" {
                rows = tool.Query_DB(
                    db,
                    "select id, data, date, ip, block, top from topic where code = ? and ? + 0 <= id + 0 and id + 0 <= ? + 0 order by id + 0 asc",
                    other_set["topic_num"], other_set["s_num"], other_set["e_num"],
                )
            } else {
                rows = tool.Query_DB(
                    db,
                    "select id, data, date, ip, block, top from topic where code = ? order by id + 0 asc",
                    other_set["topic_num"],
                )
            }
        }
        defer rows.Close()

        data_list := [][]string{}
        ip_parser_temp := map[string][]string{}

        for rows.Next() {
            var id, data, date, ip, block, top string

            err := rows.Scan(&id, &data, &date, &ip, &block, &top)
            if err != nil {
                panic(err)
            }

            data_list = append(data_list, []string{id, data, date, ip, block, top})
        }

        new_data := make(map[string]any)
        new_data["data"] = []map[string]string{}
        data_slice := []map[string]string{}

        admin_auth := tool.Check_acl(db, "", "", "toron_auth", config.IP)

        var ip_pre string
        var ip_render string

        for for_a := 0; for_a < len(data_list); for_a++ {
            data := ""
            if data_list[for_a][4] != "O" || admin_auth {
                data = data_list[for_a][1]
            }

            if _, ok := ip_parser_temp[data_list[for_a][3]]; ok {
                ip_pre = ip_parser_temp[data_list[for_a][3]][0]
                ip_render = ip_parser_temp[data_list[for_a][3]][1]
            } else {
                ip_pre = tool.IP_preprocess(db, data_list[for_a][3], config.IP)[0]
                ip_render = tool.IP_parser(db, data_list[for_a][3], config.IP)

                ip_parser_temp[data_list[for_a][3]] = []string{ip_pre, ip_render}
            }

            data_slice = append(data_slice, map[string]string{
                "id":        data_list[for_a][0],
                "data":      data,
                "date":      data_list[for_a][2],
                "ip":        ip_pre,
                "ip_render": ip_render,
                "blind":     data_list[for_a][4],
            })
        }

        new_data["data"] = data_slice
        new_data["language"] = map[string]string{
            "tool":   tool.Get_language(db, "tool", false),
            "render": tool.Get_language(db, "render", false),
        }

        json_data, _ := json.Marshal(new_data)
        return string(json_data)
    }
}
