package route

import (
	"opennamu/route/tool"
)

func Api_func_ban(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    ip := config.IP
    if _, exist := other_set["ip"]; exist {
        ip = other_set["ip"]
    }

    ip_data := tool.Get_user_ban(db, ip, other_set["type"])

    new_data := make(map[string]any)
    new_data["response"] = "ok"
    new_data["ban"] = ip_data[0]
    new_data["ban_type"] = ip_data[1]

    json_data, _ := json.Marshal(new_data)
    return string(json_data)
}
