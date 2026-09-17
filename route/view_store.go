package route

import (
	"strconv"

	"opennamu/route/tool"
)

func View_store(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	api_data := Api_store_list(config)
	data, ok := api_data["data"].(map[string]any)
	if api_data["response"] != "ok" || !ok {
		return tool.Get_error_page(db, config, "error")
	}

	point, _ := data["point"].(int)
	items, _ := data["items"].([]map[string]string)
	body := `<h2>` + tool.Get_language(db, "point", true) + `</h2><div>` + strconv.Itoa(point) + `</div><hr class="main_hr"><h2>` + tool.Get_language(db, "store", true) + `</h2>`
	if len(items) == 0 {
		body += `<div>` + tool.Get_language(db, "store_empty", true) + `</div>`
	}

	return user_form_page(db, config, tool.Get_language(db, "store", true), body)
}
