package route

import (
	"opennamu/route/tool"
)

func Api_func_ip_post(config tool.Config, data_list []string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    ip_data := map[string]string{}

    for _, val := range data_list {
        if val != "" {
            ip_data[val] = tool.IP_parser(db, val, config.IP)
        }
    }

    new_data := make(map[string]any)
    new_data["response"] = "ok"
    new_data["data"] = ip_data

    return new_data
}
