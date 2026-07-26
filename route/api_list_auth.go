package route

import (
	"opennamu/route/tool"
)

func Api_list_auth(config tool.Config) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    data := tool.List_auth(db)

    return_data := make(map[string]any)
    return_data["response"] = "ok"
    return_data["language"] = map[string]string{
        "send" : tool.Get_language(db, "send", false),
        "many_delete_help" : tool.Get_language(db, "many_delete_help", false),
    }
    return_data["data"] = data

    return return_data
}
