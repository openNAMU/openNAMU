package route

import (
	"opennamu/route/tool"
)

func Api_func_email_post(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    err := tool.Send_email(db, config.IP, other_set["who"], other_set["title"], other_set["data"])
    if err == nil {
        new_data := make(map[string]any)
        new_data["response"] = "ok"
    
        json_data, _ := json.Marshal(new_data)
        return string(json_data)
    }

    new_data := make(map[string]any)
    new_data["response"] = "err"
    new_data["data"] = err.Error()

    json_data, _ := json.Marshal(new_data)
    return string(json_data)
}