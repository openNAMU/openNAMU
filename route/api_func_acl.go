package route

import (
	"opennamu/route/tool"
)

func Api_func_acl(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    ip := config.IP
    if _, exist := other_set["ip"]; exist {
        ip = other_set["ip"]
    }

    new_data := make(map[string]any)
    new_data["response"] = "ok"
    new_data["data"] = tool.Check_acl(db, other_set["name"], other_set["topic_number"], other_set["tool"], ip)

    json_data, _ := json.Marshal(new_data)
    return string(json_data)
}
