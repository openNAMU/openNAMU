package route

import (
	"opennamu/route/tool"
)

func Api_bbs_w(config tool.Config, set_id string, set_code string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    rows := tool.Query_DB(
        db,
        "select set_name, set_data from bbs_data where set_id = ? and set_code = ?",
        set_id,
        set_code,
    )
    defer rows.Close()
    
    data_list := map[string]string{}

    for rows.Next() {
        var set_name string
        var set_data string

        err := rows.Scan(&set_name, &set_data)
        if err != nil {
            panic(err)
        }

        if set_name == "user_id" {
            var ip_pre string
            var ip_render string

            ip_pre = tool.IP_preprocess(db, set_data, config.IP)[0]
            ip_render = tool.IP_parser(db, set_data, config.IP)

            data_list["user_id"] = ip_pre
            data_list["user_id_render"] = ip_render
        } else {
            data_list[set_name] = set_data
        }
    }

    return_data := make(map[string]any)

    if !tool.Check_acl(db, "", "", "bbs_view", config.IP) {
        return_data["response"] = "require auth"
        return_data["data"] = map[string]string{}
    } else {
        return_data["response"] = "ok"
        return_data["data"] = data_list
    }

    return return_data
}
