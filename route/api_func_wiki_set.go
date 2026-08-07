package route

import (
	"opennamu/route/tool"
)

func Api_func_wiki_set(config tool.Config) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	wiki_set := tool.Get_wiki_set(db, config.IP, config.Cookies)

	new_data := make(map[string]any)
	new_data["response"] = "ok"
	new_data["data"] = wiki_set

	return new_data
}
