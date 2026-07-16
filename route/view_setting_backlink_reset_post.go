package route

import (
	"strconv"

	"opennamu/route/tool"
)

func View_setting_backlink_reset_post(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	api_data := Api_setting_backlink_reset_post(config)
	response := api_data["response"].(string)
	if response == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	} else if response != "ok" {
		return tool.Get_error_page(db, config, "error")
	}

	document_count := strconv.Itoa(api_data["document_count"].(int))
	error_count := api_data["error_count"].(int)
	data_html := "<ul><li>" + tool.Get_language(db, "reset_all_backlink", true) + " : " + document_count + "</li>"
	if error_count > 0 {
		data_html += "<li>" + tool.Get_language(db, "error", true) + " : " + strconv.Itoa(error_count) + "</li>"
	}
	data_html += "</ul>"

	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, "reset_all_backlink", true),
		data_html,
		[]any{},
		[][]any{
			{"setting", tool.Get_language(db, "return", true)},
		},
		map[string]string{},
	)
}
