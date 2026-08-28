package route

import (
	stdjson "encoding/json"
	"strings"

	"opennamu/route/tool"
	"opennamu/route/tool/markup"
)

func Api_w_render(config tool.Config, doc_name string, raw_data string, render_type string, option string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if render_type == "backlink" && !tool.Check_acl(db, doc_name, "", "document_edit", config.IP) {
		return map[string]any{"response": "require auth"}
	}

	parameter_data := map[string]any{}
	if option != "" {
		_ = stdjson.Unmarshal([]byte(option), &parameter_data)
	}
	if parameter_data == nil {
		parameter_data = map[string]any{}
	}
	parameter_data["__opennamu_skin_set"] = get_render_setting_parameter(db, config)
	if strings.Contains(strings.ToLower(raw_data), "[username") {
		parameter_data["ip"] = config.IP
	}

	data := markup.Get_render(db, doc_name, raw_data, render_type, parameter_data)
	if render_type != "backlink" {
		rendered_data := data["data"]
		rendered_data = apply_render_setting_data(db, config, rendered_data)
		data["data"] = get_render_setting_css(db, config) + rendered_data
	}

	return_data := make(map[string]any)
	return_data["response"] = "ok"
	for key, value := range data {
		return_data[key] = value
	}

	return return_data
}
