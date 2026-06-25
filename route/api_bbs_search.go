package route

import "opennamu/route/tool"

func Api_bbs_search(config tool.Config, keyword string, set_id string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    return_data := make(map[string]any)
    return_data["response"] = "ok"

    return return_data
}