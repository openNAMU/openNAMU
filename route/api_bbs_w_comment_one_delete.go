package route

import "opennamu/route/tool"

func Api_bbs_w_comment_one_delete(config tool.Config, set_id string, set_code string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    return_data := make(map[string]any)

    api_data := Api_bbs_w_comment_one(config, false, "", set_id + "-" + set_code)
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

    new_id := api_data_in["id"]
    new_code := api_data_in["code"]

    tool.Exec_DB(
        db,
        "update bbs_data set set_data = '' where set_id = ? and set_code = ?",
        new_id,
        new_code,
    )

    return_data["response"] = "ok"

    return return_data
}