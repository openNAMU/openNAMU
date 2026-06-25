package route

import "opennamu/route/tool"

func Api_edit_convert(config tool.Config, doc_name string, markup string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    data := Api_w_raw(config, doc_name, "", "")
    raw_data := ""
    if data["response"].(string) == "ok" {
        raw_data = data["data"].(string)
    }

    return_data := make(map[string]any)
    return_data["response"] = "ok"
    return_data["data"] = raw_data

    return return_data
}