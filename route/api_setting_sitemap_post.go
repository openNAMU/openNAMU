package route

import (
	"strings"

	"opennamu/route/tool"
)

func Api_setting_sitemap_post(config tool.Config) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)

	if !tool.Check_permission(db, "setting_sitemap", config.IP) {
		return_data["response"] = "require auth"

		return return_data
	}

	sitemap_count, document_count, err := tool.Make_sitemap(db)
	if err != nil {
		if !strings.HasPrefix(err.Error(), "file write error : ") {
			panic(err)
		}

		return_data["response"] = "error"
		return_data["error"] = err.Error()

		return return_data
	}

	return_data["response"] = "ok"
	return_data["sitemap_count"] = sitemap_count
	return_data["document_count"] = document_count

	return return_data
}
