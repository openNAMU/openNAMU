package route

import (
	"opennamu/route/tool"
	"opennamu/route/tool/markup"
)

func Api_w_render(config tool.Config, doc_name string, raw_data string, render_type string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	data := markup.Get_render(db, doc_name, raw_data, render_type)

	return_data := make(map[string]any)
	return_data["response"] = "ok"
	for key, value := range data {
		return_data[key] = value
	}

	return return_data
}
