package route

import (
	"opennamu/route/tool"
	"strconv"
	"strings"
)

func Api_bbs_w_tabom_post(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    sub_code := other_set["sub_code"]
    sub_code_parts := strings.Split(sub_code, "-")

    bbs_num := ""
    post_num := ""

    if len(sub_code_parts) > 1 {
        bbs_num = sub_code_parts[0]
        post_num = sub_code_parts[1]
    }

    return_data := make(map[string]any)

    if !tool.Check_acl(db, "", "", "bbs_comment", config.IP) {
        return_data["response"] = "require auth"
    } else {
        no_data := ""
        exist := tool.QueryRow_DB(
            db,
            "select set_data from bbs_data where set_name = 'tabom_list' and set_data = ? and set_id = ? and set_code = ?",
            []any{ &no_data },
            config.IP, bbs_num, post_num,
        )

        if !exist {
            return_data["response"] = "ok"
        
            tabom_count := ""
            exsit := tool.QueryRow_DB(
                db,
                "select set_data from bbs_data where set_name = 'tabom_count' and set_id = ? and set_code = ?",
                []any{ &tabom_count },
                bbs_num, post_num,
            )
        
            if !exsit {
                tool.Exec_DB(
                    db,
                    "insert into bbs_data (set_name, set_data, set_id, set_code) values ('tabom_count', ?, ?, ?)",
                    tabom_count, bbs_num, post_num,
                )
            }

            tabom_count_int := tool.Str_to_int(tabom_count)
            tabom_count_int += 1

            tabom_count_str := strconv.Itoa(tabom_count_int)

            tool.Exec_DB(
                db,
                "update bbs_data set set_data = ? where set_name = 'tabom_count' and set_id = ? and set_code = ?",
                tabom_count_str, bbs_num, post_num,
            )
            tool.Exec_DB(
                db,
                "insert into bbs_data (set_name, set_data, set_id, set_code) values ('tabom_list', ?, ?, ?)",
                config.IP, bbs_num, post_num,
            )
        } else {
            return_data["response"] = "same user exist"
        }
    }

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}