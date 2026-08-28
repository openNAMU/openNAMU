package route

import (
	"strings"

	"opennamu/route/tool"
)

func View_edit_preview(config tool.Config, doc_name string, data string, mode string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, doc_name, "", "document_edit", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	darkmode := "0"
	preview_name := "preview_normal"
	if mode == "dark" {
		darkmode = "1"
		preview_name = "preview_dark"
	}

	preview_config := config
	preview_config.Cookies += "; main_css_darkmode=" + darkmode

	api_data := Api_w_render(preview_config, doc_name, data, "api_view", "")
	if api_data["response"] != "ok" {
		return tool.Get_error_page(db, config, "error")
	}

	rendered_data, ok := api_data["data"].(string)
	if !ok {
		return tool.Get_error_page(db, config, "error")
	}

	out := tool.Get_template(
		db,
		preview_config,
		doc_name,
		rendered_data,
		[]any{"(" + tool.Get_language(db, preview_name, true) + ")"},
		[][]any{
			{"edit/" + tool.Url_parser(doc_name), tool.Get_language(db, "return", true)},
		},
		map[string]string{},
	)

	return strings.ReplaceAll(out, `<script src="/views/ringo/js/skin_set_do.js.cache_v4"></script>`, "")
}
