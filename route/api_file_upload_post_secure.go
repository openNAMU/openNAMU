package route

import (
	"strings"

	"opennamu/route/tool"
)

func Api_file_upload_post_secure(config tool.Config, file_name string, file_data string, file_ext string, license string, license_text string, captcha string) map[string]any {
	decoded, err := tool.Get_base64_decode(strings.TrimSpace(file_data))
	if err != nil || len(decoded) == 0 {
		return map[string]any{"response": "error", "data": "invalid data"}
	}
	return api_file_upload_post(config, file_name, []byte(decoded), file_ext, license, license_text, captcha, true, false, false)
}
