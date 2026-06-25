package route

import "opennamu/route/tool"

func Api_user_setting_skin_set_main(config tool.Config) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    set_list := Get_main_skin_set_list(db)

    return_data := make(map[string]any)
    return_data["response"] = "ok"

    db_set_list := map[string]string{}
    for k := range set_list {
        var data string
        tool.QueryRow_DB(
            db,
            "select data from user_set where name = ? and id = ?",
            []any{ &data },
            k,
            config.IP,
        )

        db_set_list[k] = data
    }

    return_data["data"] = db_set_list

    return return_data
}