package route

import (
	"opennamu/route/tool"
)

func Api_func_auth(config tool.Config, ip string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    if ip == "" {
        ip = config.IP
    }

    auth_name := tool.Get_user_auth(db, ip)
    auth_info := tool.Get_auth_group_info(db, auth_name)

    return_data := make(map[string]any)
    return_data["response"] = "ok"
    return_data["name"] = auth_name
    return_data["info"] = auth_info

    return return_data
}
