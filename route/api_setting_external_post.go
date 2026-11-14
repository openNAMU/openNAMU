package route

import (
	"database/sql"
	"strings"

	"opennamu/route/tool"
)

func Api_setting_external_post(config tool.Config, form map[string]string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "setting_external", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	provider_raw, provider_exists := form["recaptcha_ver"]
	provider, provider_supported := tool.Captcha_provider_normalize(provider_raw)
	if !provider_exists || !provider_supported {
		return_data["response"] = "error"
		return_data["data"] = "captcha configuration"
		return return_data
	}
	if provider == "" {
		provider = "v2"
	}
	form["recaptcha_ver"] = provider

	for _, name := range []string{"recaptcha", "sec_re", "altcha_sec_re"} {
		value, exists := form[name]
		preserve_current := !exists || (name != "recaptcha" && strings.TrimSpace(value) == "")
		if !preserve_current {
			continue
		}
		current_value, err := tool.Captcha_setting_value_unique(db, name)
		if err != nil {
			return_data["response"] = "error"
			return_data["data"] = "captcha configuration"
			return return_data
		}
		form[name] = current_value
	}
	if err := tool.Captcha_config_validate(provider, form["recaptcha"], form["sec_re"], form["altcha_sec_re"]); err != nil {
		return_data["response"] = "error"
		return_data["data"] = "captcha configuration"
		return return_data
	}
	if form["ai_provider"] != "openai" && form["ai_provider"] != "google" {
		form["ai_provider"] = "ollama"
	}
	for _, name := range []string{"openai_api_key", "google_api_key"} {
		if strings.TrimSpace(form[name]) == "" {
			form[name] = tool.Get_setting_value(db, name, "", "")
		}
	}
	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		Setting_save_fields(tx, Setting_external_fields(), form)
		tool.Do_insert_auth_history(tx, config.IP, "edit_set (external)")
		return nil
	}); err != nil {
		panic(err)
	}
	return_data["response"] = "ok"
	return return_data
}
