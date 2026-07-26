package route

import (
	"opennamu/route/tool"
)

func Api_func_wiki_custom(config tool.Config) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    wiki_custom_set := tool.Get_wiki_custom(db, config.IP, config.Session, config.Cookies)

    new_data := make(map[string]any)
    new_data["response"] = "ok"
    new_data["data"] = wiki_custom_set

    return new_data
}
