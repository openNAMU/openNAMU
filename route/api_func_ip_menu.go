package route

import (
	"opennamu/route/tool"
)

func Api_func_ip_menu(config tool.Config, ip string, option string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    ip_data := tool.IP_menu(db, ip, config.IP, option)

    result_data := make(map[string]any)
    result_data["response"] = "ok"
    result_data["data"] = ip_data

    return result_data
}
