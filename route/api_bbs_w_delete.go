package route

import "opennamu/route/tool"

func Api_bbs_w_delete(config tool.Config, set_id string, set_code string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    return_data := make(map[string]any)

    api_data := Api_bbs_w(config, set_id, set_code)
    api_data_in := api_data["data"].(map[string]string)

    if len(api_data_in) == 0 {
        return_data["response"] = "error"
        return_data["data"] = "no data"

        return return_data
    }

    if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
        return_data["response"] = "require auth"
        
        return return_data
    }

    tool.Exec_DB(
        db,
        "delete from bbs_data where set_id = ? and set_code = ?",
        set_id,
        set_code,
    )
    tool.Exec_DB(
        db,
        "delete from bbs_set where set_id = ? and set_code = ?",
        set_id,
        set_code,
    )
    tool.Exec_DB(
        db,
        "delete from bbs_data where set_id = ? or set_id like ?",
        set_id,
        set_code,
    )

    return_data["response"] = "ok"

    return return_data
}