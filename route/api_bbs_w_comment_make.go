package route

import "opennamu/route/tool"

func Api_bbs_w_comment_make(config tool.Config, doc_name string) map[string]any {
	config_copy := config
	config_copy.IP = "Tool:System"

	data_api := Api_bbs_w_post(config_copy, "0", doc_name, "")
	data_api_in, _ := data_api["data"].(string)

	return map[string]any{
		"response": "ok",
		"data":     data_api_in,
	}
}
