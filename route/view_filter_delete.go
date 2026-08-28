package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_filter_delete(config tool.Config, kind string, name string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	_, ok := get_filter_spec(kind)
	if !ok {
		return tool.Get_error_page(db, config, "error")
	}
	if values == nil && !tool.Check_permission(db, "filter_manage", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	if values != nil {
		result := Api_filter_delete_post(config, kind, name)
		if result["response"] != "ok" {
			return tool.Get_error_page(db, config, "error")
		}
		return tool.Get_redirect("/filter/" + kind)
	}

	body := `<form method="post"><span>` + tool.Get_language(db, "delete_warning", true) + `</span><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "delete", true) + `</button></form>`
	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, "delete", true),
		body,
		[]any{"(" + tool.HTML_escape(name) + ")"},
		[][]any{{"filter/" + kind, tool.Get_language(db, "return", true)}},
		map[string]string{},
	)
}
