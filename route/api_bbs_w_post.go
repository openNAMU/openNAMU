package route

import (
	"opennamu/route/tool"
	"strconv"
)

func Api_bbs_w_post(config tool.Config, set_id string, title string, data string) map[string]string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    if !tool.Check_acl(db, "", "", "bbs_comment", config.IP) {
        return_data := make(map[string]string)
        return_data["response"] = "require auth"

        return return_data
    }

    set_code := ""
    tool.QueryRow_DB(
        db,
        "select set_code from bbs_data where set_name = 'title' and set_id = ? order by set_code + 0 desc",
        []any{ &set_code },
        set_id,
    )

    set_code_int := tool.Str_to_int(set_code)
    set_code_int += 1

    set_code_str := strconv.Itoa(set_code_int)

    date_now := tool.Get_time()

    insert_db := [][]string{
        { "title", title },
        { "data", data },
        { "date", date_now },
        { "user_id", config.IP },
    }
    for _, v := range insert_db {
        tool.Exec_DB(
            db,
            "insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)",
            v[0],
            set_code_str,
            set_id,
            v[1],
        )
    }

    return_data := make(map[string]string)
    return_data["response"] = "ok"
    return_data["data"] = set_code_str

    return return_data
}