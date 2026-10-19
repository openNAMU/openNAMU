package route

import (
	"strconv"

	"opennamu/route/tool"
)

func View_setting_backlink_reset_post(config tool.Config, form map[string]string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	load := form["backlink_reset_load"]
	if load != "normal" && load != "slow" {
		load = "fast"
	}
	api_data := Api_setting_backlink_reset_post_internal(config, load)
	response, _ := api_data["response"].(string)
	if response == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	} else if response != "ok" && response != "partial_failure" && response != "failure" {
		return tool.Get_error_page(db, config, "error")
	}

	document_count_value, _ := api_data["document_count"].(int)
	error_count, _ := api_data["error_count"].(int)
	data_html := "<ul>"
	if response == "partial_failure" {
		data_html += "<li>" + tool.Get_language(db, "backlink_reset_partial", true) + "</li>"
	} else if response == "failure" {
		data_html += "<li>" + tool.Get_language(db, "backlink_reset_failed", true) + "</li>"
	}
	data_html += "<li>" + tool.Get_language(db, "reset_all_backlink", true) + " : " + strconv.Itoa(document_count_value) + "</li>"
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
