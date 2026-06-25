package route

import (
	"opennamu/route/tool"
)

func Api_w_watch_list_post(config tool.Config, name string, do_type string) map[string]string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    if do_type != "watchlist" {
        do_type = "star_doc"
    }

    var data string

    exist := tool.QueryRow_DB(
        db,
        "select data from user_set where name = ? and id = ? and data = ?",
        []any{ &data },
        do_type,
        config.IP,
        name,
    )
    if exist {
        tool.Exec_DB(
            db,
            "delete from user_set where name = ? and id = ? and data = ?",
            do_type,
            config.IP,
            name,
        )
    } else {
        tool.Exec_DB(
            db,
            "insert into user_set (id, name, data) values (?, ?, ?)",
            config.IP,
            do_type,
            name,
        )
    }

    return_data := make(map[string]string)
    return_data["response"] = "ok"

    return return_data
}