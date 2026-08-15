package route

import (
	"strings"

	"opennamu/route/tool"
	"opennamu/route/tool/markup"
)

func View_render(config tool.Config, doc_name string, rev string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	data_api := Api_w_raw(config, doc_name, "", rev)
	response, _ := data_api["response"].(string)
	if response == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	}
	if response != "ok" {
		return tool.Get_redirect("/history/" + tool.Url_parser(doc_name))
	}

	raw_data, _ := data_api["data"].(string)
	parameter_data := map[string]any{}
	parameter_data["__opennamu_skin_set"] = get_render_setting_parameter(db, config)
	if strings.Contains(strings.ToLower(raw_data), "[username") {
		parameter_data["ip"] = config.IP
	}
	rendered_data := markup.Get_render(db, doc_name, raw_data, "normal", parameter_data)["data"]
	rendered_data = apply_render_setting_data(db, config, rendered_data)
	render_data := get_render_setting_css(db, config) + rendered_data
	return tool.Get_template(
		db,
		config,
		doc_name,
		render_data,
		[]any{"(r" + tool.HTML_escape(rev) + ")"},
		[][]any{{"history/" + tool.Url_parser(doc_name), tool.Get_language(db, "return", true)}},
		map[string]string{},
	)
}
