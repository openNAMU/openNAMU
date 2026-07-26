package route

import (
	"opennamu/route/tool"
)

func Api_func_ban(config tool.Config, ip string, ban_type string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    if ip == "" {
        ip = config.IP
    }

    ip_data := tool.Get_user_ban(db, ip, ban_type)

    new_data := make(map[string]any)
    new_data["response"] = "ok"
    new_data["ban"] = ip_data[0]
    new_data["ban_type"] = ip_data[1]

    return new_data
}
