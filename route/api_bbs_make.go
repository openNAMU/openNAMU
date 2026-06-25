package route

import (
	"opennamu/route/tool"
	"strconv"
)

func Api_bbs_make(config tool.Config, bbs_name string, bbs_type string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    return_data := make(map[string]any)

    if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
        return_data["response"] = "require auth"
        
        return return_data
    }

    set_id := "1"
    tool.QueryRow_DB(
        db,
        `select set_id from bbs_set where set_name = "bbs_name" order by set_id + 0 desc`,
        []any{ &set_id },
    )

    set_id_int := tool.Str_to_int(set_id)
    set_id_int += 1

    set_id = strconv.Itoa(set_id_int)

    if !tool.Arr_in_str([]string{ "comment", "thread" }, bbs_type) {
        bbs_type = "comment"
    }

    tool.Exec_DB(
        db,
        "insert into bbs_set (set_name, set_code, set_id, set_data) values ('bbs_name', '', ?, ?)",
        set_id,
        bbs_name,
    )
    tool.Exec_DB(
        db,
        "insert into bbs_set (set_name, set_code, set_id, set_data) values ('bbs_type', '', ?, ?)",
        set_id,
        bbs_type,
    )

    return_data["response"] = "ok"

    return return_data
}