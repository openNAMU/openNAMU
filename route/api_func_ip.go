package route

import (
	"opennamu/route/tool"
)

func Api_func_ip(config tool.Config, raw_ip string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    ip_data := tool.IP_parser(db, raw_ip, config.IP)

    result_data := make(map[string]any)
    result_data["response"] = "ok"
    result_data["data"] = ip_data

    return result_data
}
