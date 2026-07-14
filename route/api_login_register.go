package route

import "opennamu/route/tool"

func Api_login_register(config tool.Config) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    return_data := make(map[string]any)
    return_data["response"] = "ok"

    return return_data
}