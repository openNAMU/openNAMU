package route

import (
    "database/sql"
    "log"

    "opennamu/route/tool"

    jsoniter "github.com/json-iterator/go"
)

func Api_topic(call_arg []string) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(call_arg[0]), &other_set)

    db := tool.DB_connect()
    defer db.Close()

    if other_set["tool"] == "length" {
        stmt, err := db.Prepare(tool.DB_change("select id from topic where code = ? order by id + 0 desc limit 1"))
        if err != nil {
            log.Fatal(err)
        }
        defer stmt.Close()

        var length string
        err = stmt.QueryRow(other_set["topic_num"]).Scan(&length)
        if err != nil {
            if err == sql.ErrNoRows {
                length = "0"
            } else {
                log.Fatal(err)
            }
        }

        new_data := map[string]string{}
        new_data["length"] = length

        json_data, _ := json.Marshal(new_data)
        return string(json_data)
    } else {
        var rows *sql.Rows

        if other_set["tool"] == "top" {
            stmt, err := db.Prepare(tool.DB_change("select id, data, date, ip, block, top from topic where code = ? and top = 'O' order by id + 0 asc"))
            if err != nil {
                log.Fatal(err)
            }
            defer stmt.Close()

            rows, err = stmt.Query(other_set["topic_num"])
            if err != nil {
                log.Fatal(err)
            }
        } else {
            if other_set["s_num"] != "" && other_set["e_num"] != "" {
                stmt, err := db.Prepare(tool.DB_change("select id, data, date, ip, block, top from topic where code = ? and ? + 0 <= id + 0 and id + 0 <= ? + 0 order by id + 0 asc"))
                if err != nil {
                    log.Fatal(err)
                }
                defer stmt.Close()

                rows, err = stmt.Query(other_set["topic_num"], other_set["s_num"], other_set["e_num"])
                if err != nil {
                    log.Fatal(err)
                }
            } else {
                stmt, err := db.Prepare(tool.DB_change("select id, data, date, ip, block, top from topic where code = ? order by id + 0 asc"))
                if err != nil {
                    log.Fatal(err)
                }
                defer stmt.Close()

                rows, err = stmt.Query(other_set["topic_num"])
                if err != nil {
                    log.Fatal(err)
                }
            }
        }
        defer rows.Close()

        data_list := [][]string{}
        ip_parser_temp := map[string][]string{}

        for rows.Next() {
            var id, data, date, ip, block, top string

            err := rows.Scan(&id, &data, &date, &ip, &block, &top)
            if err != nil {
                log.Fatal(err)
            }

            data_list = append(data_list, []string{id, data, date, ip, block, top})
        }

        new_data := make(map[string]interface{})
        new_data["data"] = []map[string]string{}
        data_slice := []map[string]string{}

        admin_auth := tool.Check_acl(db, "", "", "toron_auth", other_set["ip"])

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
                ip_pre = tool.IP_preprocess(db, data_list[for_a][3], other_set["ip"])[0]
                ip_render = tool.IP_parser(db, data_list[for_a][3], other_set["ip"])

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
