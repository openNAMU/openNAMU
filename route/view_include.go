package route

import (
	stdjson "encoding/json"

	"opennamu/route/tool"
)

type include_payload struct {
	Version    int               `json:"version"`
	Name       string            `json:"name"`
	Parameters map[string]string `json:"parameters"`
}

func View_include(config tool.Config, payload_data string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	decoded_data, err := tool.Get_base64_decode(payload_data)
	payload := include_payload{}
	if err != nil || stdjson.Unmarshal([]byte(decoded_data), &payload) != nil || payload.Version != 1 || payload.Name == "" {
		return tool.Get_error_page(db, config, "error")
	}

	raw_data_api := Api_w_raw(config, payload.Name, "", "")
	raw_response, _ := raw_data_api["response"].(string)
	if raw_response == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	}
	if raw_response != "ok" {
		error_data := tool.Get_setting_value(db, "error_404", "", "")
		if error_data == "" {
			error_data = tool.Get_language(db, "document_404_error", true)
		}
		error_data = "<h2>" + tool.Get_language(db, "error", true) + "</h2><ul><li>" + error_data + "</li></ul>"
		return tool.Get_template(db, config, tool.Get_language(db, "error", true), error_data, []any{}, [][]any{}, map[string]string{})
	}

	raw_data, _ := raw_data_api["data"].(string)
	parameter_data := map[string]any{"ip": config.IP}
	for key, value := range payload.Parameters {
		parameter_data[key] = value
	}
	option_data, _ := stdjson.Marshal(parameter_data)
	render_data_api := Api_w_render(config, payload.Name, raw_data, "api_include", string(option_data))
	if render_data_api["response"] != "ok" {
		return tool.Get_error_page(db, config, "error")
	}

	render_data, ok := render_data_api["data"].(string)
	if !ok {
		return tool.Get_error_page(db, config, "error")
	}

	return tool.Get_template(
		db,
		config,
		payload.Name,
		render_data,
		[]any{},
		[][]any{{"w/" + tool.Url_parser(payload.Name), tool.Get_language(db, "return", true)}},
		map[string]string{},
	)
}
