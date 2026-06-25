package route

import (
	"opennamu/route/tool"
	"path/filepath"
)

func Api_func_skin_name(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    skin_name := tool.Get_use_skin_name(db, config.IP)

    new_data := make(map[string]string)
    new_data["response"] = "ok"

    if other_set["set_n"] == "0" {
        new_data["data"] = filepath.Join("..", "views", skin_name, "index.html")
    } else {
        new_data["data"] = skin_name
    }

    json_data, _ := json.Marshal(new_data)
    return string(json_data)
}