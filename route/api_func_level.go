package route

import (
	"opennamu/route/tool"
)

func Api_func_level(config tool.Config, ip string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    if ip == "" {
        ip = config.IP
    }

    level := tool.Get_level(db, ip)

    new_data := make(map[string]any)
    new_data["response"] = "ok"
    new_data["data"] = level

    return new_data
}
