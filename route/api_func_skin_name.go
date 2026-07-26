package route

import (
	"opennamu/route/tool"
	"path/filepath"
)

func Api_func_skin_name(config tool.Config, set_n string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    skin_name := tool.Get_use_skin_name(db, config.IP)

    new_data := make(map[string]any)
    new_data["response"] = "ok"

    if set_n == "0" {
        new_data["data"] = filepath.Join("views", skin_name, "index.html")
    } else {
        new_data["data"] = skin_name
    }

    return new_data
}
