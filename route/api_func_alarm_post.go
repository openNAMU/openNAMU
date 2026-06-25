package route

import (
	"opennamu/route/tool"
)

func Api_func_alarm_post(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    tool.Send_alarm(db, other_set["from"], other_set["to"], other_set["data"])

    return_data := make(map[string]any)
    return_data["response"] = "ok"

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}
