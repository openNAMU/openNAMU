package route

import (
	"opennamu/route/tool"
)

func Api_func_auth_post(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    ip := config.IP
    if _, exist := other_set["ip"]; exist {
        ip = other_set["ip"]
    }

    what := other_set["what"]

    tool.Do_insert_auth_history(db, ip, what)

    new_data := make(map[string]any)
    new_data["response"] = "ok"

    json_data, _ := json.Marshal(new_data)
    return string(json_data)
}
