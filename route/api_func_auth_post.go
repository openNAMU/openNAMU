package route

import (
	"opennamu/route/tool"
)

func Api_func_auth_post(config tool.Config, ip string, what string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    if ip == "" {
        ip = config.IP
    }

    tool.Do_insert_auth_history(db, ip, what)

    new_data := make(map[string]any)
    new_data["response"] = "ok"

    return new_data
}
