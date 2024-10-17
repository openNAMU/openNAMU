package route

import (
    "opennamu/route/tool"

    jsoniter "github.com/json-iterator/go"
)

func Api_func_ban(call_arg []string) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(call_arg[0]), &other_set)

    db := tool.DB_connect()
    defer db.Close()

    ip_data := tool.Get_user_ban(db, other_set["ip"], other_set["type"])

    new_data := make(map[string]interface{})
    new_data["response"] = "ok"
    new_data["ban"] = ip_data[0]
    new_data["ban_type"] = ip_data[1]

    json_data, _ := json.Marshal(new_data)
    return string(json_data)
}
