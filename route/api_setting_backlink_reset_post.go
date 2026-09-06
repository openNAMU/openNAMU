package route

import (
	"strconv"
	"time"

	"opennamu/route/tool"
)

func Api_setting_backlink_reset_post(config tool.Config) map[string]any {
	return api_setting_backlink_reset_post(config, "")
}

func api_setting_backlink_reset_post(config tool.Config, load string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "setting_backlink", config.IP) {
		return_data["response"] = "require auth"

		return return_data
	}

	tool.Exec_DB(
		db,
		"delete from back",
	)
	tool.Exec_DB(
		db,
		"delete from data_set where set_name = 'link_count'",
	)

	document_count := 0
	error_count := 0
	page := 1
	delay := map[string]time.Duration{
		"normal": 10 * time.Millisecond,
		"slow":   100 * time.Millisecond,
	}[load]

	for {
		title_data := Api_list_title_index(config, strconv.Itoa(page))
		title_list := title_data["data"].([]string)

		for _, doc_name := range title_list {
			raw_data := Api_w_raw(config, doc_name, "", "")
			response, ok := raw_data["response"].(string)
			if !ok || response != "ok" {
				error_count++
				continue
			}

			Api_w_render(config, doc_name, raw_data["data"].(string), "backlink", "")
			document_count++
			if delay > 0 {
				time.Sleep(delay)
			}
		}

		if len(title_list) < 50 {
			break
		}

		page++
	}

	return_data["response"] = "ok"
	return_data["document_count"] = document_count
	return_data["error_count"] = error_count

	return return_data
}
