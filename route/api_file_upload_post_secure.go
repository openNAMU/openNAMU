package route

import "opennamu/route/tool"

func Api_file_upload_post_secure(config tool.Config, file_name string, file_data string, file_ext string, license string, license_text string, captcha string) map[string]any {
	return api_file_upload_post(config, file_name, file_data, file_ext, license, license_text, captcha, true)
}
